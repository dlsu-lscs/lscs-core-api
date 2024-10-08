package main

import (
	"log"
	"net/http"
	"os"

	// "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	// echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

// profile: {
//     id: '<unique>',
//     displayName: 'Edwin Sadiarin Jr.',
//     name: {
//         familyName: 'Sadiarin Jr.',
//         givenName: 'Edwin'
//     },
//     emails: [ { value: 'edwin_sadiarinjr@dlsu.edu.ph' } ]
// }

type Profile struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Name        Name
	Emails      []Email
}

type Name struct {
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
}

type Email struct {
	Value string `json:"value"`
}

// example simple struct for User (schema)
type User struct {
	GoogleID  string
	Email     string
	Name      string
	AvatarURL string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.")
	}

	goth.UseProviders(google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), "http://localhost:2323/auth/google/callback", "email", "profile"))

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use(echojwt.WithConfig(echojwt.Config{
	// 	SigningKey: []byte(os.Getenv("JWT_SECRET")),
	// }))

	// Routes
	e.GET("/", hello)
	e.GET("/login", loginHandler)
	e.GET("/auth/google/callback", googleAuthCallback)
	e.POST("/logout", logoutHandler)
	e.GET("/refresh-token", refreshTokenHandler)

	// Start server
	e.Logger.Fatal(e.Start(":2323"))
}

// **** Handlers ****//
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// func jwtTestHandler(c echo.Context) error {
// 	token, err := c.Get("user").(*jwt.Token)
// }

// GET: `/login` - redirects to Google OAuth
func loginHandler(c echo.Context) error {
	gothic.BeginAuthHandler(c.Response(), c.Request())
	return nil
}

// GET: `/auth/google/callback` - handle callback, assume user authenticated
func googleAuthCallback(c echo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error completing Google authentication")
	}
	// store user to db
	// generate JWT with custom claims
	// send the JWT signed string (with symmetric key/secret) to client
	return nil
}

// POST: `/logout` - invalidate session, client-side token invalidation
func logoutHandler(c echo.Context) error {
	// TODO: check if this is redundant
	err := gothic.Logout(c.Response(), c.Request())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Logged out successfully"})
}

func refreshTokenHandler(c echo.Context) error {
	return nil
}
