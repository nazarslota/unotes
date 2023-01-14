package oauth2

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func newHS256(secret string, exp time.Duration, userID string) (string, error) {
	if len(secret) == 0 {
		return "", errors.New("invalid secret")
	} else if exp == 0 {
		return "", errors.New("invalid exp")
	} else if len(userID) == 0 {
		return "", errors.New("invalid user id")
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"exp":     now.Add(exp).Unix(),
		"iat":     now.Unix(),
		"nbf":     now.Unix(),
		"user_id": userID,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate a token: %w", err)
	}
	return token, nil
}

func parseHS256(token, secret string) (jwt.MapClaims, error) {
	claims := make(jwt.MapClaims)
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse a token: %w", err)
	}
	return claims, nil
}
