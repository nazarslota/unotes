package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type ClaimsType interface {
	jwt.Claims
}

type ClaimsPointerType[T ClaimsType] interface {
	jwt.Claims
	*T
}

// NewHMAC creates and signs a JWT using HMAC algorithm.
func NewHMAC[T ClaimsType](secret string, claims T) (token string, err error) {
	if token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret)); err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return token, nil
}

// ParseHMAC parses and validates the signature and claims of an access token using HMAC algorithm.
func ParseHMAC[T ClaimsType, PT ClaimsPointerType[T]](secret string, tokenString string) (T, error) {
	var claims PT = new(T)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		return *new(T), fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(PT); ok && token.Valid {
		return *claims, nil
	}
	return *new(T), jwt.ErrTokenInvalidClaims
}
