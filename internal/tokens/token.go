package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

// Generates a new API token, hashes it, and returns both the hashed token and raw token, as well as the error if has any.
func GenerateToken() (hashedToken string, rawToken string, err error) {
	rawToken, err = generateRawToken()
	if err != nil {
		return "", "", err
	}

	hashedToken, err = HashRawToken(rawToken)
	if err != nil {
		return "", "", err
	}

	return hashedToken, rawToken, nil
}

// Compares a hashed token from the database with the raw token from the client.
func CompareTokens(dbHashTok, reqTok string) error {
	return bcrypt.CompareHashAndPassword([]byte(dbHashTok), []byte(reqTok))
}

// Generates base64-encoded, 32-byte string as raw token.
func generateRawToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		slog.Error("failed to generate raw token ", "error", err)
		return "", fmt.Errorf("failed to generate raw token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// Hashes rawToken using bcrypt.
func HashRawToken(rawToken string) (string, error) {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(rawToken), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedToken), nil
}
