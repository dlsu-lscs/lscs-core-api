package tokens

import (
	"fmt"
	"os"
	"time"

	"github.com/dlsu-lscs/lscs-central-auth-api/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateRefreshToken(email string) (string, error) {
	// claims := JwtCustomClaims{
	// 	ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 25 * time.Hour)),
	// 	Subject:   email,
	// }

	claims := models.JwtCustomClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 25 * time.Hour)),
		},
	}
	fmt.Printf("\nGenerated JWT Claims: %+v\n", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Printf("\nGenerated token raw: %+v\n", token)

	tokenSignedString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", nil
	}
	fmt.Printf("\nGenerated JWT Token: %s\n", tokenSignedString)

	return tokenSignedString, nil
}
