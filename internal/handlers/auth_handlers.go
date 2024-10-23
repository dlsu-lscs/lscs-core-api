package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/dlsu-lscs/lscs-central-auth-api/internal/database"
	"github.com/dlsu-lscs/lscs-central-auth-api/internal/db"
	"github.com/dlsu-lscs/lscs-central-auth-api/internal/tokens"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

// GET: `/authenticate?provider=google` - redirects to Google OAuth
func AuthenticateHandler(c echo.Context) error {
	gothic.BeginAuthHandler(c.Response(), c.Request())
	return nil
}

// GET: `/auth/google/callback` - handle callback, assume user authenticated
func GoogleAuthCallback(c echo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error completing Google authentication")
	}

	// if user.Email does not exist in database, then reject, otherwise accept and generate new JWT with refresh token
	ctx := c.Request().Context()
	dbconn := database.Connect()
	queries := db.New(dbconn)
	email, err := queries.CheckEmailIfMember(ctx, user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{
				"error": "Not an LSCS member",
				"state": "absent",
				"email": email,
			})
		}
		log.Printf("Error checking email: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Internal server error",
		})
	}

	member, err := queries.GetMember(ctx, email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Internal server error",
		})
	}

	jwt, err := tokens.GenerateJWT(email)
	if err != nil {
		log.Printf("Error generating JWT: %v\n", err)
	}
	rt, err := tokens.GenerateRefreshToken(email)
	if err != nil {
		log.Printf("Error generating Refresh Token: %v\n", err)
	}

	// send the JWT signed string (with symmetric key/secret) to client
	// -> return user profile info with JWT token in an HttpOnly cookie
	// TODO: also return member info so use ListInfoMember in queries or sum
	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  jwt,
		"refresh_token": rt,
		"email":         email,
		"success":       "Email is an LSCS member",
		"state":         "present",
		"member_info":   member,
		"google_info":   user,
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

	c.Response().Header().Set("Location", "/")
	return c.NoContent(http.StatusTemporaryRedirect)
}

// type UserFromGoth struct {
// 	RawData           map[string]interface{}
// 	Provider          string
// 	Email             string
// 	Name              string
// 	FirstName         string
// 	LastName          string
// 	NickName          string
// 	Description       string
// 	UserID            string
// 	AvatarURL         string
// 	Location          string
// 	AccessToken       string
// 	AccessTokenSecret string
// 	RefreshToken      string
// 	ExpiresAt         time.Time
// 	IDToken           string
// }
