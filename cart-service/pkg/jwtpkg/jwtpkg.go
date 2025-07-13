package jwtpkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrClaimsFailed = errors.New("claims failed")
	ErrConversion   = errors.New("conversion failed")
	ErrTokenExpired = errors.New("token expired")
)

type CartJWTpkg struct {
	JWTSecret string
}

func NewCartJWTpkg(JWTSecret string) CartJWTpkg {
	return CartJWTpkg{JWTSecret: JWTSecret}
}

type CartJWTInterface interface {
	// returns role and uuid
	ParseToken(token string) (string, string, error)
}

func (c CartJWTpkg) ParseToken(token string) (string, string, error) {
	keyfunc := func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrUnauthorized
		}
		return []byte(c.JWTSecret), nil
	}

	tkn, err := jwt.Parse(token, keyfunc)
	if err != nil {
		return "", "", err
	}

	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", ErrClaimsFailed
	}

	if claims["exp"].(float64) < float64(time.Now().Unix()) {
		return "", "", ErrTokenExpired
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", "", ErrConversion
	}
	uuid, ok := claims["uuid"].(string)
	if !ok {
		return "", "", ErrConversion
	}

	return role, uuid, nil
}
