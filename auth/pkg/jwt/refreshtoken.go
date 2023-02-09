package jwt

import (
	"github.com/golang-jwt/jwt/v4"
)

// RefreshTokenClaims represents the claims in a refresh token.
type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

// RefreshTokenManagerHMAC is a struct for managing refresh tokens using HMAC algorithm.
type RefreshTokenManagerHMAC struct {
	RefreshTokenSecret string
}

// NewRefreshTokenManagerHMAC creates and returns a new RefreshTokenManagerHMAC with the given refresh token secret.
func NewRefreshTokenManagerHMAC(accessTokenSecret string) *RefreshTokenManagerHMAC {
	return &RefreshTokenManagerHMAC{RefreshTokenSecret: accessTokenSecret}
}

// New creates and signs a new refresh token with the given claims.
func (h *RefreshTokenManagerHMAC) New(claims RefreshTokenClaims) (string, error) {
	return NewHMAC(h.RefreshTokenSecret, claims)
}

// Parse parses and validates the signature and claims of a refresh token.
func (h *RefreshTokenManagerHMAC) Parse(token string) (RefreshTokenClaims, error) {
	return ParseHMAC[RefreshTokenClaims](h.RefreshTokenSecret, token)
}
