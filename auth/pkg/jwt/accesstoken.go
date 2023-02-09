package jwt

import (
	"github.com/golang-jwt/jwt/v4"
)

// AccessTokenClaims represents the claims in an access token.
type AccessTokenClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

// AccessTokenManagerHMAC is a struct for managing access tokens using HMAC algorithm.
type AccessTokenManagerHMAC struct {
	AccessTokenSecret string
}

// NewAccessTokenManagerHMAC creates and returns a new AccessTokenManagerHMAC with the given access token secret.
func NewAccessTokenManagerHMAC(accessTokenSecret string) *AccessTokenManagerHMAC {
	return &AccessTokenManagerHMAC{AccessTokenSecret: accessTokenSecret}
}

// New creates and signs a new access token with the given claims.
func (h *AccessTokenManagerHMAC) New(claims AccessTokenClaims) (string, error) {
	return NewHMAC(h.AccessTokenSecret, claims)
}

// Parse parses and validates the signature and claims of an access token.
func (h *AccessTokenManagerHMAC) Parse(token string) (AccessTokenClaims, error) {
	return ParseHMAC[AccessTokenClaims](h.AccessTokenSecret, token)
}
