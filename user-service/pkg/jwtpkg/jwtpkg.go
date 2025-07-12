package jwtpkg

import (
	"errors"
	"fmt"
	"time"

	"github.com/Gergenus/commerce/user-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrClaimsFailed = errors.New("claims failed")
)

type UserJWTpkg struct {
	Secret    string
	AccessTTL time.Duration
}

type UserJWTInterface interface {
	GenerateAccessToken(user models.User) (string, error)
	RegenerateToken(oldToken string) (string, error)
}

func NewUserJWTpkg(Secret string, AccessTTL time.Duration) UserJWTpkg {
	return UserJWTpkg{Secret: Secret, AccessTTL: AccessTTL}
}

func (u UserJWTpkg) GenerateAccessToken(user models.User) (string, error) {
	const op = "jwtpkg.GenerateAccessToken"
	claims := jwt.MapClaims{}
	claims["exp"] = time.Now().Add(u.AccessTTL).Unix()
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

func (u UserJWTpkg) RegenerateToken(oldToken string) (string, error) {
	const op = "jwtpkg.RegenerateToken"
	token, err := jwt.Parse(oldToken, func(t *jwt.Token) (interface{}, error) { return []byte(u.Secret), nil }, jwt.WithoutClaimsValidation())
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrClaimsFailed
	}
	UUIDString := claims["uuid"].(string)
	parsedUUID, err := uuid.Parse(UUIDString)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	user := models.User{
		ID:       parsedUUID,
		Role:     claims["role"].(string),
		Verified: claims["verified"].(bool),
	}
	return u.GenerateAccessToken(user)
}
