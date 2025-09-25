package main

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"

	"github.com/dlsu-lscs/lscs-core-api/internal/server"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	env := os.Getenv("GO_ENV")
	if env == "" || env == "development" {
		err := godotenv.Load() // NOTE: production environment variables are to be placed in Coolify or hosting provider
		if err != nil {
			log.Fatalf("Error loading .env file: %v\n", err)
		}
	}

	srv := server.NewServer()

	c := color.New(color.FgGreen, color.Bold)
	c.Printf("Listening on port %s\n", srv.Addr)

	err := srv.ListenAndServe() // listen on port :42069
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
