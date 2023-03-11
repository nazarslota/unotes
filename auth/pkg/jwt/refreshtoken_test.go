package jwt

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRefreshTokenManagerHMAC(t *testing.T) {
	tm := NewRefreshTokenManagerHMAC("secret")
	assert.NotNil(t, tm)
	assert.Equal(t, "secret", tm.RefreshTokenSecret)
}

func TestRefreshTokenManagerHMAC_New(t *testing.T) {
	t.Run("should create a new refresh token", func(t *testing.T) {
		tm := NewRefreshTokenManagerHMAC("secret")
		require.NotNil(t, tm)

		token, err := tm.New(RefreshTokenClaims{UserID: "e10adb24-7179-468f-911d-cc90aacb7410"})
		assert.NoError(t, err)
		assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZTEwYWRiMjQtNzE3OS00NjhmLTkxMWQtY2M5MGFhY2I3NDEwIn0.DFyFMDoWV4HyRhA3DrNWEqIuk-dvpICKa6VhSMWdK_w", token)
	})

	t.Run("should return an error if secret is empty", func(t *testing.T) {
		tm := NewRefreshTokenManagerHMAC("")
		require.NotNil(t, tm)

		token, err := tm.New(RefreshTokenClaims{UserID: "e10adb24-7179-468f-911d-cc90aacb7410"})
		assert.ErrorIs(t, err, jwt.ErrInvalidKey)
		assert.Empty(t, token)
	})
}

func TestRefreshTokenManagerHMAC_Parse(t *testing.T) {
	t.Run("should parse refresh token", func(t *testing.T) {
		tm := NewRefreshTokenManagerHMAC("secret")
		require.NotNil(t, tm)

		claims, err := tm.Parse("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZTEwYWRiMjQtNzE3OS00NjhmLTkxMWQtY2M5MGFhY2I3NDEwIn0.DFyFMDoWV4HyRhA3DrNWEqIuk-dvpICKa6VhSMWdK_w")
		assert.NoError(t, err)
		assert.Equal(t, RefreshTokenClaims{UserID: "e10adb24-7179-468f-911d-cc90aacb7410"}, claims)
	})

	t.Run("should return an error if token is expired", func(t *testing.T) {
		tm := NewRefreshTokenManagerHMAC("secret")
		require.NotNil(t, tm)

		claims, err := tm.Parse("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyMzkwMjIsInVzZXJfaWQiOiJlMTBhZGIyNC03MTc5LTQ2OGYtOTExZC1jYzkwYWFjYjc0MTAifQ.9QDLYKxLIbqKeUkH-OR6-292RaQAwV1XbD9i7QR7p20")
		assert.ErrorIs(t, err, jwt.ErrTokenExpired)
		assert.Empty(t, claims)
	})
}
