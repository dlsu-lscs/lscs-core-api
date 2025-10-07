package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// JwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Service is the interface for the auth service.
// It can be used for mocking.
type Service interface {
	GenerateJWT(email string) (string, error)
}

type service struct {
	jwtSecret []byte
}

// NewService creates a new auth service.
func NewService(secret string) Service {
	return &service{
		jwtSecret: []byte(secret),
	}
}

// GenerateJWT generates a new JWT token and returns it as a string, as well as the error if has any.
func (s *service) GenerateJWT(email string) (string, error) {
	claims := &JwtCustomClaims{
		Email:            email,
		RegisteredClaims: jwt.RegisteredClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

