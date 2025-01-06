package handlers

import (
	"log"
	"net/http"

	"github.com/dlsu-lscs/lscs-central-auth-api/internal/database"
	"github.com/dlsu-lscs/lscs-central-auth-api/internal/repository"
	"github.com/dlsu-lscs/lscs-central-auth-api/internal/tokens"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

// GET or POST: request API key when?
func RequestAPIKey(c echo.Context) error {
	// TODO: only RND members can request API key (to use on their apps)
	dbconn := database.Connect()
	q := repository.New(dbconn)
	q.GetMemberInfo(c.Request().Context(), c.Request().Header["Authorization"])
	jwt, err := tokens.GenerateJWT("")
	if err != nil {
		log.Printf("Error generating JWT: %v\n", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Error generating API Key.",
		})
	}
	rt, err := tokens.GenerateRefreshToken()
	if err != nil {
		log.Printf("Error generating Refresh Token: %v\n", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Error generating refresh token.",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"api_key":       jwt,
		"refresh_token": rt,
	})
}

// POST: `/invalidate` - invalidate session, client-side token invalidation
func InvalidateHandler(c echo.Context) error {
	err := gothic.Logout(c.Response(), c.Request())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to log out from session"})
	}

	// token := c.Get("user").(*jwt.Token)
	// claims := token.Claims.(*JwtCustomClaims)

	// then create a query to invalidate refresh token (requires a refresh_token table in the db)
	return c.NoContent(http.StatusTemporaryRedirect)
}
