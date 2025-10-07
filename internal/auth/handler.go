package auth

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log/slog"
	"net/http"

	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/helpers"
	"github.com/dlsu-lscs/lscs-core-api/internal/repository"
	"github.com/labstack/echo/v4"
)

type EmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type Handler struct {
	authService Service
	dbService   database.Service
}

func NewHandler(authService Service, dbService database.Service) *Handler {
	return &Handler{
		authService: authService,
		dbService:   dbService,
	}
}

func (h *Handler) RequestKeyHandler(c echo.Context) error {
	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)

	// Parse body
	var req EmailRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("cannot read body", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot read body"})
	}

	// Check if user is an LSCS member and in RND
	memberInfo, err := q.GetMemberInfo(c.Request().Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			response := map[string]string{
				"error": "Not an LSCS member",
				"state": "absent",
				"email": req.Email,
			}
			return c.JSON(http.StatusNotFound, response)
		}
		slog.Error("error checking email", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	// only allow RND members and those that are AVPs and above
	if helpers.NullStringToString(memberInfo.CommitteeName) != "Research and Development" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "User is not a member of Research and Development"})
	}

	allowedPositions := map[string]bool{
		"PRES": true,
		"EVP":  true,
		"VP":   true,
		"AVP":  true,
		"CT":   false,
		"JO":   false,
		"MEM":  false,
	}
	pos := helpers.NullStringToString(memberInfo.PositionID)
	if !allowedPositions[pos] {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "User must be AVP or higher"})
	}

	// Generate JWT
	tokenString, err := h.authService.GenerateJWT(memberInfo.Email)
	if err != nil {
		slog.Error("failed to generate token", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error generating token"})
	}

	// Hash the token
	hash := sha256.Sum256([]byte(tokenString))
	hashStr := hex.EncodeToString(hash[:])

	// Store API key
	params := repository.CreateAPIKeyParams{
		MemberEmail: memberInfo.Email,
		ApiKeyHash:  hashStr,
		ExpiresAt:   sql.NullTime{Valid: false},
	}

	_, err = q.CreateAPIKey(c.Request().Context(), params)
	if err != nil {
		slog.Error("failed to store api key", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error storing API key"})
	}

	response := map[string]interface{}{
		"email":   memberInfo.Email,
		"api_key": tokenString,
	}

	return c.JSON(http.StatusOK, response)
}

