package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
)

var PEPPER = os.Getenv("PEPPER")

// Generates a new API token, hashes it, and returns both the hashed token and raw token, as well as the error if has any.
func GenerateToken() (hashedToken string, rawToken string, err error) {
	rawToken, err = generateRawToken()
	if err != nil {
		return "", "", err
	}

	hashedToken = HashRawToken(rawToken)
	return hashedToken, rawToken, nil
}

// Compares a hashed token from the database with the raw token from the client.
func CompareTokens(dbHashTok, reqTok string) bool {
	return dbHashTok == HashRawToken(reqTok)
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
func HashRawToken(rawToken string) string {
	input := rawToken + PEPPER
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
