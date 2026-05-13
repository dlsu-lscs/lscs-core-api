package server

import (
	"net/http"

	"github.com/dlsu-lscs/lscs-core-api/internal/auth"
	"github.com/dlsu-lscs/lscs-core-api/internal/middlewares"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (s *Server) RegisterRoutes(e *echo.Echo) {
	// request ID and logging middlewares
	e.Use(middlewares.RequestIDMiddleware())
	e.Use(middlewares.RequestLoggerMiddleware())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     s.cfg.AllowedOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentLength, echo.HeaderAcceptEncoding, echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true, // required for cookies
	}))

	// Public routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "it works")
	})

	// Swagger documentation
	e.GET("/docs/*", echoSwagger.WrapHandler)

	// --- OAuth routes (public) ---
	authRoutes := e.Group("/auth")
	authRoutes.GET("/google/login", s.oauthHandler.GoogleLoginHandler)
	authRoutes.GET("/google/callback", s.oauthHandler.GoogleCallbackHandler)
	authRoutes.POST("/logout", s.oauthHandler.LogoutHandler)

	// --- Session-protected routes (Web UI) ---
	sessionProtected := e.Group("/auth")
	sessionProtected.Use(middlewares.SessionMiddleware(s.sessionService, s.cfg))
	sessionProtected.GET("/me", s.memberHandler.GetMeHandler)
	sessionProtected.PUT("/me", s.memberHandler.UpdateMeHandler)
	sessionProtected.GET("/members/:id", s.memberHandler.GetMemberByIDHandler)
	sessionProtected.PUT("/members/:id", s.memberHandler.UpdateMemberByIDHandler, middlewares.RequireCanEditMember(s.rbacService))

	// --- Upload routes (Web UI) ---
	uploadProtected := e.Group("/upload")
	uploadProtected.Use(middlewares.SessionMiddleware(s.sessionService, s.cfg))
	uploadProtected.GET("/profile-image/url", s.uploadHandler.GetImageURLHandler)
	uploadProtected.POST("/profile-image", s.uploadHandler.GenerateUploadURLHandler)
	uploadProtected.POST("/profile-image/complete", s.uploadHandler.CompleteUploadHandler)
	uploadProtected.DELETE("/profile-image", s.uploadHandler.DeleteImageHandler)

	// --- API Key routes (Web UI) ---
	// Note: Authorization checks (RND AVP+) are done inline in handlers
	apiKeyProtected := e.Group("/api-keys")
	apiKeyProtected.Use(middlewares.SessionMiddleware(s.sessionService, s.cfg))
	apiKeyProtected.GET("", s.authHandler.ListAPIKeys)
	apiKeyProtected.DELETE("/:id", s.authHandler.RevokeAPIKey)

	// --- API Key Request routes (Web UI) ---
	// Uses session-based auth instead of Bearer tokens for web UI compatibility
	apiRequestKeyProtected := e.Group("/request-key")
	apiRequestKeyProtected.Use(middlewares.SessionMiddleware(s.sessionService, s.cfg))
	apiRequestKeyProtected.POST("", s.authHandler.RequestKeyHandler)

	// --- JWT Protected routes (API Keys) ---
	protected := e.Group("")
	protected.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims { return new(auth.JwtCustomClaims) },
		SigningKey:    []byte(s.cfg.JWTSecret),
		TokenLookup:   "header:Authorization:Bearer ",
		SigningMethod: "HS256",
	}))
	protected.Use(middlewares.JWTEmailMiddleware)
	protected.Use(middlewares.RequireAPIAccess(s.rbacService))

	protected.GET("/members", s.memberHandler.GetAllMembersHandler)
	protected.GET("/committees", s.committeeHandler.GetAllCommitteesHandler)
	protected.POST("/member", s.memberHandler.GetMemberInfo)
	protected.POST("/member-id", s.memberHandler.GetMemberInfoByID)
	protected.POST("/check-email", s.memberHandler.CheckEmailHandler)
	protected.POST("/check-id", s.memberHandler.CheckIDIfMember)
}
