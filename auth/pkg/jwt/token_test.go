package jwt

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewHMAC(t *testing.T) {
	tests := []struct {
		name      string
		secret    string
		claims    jwt.MapClaims
		wantToken string
		wantErr   error
	}{
		{
			name:      "should create new token",
			secret:    "secret",
			claims:    jwt.MapClaims{"sub": "1234567890", "name": "John Doe", "iat": 1516239022},
			wantToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjIsIm5hbWUiOiJKb2huIERvZSIsInN1YiI6IjEyMzQ1Njc4OTAifQ.ub7srKZNrlkC9jpqvPSYMwZp8IZQN1ZBCuld49qCqOs",
			wantErr:   nil,
		},
		{
			name:      "should return an error if secret is empty",
			secret:    "",
			claims:    jwt.MapClaims{"sub": "1234567890", "name": "John Doe", "iat": 1516239022},
			wantToken: "",
			wantErr:   jwt.ErrInvalidKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := NewHMAC[jwt.MapClaims](tt.secret, tt.claims)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantToken, token)
		})
	}
}

func TestParseHMAC(t *testing.T) {
	tests := []struct {
		name       string
		secret     string
		token      string
		wantClaims jwt.MapClaims
		wantErr    error
	}{
		{
			name:       "should parse valid token",
			secret:     "secret",
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJzdWIiOiIxMjM0NTY3ODkwIn0.LmURJCzy7NfJBggkWYrAZ8XNzFxMhrDgBJUmBcffhSw",
			wantClaims: jwt.MapClaims{"sub": "1234567890", "name": "John Doe"},
			wantErr:    nil,
		},
		{
			name:       "should return an error if secret is invalid",
			secret:     "invalid",
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJzdWIiOiIxMjM0NTY3ODkwIn0.LmURJCzy7NfJBggkWYrAZ8XNzFxMhrDgBJUmBcffhSw",
			wantClaims: nil,
			wantErr:    jwt.ErrSignatureInvalid,
		},
		{
			name:       "should return an error if token is expired",
			secret:     "secret",
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyMzkwMjIsIm5hbWUiOiJKb2huIERvZSIsInN1YiI6IjEyMzQ1Njc4OTAifQ.E9lKgKUVuCmwIylcq5bFBBQ67fWvKy1HX15qZpVDEiU",
			wantClaims: nil,
			wantErr:    jwt.ErrTokenExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseHMAC[jwt.MapClaims](tt.secret, tt.token)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantClaims, claims)
		})
	}
}
