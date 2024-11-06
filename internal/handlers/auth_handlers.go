package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/dlsu-lscs/lscs-central-auth-api/internal/database"
	"github.com/dlsu-lscs/lscs-central-auth-api/internal/repository"
	"github.com/dlsu-lscs/lscs-central-auth-api/internal/tokens"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

// GET: `/authenticate?provider=google` - redirects to Google OAuth
func AuthenticateHandler(c echo.Context) error {
	// redirectURI := c.QueryParam("redirect_uri")
	// c.Set("redirectURI", redirectURI)
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
	queries := repository.New(dbconn)
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
	rt, err := tokens.GenerateRefreshToken()
	if err != nil {
		log.Printf("Error generating Refresh Token: %v\n", err)
	}

	// c.SetCookie(&http.Cookie{
	// 	Name:   "access_token",
	// 	Value:  jwt,
	// 	Path:   "/",
	// 	Secure: true,
	// })
	//
	// c.SetCookie(&http.Cookie{
	// 	Name:   "refresh_token",
	// 	Value:  rt,
	// 	Path:   "/",
	// 	Secure: true,
	// })
	//
	// c.SetCookie(&http.Cookie{
	// 	Name:   "email",
	// 	Value:  email,
	// 	Path:   "/",
	// 	Secure: true,
	// })

	// redirectURI := c.Get("redirectURI").(string)
	// c.Response().Header().Set("Location", redirectURI)
	// return c.Redirect(http.StatusTemporaryRedirect, redirectURI+"?token="+jwt)

	c.Response().Header().Set("Location", "/successful-redirect")
	c.Set("access_token", jwt)
	c.Set("refresh_token", rt)
	c.Set("email", email)
	c.Set("success", "Email is an LSCS member")
	c.Set("state", "present")
	c.Set("member_info", member)
	c.Set("google_info", user)

	return c.NoContent(http.StatusTemporaryRedirect)
	// return c.NoContent(http.StatusTemporaryRedirect)
	// return c.JSON(http.StatusOK, echo.Map{
	// 	"email":       email,
	// 	"success":     "Email is an LSCS member",
	// 	"state":       "present",
	// 	"member_info": member,
	// 	"google_info": user,
	// })
}

func SuccessfulRedirect(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  c.Get("email"),
		"refresh_token": c.Get("refresh_token"),
		"email":         c.Get("email"),
		"success":       c.Get("success"),
		"state":         c.Get("state"),
		"member_info":   c.Get("member_info"),
		"google_info":   c.Get("google_info"),
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
