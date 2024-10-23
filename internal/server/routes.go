package server

import (
	"net/http"
	"os"

	"github.com/dlsu-lscs/lscs-central-auth-api/internal/handlers"
	"github.com/dlsu-lscs/lscs-central-auth-api/internal/tokens"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// v1 := e.Group("/v1") // NOTE: if versioning APIs, change param of funcs to group type and use nested groups for routes

	registerAuthRoutes(e)
	registerAdminRoutes(e)

	return e
}

/* Auth Routes */
func registerAuthRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "good")
	})
	e.GET("/authenticate", handlers.AuthenticateHandler) // `/authenticate?provider=google`
	e.GET("/auth/google/callback", handlers.GoogleAuthCallback)
	e.POST("/invalidate", handlers.InvalidateHandler)
}

func registerAdminRoutes(e *echo.Echo) {
	/* Protected Routes */
	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(tokens.JwtCustomClaims)
		},
		SigningMethod: "HS256",
		SigningKey:    []byte(os.Getenv("JWT_SECRET")),
	})

	admin := e.Group("/admin", jwtMiddleware)

	e.GET("/members", handlers.GetAllMembersHandler)
	e.POST("/check-email", handlers.CheckEmailHandler)
	admin.POST("/refresh-token", handlers.RefreshTokenHandler)
	admin.GET("/protected-test", handlers.GetAllMembersHandler)
}
