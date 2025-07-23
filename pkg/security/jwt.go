package security

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtCustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func JwtGenerateAccessToken(id string, role string) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET_KEY"))

	claims := JwtCustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    "auth",
			Subject:   id,
			Audience:  []string{"web", "mobile"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add((24 * time.Hour) * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}
