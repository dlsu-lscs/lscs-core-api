package tokens

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

func GenerateRefreshToken() (string, error) {
	rawToken, err := generateRawToken()
	if err != nil {
		return "", nil
	}
	protectedToken, err := hashRawToken(rawToken)

	return protectedToken, nil
}

func CompareTokens(dbHashTok, reqTok string) error {
	return bcrypt.CompareHashAndPassword([]byte(dbHashTok), []byte(reqTok))
}

func generateRawToken() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", nil
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func hashRawToken(password string) (string, error) {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedToken), nil
}
