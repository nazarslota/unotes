package oauth2

import (
	"context"
	"errors"
	"fmt"

	"github.com/nazarslota/unotes/auth/internal/domain/refreshtoken"
)

// SignOutRequest represents a sign out request.
type SignOutRequest struct {
	AccessToken string `json:"access_token"`
}

// SignOutResponse represents a sign out response.
type SignOutResponse struct {
}

// LogOutRequestHandler is an interface that defines a sign out request handler.
type LogOutRequestHandler interface {
	Handle(ctx context.Context, request *SignOutRequest) (*SignOutResponse, error)
}

// signOutRequestHandler is a sign out request handler that deletes the user's refresh tokens.
type signOutRequestHandler struct {
	AccessTokenSecret      string
	RefreshTokenRepository refreshtoken.Repository
}

// ErrSignOutInvalidOrExpiredToken is returned when the access token is invalid or expired.
var ErrSignOutInvalidOrExpiredToken = errors.New("invalid or expired token")

// NewSignOutRequestHandler creates a new sign out request handler.
func NewSignOutRequestHandler(
	accessTokenSecret string, refreshTokenRepository refreshtoken.Repository,
) LogOutRequestHandler {
	return &signOutRequestHandler{
		AccessTokenSecret:      accessTokenSecret,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

// Handle handles a sign-out request and returns a response.
//
// It can return the following errors:
//   - ErrSignOutInvalidOrExpiredToken: if the access token is invalid or expired
//   - other errors: if an error occurred while parsing the access token or deleting the refresh tokens
func (h *signOutRequestHandler) Handle(ctx context.Context, request *SignOutRequest) (*SignOutResponse, error) {
	// Parse the access token to get the user ID.
	claims, err := parseHS256(request.AccessToken, h.AccessTokenSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the access token: %w", ErrSignOutInvalidOrExpiredToken)
	}

	userID := ""
	if _, ok := claims["user_id"]; !ok {
		return nil, fmt.Errorf("failed to get the user ID from the token: %w", ErrSignOutInvalidOrExpiredToken)
	} else if userID, ok = claims["user_id"].(string); !ok {
		return nil, fmt.Errorf("failed to convert user ID to string: %w", ErrSignOutInvalidOrExpiredToken)
	}

	// Get the user's refresh tokens.
	tokens, err := h.RefreshTokenRepository.FindMany(ctx, userID)
	if errors.Is(err, refreshtoken.ErrTokensNotFound) {
		return &SignOutResponse{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to search for the user's refresh tokens: %w", err)
	}

	for _, token := range tokens {
		if err := h.RefreshTokenRepository.DeleteOne(ctx, userID, &token); err != nil {
			return nil, fmt.Errorf("failed to delete the refresh token: %w", err)
		}
	}
	return &SignOutResponse{}, nil
}
