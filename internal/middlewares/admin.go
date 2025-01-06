package middlewares

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/dlsu-lscs/lscs-central-auth-api/internal/database"
	"github.com/dlsu-lscs/lscs-central-auth-api/internal/repository"
	"github.com/dlsu-lscs/lscs-central-auth-api/internal/tokens"
)

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbconn := database.Connect()
		defer dbconn.Close()
		q := repository.New(dbconn)

		// parse header Authorization
		bearer := r.Header.Get("Authorization")
		if bearer == "" {
			slog.Error("Authorization header missing")
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(bearer, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			slog.Error("Invalid Authorization header format")
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		reqToken := parts[1]
		slog.Info("[Token from Request headers]", "reqToken", reqToken)

		// hash token
		hashedToken, err := tokens.HashRawToken(reqToken)
		if err != nil {
			slog.Error("failed to hash token - adminMiddleware")
			http.Error(w, "failed to hash token", http.StatusInternalServerError)
			return
		}

		info, err := q.GetAPIKeyInfo(r.Context(), hashedToken)
		if err != nil {
			slog.Error("Failed to retrieve API key - Invalid or expired API key")
			http.Error(w, "Invalid or expired API key", http.StatusNotFound)
			return
		}

		// compare hashed token and token from request
		err = tokens.CompareTokens(info.ApiKeyHash, reqToken)
		if err != nil {
			http.Error(w, "Invalid or expired API key", http.StatusNotFound)
			slog.Error("Token comparison failed: ", "error", err)
			return
		}

		if time.Now().After(info.ExpiresAt.Time) {
			http.Error(w, "Token expired, request a new one", http.StatusUnauthorized)
			slog.Error("API Key expired")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Test the middleware with:
//
//    1. A valid token that is not expired.
//    2. An expired token.
//    3. An invalid token (wrong format or tampered).
//    4. Requests without the Authorization header.
