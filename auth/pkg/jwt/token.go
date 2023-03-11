// Package jwt provides functionality for creating and parsing JSON Web Tokens (JWTs) using HMAC-SHA256 for signature
// validation.
package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// ClaimsType is an interface that extends the jwt.Claims interface.
// It is used to specify the type of the claims to be used with the JWT.
type ClaimsType interface {
	jwt.Claims
}

// ClaimsPointerType is an interface that extends the jwt.Claims interface and also includes a pointer to a ClaimsType.
// It is used to specify the type of the claims to be used with the JWT.
type ClaimsPointerType[T ClaimsType] interface {
	jwt.Claims
	*T
}

// NewHMAC creates and returns a new HMAC-SHA256 signed JWT token.
// It takes in a secret string and a set of claims of a type that implements the ClaimsType interface.
// It returns a signed JWT token as a string or an error if the signing process fails.
func NewHMAC[T ClaimsType](secret string, claims T) (token string, err error) {
	if len(secret) == 0 {
		return "", fmt.Errorf("empty secret: %w", jwt.ErrInvalidKey)
	}

	if token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret)); err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return token, nil
}

// ParseHMAC parses a given JWT token string and returns a value with the parsed claims or an error.
// It takes in a secret string, a JWT token string, and two type parameters: T and PT.
// T is the type of the claims, and PT is a pointer to T and implements the ClaimsType interface.
// It returns a T value with the parsed claims or an error if the parsing process fails.
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
