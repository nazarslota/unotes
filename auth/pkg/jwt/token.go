package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// NewHMAC creates and signs a JWT using HMAC algorithm.
func NewHMAC[T jwt.Claims](secret string, claims T) (token string, err error) {
	if token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret)); err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return token, nil
}

// ParseHMAC parses and validates the signature and claims of an access token using HMAC algorithm.
func ParseHMAC[T jwt.Claims](secret string, token string) (claims T, err error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if claims, ok := parsed.Claims.(T); ok && parsed.Valid {
		return claims, nil
	}
	return claims, fmt.Errorf("failed to parse token: %w", err)
}
