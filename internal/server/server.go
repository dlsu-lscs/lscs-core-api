package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dlsu-lscs/lscs-core-api/internal/auth"
	"github.com/dlsu-lscs/lscs-core-api/internal/committee"
	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/member"
	"github.com/labstack/echo/v4"
)

type Server struct {
	port int

	db database.Service

	authHandler      *auth.Handler
	memberHandler    *member.Handler
	committeeHandler *committee.Handler
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	dbService := database.New()

	NewServer := &Server{
		port:             port,
		db:               dbService,
		authHandler:      auth.NewHandler(auth.NewService(os.Getenv("JWT_SECRET")), dbService),
		memberHandler:    member.NewHandler(dbService),
		committeeHandler: committee.NewHandler(dbService),
	}

	// Declare Server config
	e := echo.New()
	NewServer.RegisterRoutes(e)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      e,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
