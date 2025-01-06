package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/dlsu-lscs/lscs-central-auth-api/internal/database"
	"github.com/dlsu-lscs/lscs-central-auth-api/internal/repository"
	"github.com/dlsu-lscs/lscs-central-auth-api/internal/tokens"
)

func AdminMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbconn := database.Connect()
		q := repository.New(dbconn)

		// parse header Authorization
		bearer := r.Header.Get("Authorization")
		if bearer == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			slog.Error("Authorization header missing")
			return
		}
		parts := strings.Split(bearer, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			slog.Error("Invalid Authorization header format")
			return
		}
		reqToken := parts[1]
		slog.Info("[Token from Request headers]", "reqToken", reqToken)

		// hash token
		hashedToken, err := tokens.HashRawToken(reqToken)
		if err != nil {
			http.Error(w, "failed to hash token", http.StatusInternalServerError)
			slog.Error("failed to hash token - adminMiddleware")
			return
		}

		info, err := q.GetAPIKeyInfo(r.Context(), hashedToken)
		if err != nil {
			http.Error(w, "API Key not found", http.StatusNotFound)
			slog.Error("API Key not found")
			return
		}

		// compare hashed token and token from request
		err = tokens.CompareTokens(info.ApiKeyHash, reqToken)
		if err != nil {
			http.Error(w, "API Key is not valid", http.StatusNotFound)
			slog.Error("API Key is not valid")
			return
		}

		if !time.Now().After(info.ExpiresAt.Time) {
			http.Error(w, "Token expired, request a new one", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
