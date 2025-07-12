package jwtpkg

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTpkg struct {
	JWTSecret string
	TokenTTL  time.Duration
}

func NewJWTpkg(Secret string, TokenTTL time.Duration) JWTpkg {
	return JWTpkg{JWTSecret: Secret, TokenTTL: TokenTTL}
}

func (j *JWTpkg) GenerateToken(email string) (string, error) {
	const op = "jwtpkg.GenerateToken"
	tkn := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(j.TokenTTL).Unix(),
	})

	token, err := tkn.SignedString([]byte(j.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}
