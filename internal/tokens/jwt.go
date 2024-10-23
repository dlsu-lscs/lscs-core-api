package tokens

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Email string `json:"email"`
	// default claims like iss, sub, aud, expiresAt, jwtID, etc.
	jwt.RegisteredClaims
}

// Generates JWT with custom claims and returns a signed token string
func GenerateJWT(email string) (string, error) {
	claims := &JwtCustomClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenSignedString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenSignedString, nil
}
