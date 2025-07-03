package jwtpkg

import (
	"fmt"
	"time"

	"github.com/Gergenus/commerce/user-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type UserJWTpkg struct {
	Secret string
}

type UserJWTInterface interface {
	GenerateAccessToken(user models.User) (string, error)
}

func NewUserJWTpkg(Secret string) UserJWTpkg {
	return UserJWTpkg{Secret: Secret}
}

func (u UserJWTpkg) GenerateAccessToken(user models.User) (string, error) {
	const op = "jwtpkg.GenerateAccessToken"
	claims := jwt.MapClaims{}
	claims["exp"] = time.Now().Add(15 * time.Minute).Unix()
	claims["uuid"] = user.ID.String()
	claims["verified"] = user.Verified
	claims["role"] = user.Role
	AccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	AccessTokenString, err := AccessToken.SignedString([]byte(u.Secret))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return AccessTokenString, nil
}
