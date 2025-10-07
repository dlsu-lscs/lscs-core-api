package server

import (
	"net/http"
	"os"

	"github.com/dlsu-lscs/lscs-core-api/internal/auth"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes(e *echo.Echo) {

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://*", "http://*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentLength, echo.HeaderAcceptEncoding, echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	// Public routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "it works")
	})
	e.POST("/request-key", s.authHandler.RequestKeyHandler)

	// --- Protected routes ----
	protected := e.Group("")
	protected.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims { return new(auth.JwtCustomClaims) },
		SigningKey:    []byte(os.Getenv("JWT_SECRET")),
		TokenLookup:   "header:Authorization:Bearer ",
		SigningMethod: "HS256",
	}))

	protected.GET("/members", s.memberHandler.GetAllMembersHandler)
	protected.GET("/committees", s.committeeHandler.GetAllCommitteesHandler)
	protected.POST("/member", s.memberHandler.GetMemberInfo)
	protected.POST("/member-id", s.memberHandler.GetMemberInfoByID)
	protected.POST("check-email", s.memberHandler.CheckEmailHandler)
	protected.POST("check-id", s.memberHandler.CheckIDIfMember)
}
