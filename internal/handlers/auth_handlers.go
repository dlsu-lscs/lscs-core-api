package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dlsu-lscs/lscs-central-auth-api/internal/database"
	"github.com/dlsu-lscs/lscs-central-auth-api/internal/repository"
	"github.com/dlsu-lscs/lscs-central-auth-api/internal/tokens"
	"github.com/labstack/echo/v4"
)

// only RND members can request API key (to use on their apps)
func RequestAPIKey(c echo.Context) error {
	dbconn := database.Connect()
	q := repository.New(dbconn)

	// parse body
	var req EmailRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		// http.Error(c.Response().Writer, "invalid request body", http.StatusBadRequest)
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "cannot read body",
		})
	}

	memEmail, err := q.CheckEmailIfMember(c.Request().Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{
				"error": "Not an LSCS member",
				"state": "absent",
				"email": memEmail,
			})
		}
		log.Printf("Error checking email: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Internal server error",
		})
	}

	memberInfo, err := q.GetFullMemberInfo(c.Request().Context(), memEmail)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Internal server error",
		})
	}

	// check if RND
	if memberInfo.CommitteeID != "RND" {
		// TODO: maybe return json response then shortcircuit
		http.Error(c.Response().Writer, "Not a Research and Development committee member", http.StatusForbidden)
	}

	// TODO: only generate token if
	hashedToken, rawToken, err := tokens.GenerateToken()
	if err != nil {
		log.Printf("Error generating Refresh Token: %v\n", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Error generating refresh token.",
		})
	}
	// TODO: send rawToken (plaintext) to the client
	// TODO: store hashed token (hashedToken) in db [StoreAPIKey(memEmail, hashedToken, expires_at)]
	return c.JSON(http.StatusOK, echo.Map{
		"api_key": rawToken,
	})
}

func RevokeAPIKey() {
}
