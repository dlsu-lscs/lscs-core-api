package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"google.golang.org/api/idtoken"
)

// TODO: use this for verifying google emails
func GoogleAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization header is required"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid Authorization header format"})
		}

		tokenString := parts[1]
		audience := os.Getenv("GOOGLE_CLIENT_ID")
		if audience == "" {
			slog.Error("GOOGLE_CLIENT_ID environment variable not set")
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}

		payload, err := idtoken.Validate(context.Background(), tokenString, audience)
		if err != nil {
			slog.Error("failed to validate token", "error", err)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid ID token"})
		}

		email, ok := payload.Claims["email"].(string)
		if !ok || email == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Email not found in token"})
		}

		c.Set("user_email", email)
		return next(c)
	}
}
