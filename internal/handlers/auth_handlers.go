package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/dlsu-lscs/lscs-central-auth-api/internal/database"
	"github.com/dlsu-lscs/lscs-central-auth-api/internal/repository"
	"github.com/dlsu-lscs/lscs-central-auth-api/internal/tokens"
)

type RevokeRequest struct {
	Email           string `json:"email"`
	PepperSecretKey string `json:"pepper"`
}

// only RND members can request API key (to use on their apps)
func RequestAPIKey(w http.ResponseWriter, r *http.Request) {
	dbconn := database.Connect()
	defer dbconn.Close()
	q := repository.New(dbconn)

	// parse body
	var req EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `"error": "cannot read body`, http.StatusNotFound)
		slog.Error("cannot read body", "error", err)
		return
	}

	memEmail, err := q.CheckEmailIfMember(r.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			response := map[string]string{
				"error": "Not an LSCS member",
				"state": "absent",
				"email": memEmail,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}
		http.Error(w, `"error": "Internal server error"`, http.StatusInternalServerError)
		slog.Error("error checking email", "error", err)
		return
	}

	memberInfo, err := q.GetFullMemberInfo(r.Context(), memEmail)
	if err != nil {
		http.Error(w, `"error": "Internal server error"`, http.StatusInternalServerError)
		slog.Error("GetFullMemberInfo error", "error", err)
		return
	}

	// check if RND
	if memberInfo.CommitteeID != "RND" {
		http.Error(w, "Not a Research and Development committee member", http.StatusForbidden)
		slog.Error("Not a Research and Development committee member")
		return
	}

	hashedToken, rawToken, err := tokens.GenerateToken()
	if err != nil {
		http.Error(w, `"error": "Error generating refresh token"`, http.StatusInternalServerError)
		slog.Error("failed to generate token", "error", err)
		return
	}
	slog.Info("generated token", "hashedToken", hashedToken, "rawToken", rawToken)

	newAPIKey := repository.StoreAPIKeyParams{
		MemberEmail: memEmail,
		ApiKeyHash:  hashedToken,
		ExpiresAt: sql.NullTime{
			Time:  time.Time{}, // never expire, let's make it simple for now
			Valid: false,
			// Time:  time.Now().Add(30 * 24 * time.Hour), // 30 days
			// Valid: true,
		},
	}

	if err = q.StoreAPIKey(r.Context(), newAPIKey); err != nil {
		http.Error(w, `"error": "Error storing API key.",`, http.StatusInternalServerError)
		slog.Error("error storing API key", "error", err)
	}

	response := map[string]interface{}{
		"email":   memberInfo.Email,
		"api_key": rawToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// POST: email
func RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	dbconn := database.Connect()
	defer dbconn.Close()
	q := repository.New(dbconn)
	var req RevokeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `"error": "cannot read body`, http.StatusNotFound)
		slog.Error("cannot read body", "error", err)
		return
	}

	if req.PepperSecretKey != os.Getenv("PEPPER") {
		if req.PepperSecretKey == "" {
			http.Error(w, `"error": "no pepper key provided"`, http.StatusForbidden)
			slog.Error("no pepper key provided")
			return
		}
		http.Error(w, `"error": "invalid pepper key"`, http.StatusForbidden)
		slog.Error("invalid pepper key")
		return
	}

	if err := q.DeleteAPIKey(r.Context(), req.Email); err != nil {
		http.Error(w, `"error": "cannot revoke API key"`, http.StatusNotFound)
		slog.Error("cannot revoke API key", "error", err)
		return
	}

	response := fmt.Sprintf("API key for %s is successfully revoked", req.Email)
	// w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
