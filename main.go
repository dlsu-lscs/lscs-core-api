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
	"github.com/markbates/goth/providers/google"
)

type User struct {
	GoogleID string
	Email    string
	Name     string
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
	// e.GET("/login", loginHandler)
	// e.GET("/auth/google/callback", googleAuthCallback)
	// e.GET("/logout", logoutHandler)
	// e.GET("/refresh-token", refreshTokenHandler)

	// Start server
	e.Logger.Fatal(e.Start(":2323"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// func jwtTestHandler(c echo.Context) error {
// 	token, err := c.Get("user").(*jwt.Token)
// }

func loginHandler(c echo.Context) error {
	// call googleAuthCallback
	// if successful login:
	// - store user to db
	// - generate JWT with custom claims
	// - send the JWT signed string (with symmetric key/secret) to client
	return nil
}
