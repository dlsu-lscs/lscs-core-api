package tokens

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(email string) (string, error) {
	// generate JWT with custom claims and sign it (symmetric key)
	claims := JwtCustomClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	fmt.Printf("\nGenerated JWT Claims: %+v\n", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Printf("\nGenerated token raw: %+v\n", token)

	tokenSignedString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	fmt.Printf("\nGenerated JWT Token: %s\n", tokenSignedString)

	return tokenSignedString, nil
}

func Parse() {
}
