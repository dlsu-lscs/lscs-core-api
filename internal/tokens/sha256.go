package tokens

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ( // TODO: change this to something more secure
	JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))
	PEPPER     = os.Getenv("PEPPER")
)

// Generates a new JWT token and returns it as a string, as well as the error if has any.
func GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
	})

	tokenString, err := token.SignedString(JWT_SECRET)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

// Hashes rawToken using bcrypt.
func HashRawToken(rawToken string) string {
	input := rawToken + PEPPER
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}