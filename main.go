package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

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

	dbpool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create database connection pool: %v", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	err = createUsersTable()
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	goth.UseProviders(google.New(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		"http://localhost:2323/auth/google/callback", // TODO: add prod callback in google console
		"email", "profile", // scopes - can add openIDConnect if necessary
	))

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/login", loginHandler) // NOTE: use `/login?provider=google` when calling
	e.GET("/auth/google/callback", googleAuthCallback)
	e.POST("/logout", logoutHandler)
	// e.GET("/refresh-token", refreshTokenHandler)

	needsJWT := e.Group("/auth")
	needsJWT.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))
	needsJWT.GET("/profile", profileHandler)

	// Start server
	e.Logger.Fatal(e.Start(":2323"))
}

// **** Handlers ****//
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func createUsersTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS users (
            google_id VARCHAR(255) PRIMARY KEY,
            email VARCHAR(255) NOT NULL,
            name VARCHAR(255),
            avatar_url VARCHAR(255)
        );
    `
	_, err := dbpool.Exec(context.Background(), query)
	if err != nil {
		log.Printf("Failed to create users table")
		return fmt.Errorf("Failed to create users table: %w", err)
	}
	return nil
}

// GET: `/login?provider=google` - redirects to Google OAuth
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
	fmt.Printf("\nUser: %v\n", user)
	// store user to db --> should save to postgresql
	err = saveUser(&user)
	if err != nil {
		log.Printf("Error saving user to database: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error when saving user to database"})
	}

	// generate JWT with custom claims and sign it (symmetric key)
	claims := JwtCustomClaims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)), // TODO: is this proper expiration time for the JWT
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSignedString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
		return c.JSON(http.StatusInternalServerError, "Error generating JWT")
	}
	fmt.Printf("\nGenerated JWT Token: %s\n", tokenSignedString)

	// send the JWT signed string (with symmetric key/secret) to client
	// -> return user profile info with JWT token in an HttpOnly cookie
	return c.JSON(http.StatusOK, echo.Map{
		"token": tokenSignedString,
		"user":  user,
	})
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

// func refreshTokenHandler(c echo.Context) error {
// 	return nil
// }

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
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)
	// - [ ] get claims -> retrieve the email
	// - [ ] query SELECT the user profile info
	// - [ ] return JSON
	fmt.Printf("Received JWT Claims: %v\n", claims)
	var google_id, email, name, avatar_url string
	query := `SELECT google_id, email, name, avatar_url FROM users WHERE email = $1 `
	err := dbpool.QueryRow(context.Background(), query, claims.Email).Scan(&google_id, &email, &name, &avatar_url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error retrieving info from database."})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"googleID":  google_id,
		"email":     email,
		"name":      name,
		"avatarURL": avatar_url,
	})
}
