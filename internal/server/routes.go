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
	"golang.org/x/time/rate"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))

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
	e.GET("/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"test": "tseter"})
	})
	e.GET("/request-key", handlers.RequestAPIKey) // `/request_key?` TODO: change to POST if need condition before able to request, ex. need to be admin email only
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

	e.GET("/members", handlers.GetAllMembersHandler, jwtMiddleware)
	e.POST("/member", handlers.GetMemberInfo, jwtMiddleware)
	e.POST("/check-email", handlers.CheckEmailHandler, jwtMiddleware)
	e.POST("/refresh-token", handlers.RefreshTokenHandler, jwtMiddleware)
	e.GET("/protected-test", handlers.GetAllMembersHandler, jwtMiddleware)
}
