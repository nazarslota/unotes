package oauth2

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrInvalidPassword  = errors.New("invalid password")

	ErrInvalidOrExpiredToken = errors.New("invalid or expired token")
)

type tokens struct{}

func (tokens) NewHS256(secret string, exp time.Duration, userID string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"exp":     now.Add(exp).Unix(),
		"iat":     now.Unix(),
		"nbf":     now.Unix(),
		"user_id": userID,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("filed to create a token: %w", err)
	}
	return token, nil
}

func (tokens) ParseHS256(token, secret string) (jwt.MapClaims, error) {
	claims := make(jwt.MapClaims)
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid or expired token")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse a token: %w", err)
	}
	return claims, nil
}
