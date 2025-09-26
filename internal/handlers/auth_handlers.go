package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/repository"
	"github.com/dlsu-lscs/lscs-core-api/internal/tokens"
)

func RequestKeyHandler(w http.ResponseWriter, r *http.Request) {
	dbconn := database.Connect()
	defer dbconn.Close()
	q := repository.New(dbconn)

	// Parse body
	var req EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "cannot read body"}`, http.StatusBadRequest)
		slog.Error("cannot read body", "error", err)
		return
	}

	// Check if user is an LSCS member
	memEmail, err := q.CheckEmailIfMember(r.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			response := map[string]string{
				"error": "Not an LSCS member",
				"state": "absent",
				"email": req.Email,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		slog.Error("error checking email", "error", err)
		return
	}

	// Generate JWT
	tokenString, err := tokens.GenerateJWT(memEmail)
	if err != nil {
		http.Error(w, `{"error": "Error generating token"}`, http.StatusInternalServerError)
		slog.Error("failed to generate token", "error", err)
		return
	}

	response := map[string]interface{}{
		"email": memEmail,
		"token": tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
