package routes

import (
	"net/http"
	"os"

	"github.com/dlsu-lscs/lscs-central-auth-api/internal/controllers"
	"github.com/dlsu-lscs/lscs-central-auth-api/internal/tokens"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitializeRoutes() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS, echo.PATCH},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	// v1 := e.Group("/v1")
	registerAuthRoutes(e)
	registerAdminRoutes(e)

	return e
}

func registerAuthRoutes(e *echo.Echo) {
	// test route
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "good")
	})
	e.GET("/authenticate", controllers.AuthenticateHandler) // `/authenticate?provider=google`
	e.GET("/auth/google/callback", controllers.GoogleAuthCallback)
	e.POST("/invalidate", controllers.InvalidateHandler)
}

func registerAdminRoutes(e *echo.Echo) {
	// TODO: should these routes be protected?
	// --> reasoning: only authorized roles can access this?
	// NOTE: these should also be protected if ever
	e.GET("/members", controllers.GetAllMembersHandler)
	e.POST("/check-email", controllers.CheckEmailHandler)
	e.POST("/refresh-token", controllers.RefreshTokenHandler)

	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(tokens.JwtCustomClaims)
		},
		SigningMethod: "HS256",
		SigningKey:    []byte(os.Getenv("JWT_SECRET")),
	})
	e.GET("/protected-test", controllers.GetAllMembersHandler, jwtMiddleware) // NOTE: exmaple protected /admin/profile route for testing
}
