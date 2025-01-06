package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/dlsu-lscs/lscs-central-auth-api/internal/database"
	"github.com/dlsu-lscs/lscs-central-auth-api/internal/repository"
)

type EmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func GetMemberInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dbconn := database.Connect()
	defer dbconn.Close()
	q := repository.New(dbconn)

	req := new(EmailRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		slog.Error("Failed to parse request body", "err", err)
		http.Error(w, `"error": "Invalid request format"`, http.StatusBadRequest)
		return
	}

	memberInfo, err := q.GetMemberInfo(ctx, req.Email)
	if err != nil {
		slog.Error("email is not an LSCS member", "err", err)
		http.Error(w, `"error": "Email is not an LSCS member"`, http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"email":          memberInfo.Email,
		"full_name":      memberInfo.FullName,
		"committee_name": memberInfo.CommitteeName,
		"division_name":  memberInfo.DivisionName,
		"position_name":  memberInfo.PositionName,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
	}
}

type MemberResponse struct {
	ID          int32            `json:"id"`
	FullName    string           `json:"full_name"`
	Nickname    CustomNullString `json:"nickname"`
	Email       string           `json:"email"`
	Telegram    CustomNullString `json:"telegram"`
	PositionID  CustomNullString `json:"position_id"`
	CommitteeID CustomNullString `json:"committee_id"`
	College     CustomNullString `json:"college"`
	Program     CustomNullString `json:"program"`
	Discord     CustomNullString `json:"discord"`
}

// NOTE: wraps sql.NullString to customize JSON marshaling - to omit sql.NullString struct values on response
type CustomNullString struct {
	sql.NullString
}

// NOTE: implements the json.Marshaler interface - if nil (sql.Null) then "", otherwise return orig val
func (cns CustomNullString) MarshalJSON() ([]byte, error) {
	if cns.Valid {
		return json.Marshal(cns.String)
	}
	return json.Marshal("")
}

func toMemberResponse(m repository.Member) MemberResponse {
	return MemberResponse{
		ID:          m.ID,
		FullName:    m.FullName,
		Nickname:    CustomNullString{m.Nickname},
		Email:       m.Email,
		Telegram:    CustomNullString{m.Telegram},
		PositionID:  CustomNullString{m.PositionID},
		CommitteeID: CustomNullString{m.CommitteeID},
		College:     CustomNullString{m.College},
		Program:     CustomNullString{m.Program},
		Discord:     CustomNullString{m.Discord},
	}
}

func GetAllMembersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dbconn := database.Connect()
	defer dbconn.Close()
	queries := repository.New(dbconn)

	members, err := queries.ListMembers(ctx)
	if err != nil {
		slog.Error("Failed to list members", "err", err)
		http.Error(w, `"error": "Failed to list members"`, http.StatusInternalServerError)
		return
	}

	response := make([]MemberResponse, 0, len(members))
	for _, m := range members {
		response = append(response, toMemberResponse(m))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
	}
}

// this will be included in the google auth callback handler
func CheckEmailHandler(w http.ResponseWriter, r *http.Request) {
	var req EmailRequest // req := new(EmailRequest)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("invalid request body")
		http.Error(w, `"error": "Invalid request body"`, http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		slog.Error("email is required")
		http.Error(w, `"error": "Email is required"`, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	dbconn := database.Connect()
	defer dbconn.Close()
	queries := repository.New(dbconn)
	memberEmail, err := queries.CheckEmailIfMember(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			response := map[string]string{
				"error": "Not an LSCS member",
				"state": "absent",
				"email": memberEmail,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}
		log.Printf("Error checking email: %v", err)
		http.Error(w, `"error": "Internal server error"`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": "Email is an LSCS member",
		"state":   "present",
		"email":   memberEmail,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
	}
}

// GetAllCommitteesHandler retrieves and returns all committees.
func GetAllCommitteesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dbconn := database.Connect()
	defer dbconn.Close()
	q := repository.New(dbconn)

	committees, err := q.GetAllCommittees(ctx)
	if err != nil {
		response := map[string]string{
			"error": fmt.Sprintf("Internal server error: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"committees": committees,
	})
}
