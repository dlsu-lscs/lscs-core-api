package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dlsu-lscs/lscs-core-api/internal/auth"
	"github.com/dlsu-lscs/lscs-core-api/internal/committee"
	"github.com/dlsu-lscs/lscs-core-api/internal/config"
	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/member"
	"github.com/dlsu-lscs/lscs-core-api/internal/storage"
	"github.com/labstack/echo/v4"
)

type Server struct {
	port int
	cfg  *config.Config

	db database.Service

	// handlers
	authHandler      *auth.Handler
	oauthHandler     *auth.OAuthHandler
	memberHandler    *member.Handler
	committeeHandler *committee.Handler
	uploadHandler    *storage.UploadHandler

	// services
	sessionService auth.SessionService
	rbacService    *auth.RBACService
	s3Service      *storage.S3Service
}

func NewServer(cfg *config.Config) *http.Server {
	dbService := database.New(cfg)

	// create session service
	sessionService := auth.NewSessionService(dbService.GetConnection(), cfg)

	// create RBAC service
	rbacService := auth.NewRBACService(dbService)

	// create S3 storage service
	s3Service, err := storage.NewS3Service(cfg)
	if err != nil {
		// log but continue - storage is optional
		_ = err
	}

	// create upload handler
	uploadHandler := storage.NewUploadHandler(s3Service, dbService, cfg)

	// start session cleanup job (runs every hour)
	ctx := context.Background()
	auth.StartCleanupJob(ctx, sessionService, 1*time.Hour)

	NewServer := &Server{
		port:             cfg.Port,
		cfg:              cfg,
		db:               dbService,
		sessionService:   sessionService,
		rbacService:      rbacService,
		s3Service:        s3Service,
		uploadHandler:    uploadHandler,
		authHandler:      auth.NewHandler(auth.NewService(cfg.JWTSecret, cfg), dbService, rbacService),
		oauthHandler:     auth.NewOAuthHandler(cfg, sessionService, dbService),
		memberHandler:    member.NewHandler(dbService),
		committeeHandler: committee.NewHandler(dbService),
	}

	// Declare Server config
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	NewServer.RegisterRoutes(e)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      e,
		IdleTimeout:  cfg.ServerIdleTimeout,
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
	}

	return server
}
