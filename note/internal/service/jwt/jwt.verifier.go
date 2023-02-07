package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type Verifier interface {
	Verify(token string) (jwt.MapClaims, error)
}

type verifier struct {
	Secret string
}

func NewVerifier(secret string) Verifier {
	return &verifier{Secret: secret}
}

func (v verifier) Verify(token string) (jwt.MapClaims, error) {
	claims := make(jwt.MapClaims)
	if _, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(v.Secret), nil
	}); err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	return claims, nil
}
