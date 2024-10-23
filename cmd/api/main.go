package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"

	"github.com/dlsu-lscs/lscs-central-auth-api/internal/server"

	_ "github.com/go-sql-driver/mysql"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func main() {
	env := os.Getenv("GO_ENV") // NOTE: environment variables are to be placed in Coolify or hosting provider
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
		os.Getenv("CALLBACK_URL"), // TODO: add prod callback url in google console and put in .env
		"email", "profile",
	))

	srv := server.NewServer()

	c := color.New(color.FgGreen, color.Bold)
	c.Printf("Listening on port %s\n", srv.Addr)
	// Start server on port :42069 ...yeah
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf(fmt.Sprintf("Server error: %v", err))
	}
}
