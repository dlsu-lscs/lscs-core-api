package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dlsu-lscs/lscs-core-api/internal/database"
)

type Server struct {
	dbconn *sql.DB
	port   int
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	// port := os.Getenv("PORT")

	NewServer := &Server{
		dbconn: database.Connect(),
		port:   port,
	}

	// server configs here
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
