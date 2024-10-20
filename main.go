package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"

	"github.com/dlsu-lscs/lscs-central-auth-api/internal/db"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

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

var dbconn *sql.DB

func main() {
	env := os.Getenv("GO_ENV") // this is included in Dockerfile or in Env vars in Coolify
	if env == "" {
		env = "development"
	}

	var envFile string
	switch env {
	case "production":
		envFile = ".env.production" // TODO: GO_ENV = "production" should be set on coolify or dockerfile for production to work
	default:
		envFile = ".env.development"
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	fmt.Printf("ENV: %v (using %v)\n", env, envFile)

	// MySQL connection string format: username:password@tcp(host:port)/dbname
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	dbconn, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	dbconn.SetConnMaxLifetime(0)
	dbconn.SetMaxIdleConns(50)
	dbconn.SetMaxOpenConns(50)
	defer dbconn.Close()

	// test connection
	if err := dbconn.Ping(); err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
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
	e.GET("/members", getAllMembersHandler)
	e.POST("/check-email", checkEmailHandler)
	e.GET("/login?provider=google", loginHandler)
	e.GET("/auth/google/callback", googleAuthCallback)
	e.POST("/invalidate", invalidateHandler)
	// e.GET("/refresh", refreshTokenHandler)

	needsJWT := e.Group("/auth")
	needsJWT.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		SigningMethod: "HS256",
		SigningKey:    []byte(os.Getenv("JWT_SECRET")),
	}))
	// needsJWT.GET("/profile", profileHandler) // NOTE: exmaple protected /auth/profile route for testing

	// Start server
	e.Logger.Fatal(e.Start(":2323"))
}

// **** Handlers ****//
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// GET: `/login?provider=google` - redirects to Google OAuth
func loginHandler(c echo.Context) error {
	// TODO: add redirection to `/login?provider=google` when hitting raw `/login`
	gothic.BeginAuthHandler(c.Response(), c.Request())
	return nil
}

// GET: `/auth/google/callback` - handle callback, assume user authenticated
func googleAuthCallback(c echo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error completing Google authentication")
	}
	fmt.Printf("\nUser: %v\n", user)

	// TODO: add email check here if user that verified is an LSCS Officer
	// if user.Email does not exist in database, then reject, otherwise accept and generate new JWT with refresh token
	// q := `SELECT `

	// generate JWT with custom claims and sign it (symmetric key)
	claims := &JwtCustomClaims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	fmt.Printf("\nGenerated JWT Claims: %+v\n", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Printf("\nGenerated token raw: %+v\n", token)

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

// POST: `/invalidate` - invalidate session, client-side token invalidation
func invalidateHandler(c echo.Context) error {
	// TODO: logoutHandler: check if this is redundant
	err := gothic.Logout(c.Response(), c.Request())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Logged out successfully"})
}

func getAllMembersHandler(c echo.Context) error {
	ctx := c.Request().Context()

	queries := db.New(dbconn)

	members, err := queries.ListMembers(ctx)
	if err != nil {
		log.Printf("Failed to list members: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to list members"})
	}

	return c.JSON(http.StatusOK, members)
}

type EmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func checkEmailHandler(c echo.Context) error {
	req := new(EmailRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}
	if req.Email == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Email is required"})
	}

	ctx := c.Request().Context()
	queries := db.New(dbconn)
	memberEmail, err := queries.CheckEmailIfMember(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{
				"error": "Not an LSCS member",
				"state": "absent",
				"email": memberEmail,
			})
		}
		log.Printf("Error checking email: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"success": "Email is an LSCS member",
		"state":   "present",
		"email":   memberEmail,
	})
}
