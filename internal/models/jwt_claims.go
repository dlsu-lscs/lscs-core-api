package models

import "github.com/golang-jwt/jwt/v5"

// NOTE: can add API tokens / custom keys here
type JwtCustomClaims struct {
	Email string `json:"email"`
	// default claims like iss, sub, aud, expiresAt, jwtID, etc.
	jwt.RegisteredClaims
}
