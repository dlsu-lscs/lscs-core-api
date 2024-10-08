package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/golang-jwt/jwt/v5"
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

// TODO: DO I NEED THIS WTFOK?!?!?!???????
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

// this is returned by goth.User struct:

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

// NOTE: can add API tokens / custom keys here
type JwtCustomClaims struct {
	Email string `json:"email"`
	// default claims like iss, sub, aud, expiresAt, jwtID, etc.
	jwt.RegisteredClaims
}

var dbpool *pgxpool.Pool

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.")
	}

	// jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create database connection pool: %v", err)
		os.Exit(1)
	}
	defer dbpool.Close()

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
	e.GET("/profile", profileHandler)

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
    // TODO: googleAuthCallback: add JWT generation
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error completing Google authentication")
	}
	fmt.Printf("User: %v", user)
	// store user to db
	// - saveUser(&user)
	// --> should save to postgresql

	// generate JWT with custom claims

	// send the JWT signed string (with symmetric key/secret) to client

	return nil
}

// POST: `/logout` - invalidate session, client-side token invalidation
func logoutHandler(c echo.Context) error {
    // TODO: logoutHandler: check if this is redundant
	err := gothic.Logout(c.Response(), c.Request())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Logged out successfully"})
}

func refreshTokenHandler(c echo.Context) error {
	return nil
}

func saveUser(user *goth.User) error {
	query := `
        INSERT INTO users (google_id, email, name, avatar_url)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (google_id) DO NOTHING;
    `
	_, err := dbpool.Exec(context.Background(), query, user.UserID, user.Email, user.Name, user.AvatarURL)
	if err != nil {
		log.Printf("Error saving user to database: %v", err)
		return err
	}

	return nil
}

// protected route for testing JWT - returns user profile
func profileHandler(c echo.Context) error {
    // TODO: profileHandler: profile tasks
	// - [ ] get user token
    userToken := c.Get("")
	// - [ ] get claims -> retrieve the email
	// - [ ] query SELECT the user profile info
	// - [ ] return JSON
    query := `SELECT google_id, email, name, avatar_url FROM users WHERE email = $1 `
    dbpool.Exec(context.Background(), query, claim)

    
	return c.JSON(http.StatusOK, echo.Map{
        "googleID": ,
        "email": ,
        "name": ,
        "avatarURL": ,
    })
}
