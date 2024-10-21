package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
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

	ss := os.Getenv("SESSION_SECRET")
	if ss == "" {
		log.Fatal("No session secret configured.")
	}

	// NOTE: this is for admin things (for the future, if need admin micro-frontend for adding new LSCS members to central auth database)
	store := sessions.NewCookieStore([]byte(ss)) // maybe use redis for session management (storing session data), but thats future ppl problems kekw
	store.Options.HttpOnly = true
	store.Options.Path = "/"
	store.Options.MaxAge = 0
	gothic.Store = store

	goth.UseProviders(google.New(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		"http://localhost:42069/auth/google/callback", // TODO: add prod callback in google console
		"email", "profile",
	))

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Main Authentication Routes
	e.GET("/authenticate", authenticateHandler) // `/authenticate?provider=google`
	e.GET("/auth/google/callback", googleAuthCallback)

	// TODO: should these routes be protected?
	// --> reasoning: only authorized roles can access this?
	e.GET("/", hello)
	e.GET("/members", getAllMembersHandler)
	e.POST("/check-email", checkEmailHandler)
	e.POST("/invalidate", invalidateHandler)
	e.POST("/refresh-token", refreshTokenHandler)

	adminGroup := e.Group("/admin")
	adminGroup.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		SigningMethod: "HS256",
		SigningKey:    []byte(os.Getenv("JWT_SECRET")),
	}))
	// adminGroup.GET("/profile", profileHandler) // NOTE: exmaple protected /admin/profile route for testing

	// Start server
	e.Logger.Fatal(e.Start(":42069"))
}

// **** Handlers ****//
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// GET: `/authenticate?provider=google` - redirects to Google OAuth
func authenticateHandler(c echo.Context) error {
	// if user, err := gothic.CompleteUserAuth(c.Response(), c.Request()); err == nil {
	// 	fmt.Printf("Already authenticated: %v\n", user)
	// 	return c.JSON(http.StatusOK, echo.Map{ // TODO: if doesn't work then maybe check to database?
	// 		"msg":  "Already authenticated",
	// 		"data": user,
	// 	})
	// }
	gothic.BeginAuthHandler(c.Response(), c.Request())
	return nil
}

// GET: `/auth/google/callback` - handle callback, assume user authenticated
func googleAuthCallback(c echo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error completing Google authentication")
	}

	// if user.Email does not exist in database, then reject, otherwise accept and generate new JWT with refresh token
	ctx := c.Request().Context()
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

	// WARN: START: remove when modularizing (see code below)
	claims := &JwtCustomClaims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenSignedString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
		return c.JSON(http.StatusInternalServerError, "Error generating JWT")
	}
	// WARN: END: remove when modularizing (see code below)

	// NOTE: call this when modularizing
	// jwt, err := tokens.GenerateJWT(memberEmail)
	// if err != nil {
	// 	log.Printf("Error generating JWT: %v\n", err)
	// }
	// rt, err := tokens.GenerateRefreshToken(memberEmail)
	// if err != nil {
	// 	log.Printf("Error generating Refresh Token: %v\n", err)
	// }

	// send the JWT signed string (with symmetric key/secret) to client
	// -> return user profile info with JWT token in an HttpOnly cookie
	// TODO: also return member info so use ListInfoMember in queries or sum
	return c.JSON(http.StatusOK, echo.Map{
		// "access_token": jwt,
		// "refresh_token": rt,
		"access_token": tokenSignedString,
		"user":         user,
		"email":        email,
		"success":      "Email is an LSCS member",
		"state":        "present",
	})
}

// POST: `/invalidate` - invalidate session, client-side token invalidation
func invalidateHandler(c echo.Context) error {
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

// this will be included in the google auth callback handler
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

// TODO: implement this for admin access route groups
func refreshTokenHandler(c echo.Context) error {
	// get refresh token from request header frfr
	// get hashed token from database
	// call CompareTokens to compare
	// if valid, tokens.GenerateJWT (generate new access token)
	// --> also generate new refreshToken maybe (call tokens.GenerateRefreshToken)
	// --> then store newRefreshToken in the database
	return c.JSON(http.StatusOK, echo.Map{
		"access_token": "return new access token here", // TODO: handle refreshing tokens
		// "refresh_token": newRefreshToken,
	})
}
