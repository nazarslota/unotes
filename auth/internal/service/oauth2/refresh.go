package oauth2

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nazarslota/unotes/auth/internal/domain/refreshtoken"
)

// RefreshRequest represents a refresh request.
type RefreshRequest struct {
	RefreshToken string
}

// RefreshResponse represents a refresh response.
type RefreshResponse struct {
	AccessToken  string
	RefreshToken string
}

// RefreshRequestHandler is an interface that defines a refresh request handler.
type RefreshRequestHandler interface {
	Handle(ctx context.Context, request *RefreshRequest) (*RefreshResponse, error)
}

// refreshRequestHandler is a refresh request handler that refreshes the access token and refresh token.
type refreshRequestHandler struct {
	AccessTokenSecret      string
	RefreshTokenSecret     string
	AccessTokenExpiresIn   time.Duration
	RefreshTokenExpiresIn  time.Duration
	RefreshTokenRepository refreshtoken.Repository
}

// ErrRefreshInvalidOrExpiredToken is returned when the refresh token is invalid or expired.
var ErrRefreshInvalidOrExpiredToken = errors.New("invalid or expired token")

// NewRefreshRequestHandler returns a new refresh request handler.
func NewRefreshRequestHandler(
	accessTokenSecret, refreshTokenSecret string,
	accessTokenExpiresIn, refreshTokenExpiresIn time.Duration,
	refreshTokenRepository refreshtoken.Repository,
) RefreshRequestHandler {
	return &refreshRequestHandler{
		AccessTokenSecret:      accessTokenSecret,
		RefreshTokenSecret:     refreshTokenSecret,
		AccessTokenExpiresIn:   accessTokenExpiresIn,
		RefreshTokenExpiresIn:  refreshTokenExpiresIn,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

// Handle handles a refresh request and returns a response.
//
// It can return the following errors:
//   - ErrRefreshInvalidOrExpiredToken: if the refresh token is invalid or expired
//   - other errors: if an error occurred while parsing the refresh token, deleting the refresh token, creating the access or refresh tokens, or saving the refresh token
func (h *refreshRequestHandler) Handle(ctx context.Context, request *RefreshRequest) (*RefreshResponse, error) {
	claims, err := parseHS256(request.RefreshToken, h.RefreshTokenSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the access token: %w", ErrRefreshInvalidOrExpiredToken)
	}

	userID := ""
	if _, ok := claims["user_id"]; !ok {
		return nil, fmt.Errorf("failed to get the user ID from the token: %w", ErrRefreshInvalidOrExpiredToken)
	} else if userID, ok = claims["user_id"].(string); !ok {
		return nil, fmt.Errorf("failed to convert user ID to string: %w", ErrRefreshInvalidOrExpiredToken)
	}

	token := &refreshtoken.Token{Token: request.RefreshToken}

	_, err = h.RefreshTokenRepository.FindOne(ctx, userID, token)
	if errors.Is(err, refreshtoken.ErrTokenNotFound) {
		return nil, fmt.Errorf("there is no such refreshtoken: %w", ErrRefreshInvalidOrExpiredToken)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get the refresh token from the repository: %w", err)
	}

	if err := h.RefreshTokenRepository.DeleteOne(ctx, userID, token); err != nil {
		return nil, fmt.Errorf("failed to delete the refresh token: %w", err)
	}

	response := new(RefreshResponse)
	if response.AccessToken, err = newHS256(h.AccessTokenSecret, h.AccessTokenExpiresIn, userID); err != nil {
		return nil, fmt.Errorf("failed to create an access token: %w", err)
	}

	if response.RefreshToken, err = newHS256(h.RefreshTokenSecret, h.RefreshTokenExpiresIn, userID); err != nil {
		return nil, fmt.Errorf("failed to create a refresh token: %w", err)
	}

	token = &refreshtoken.Token{Token: response.RefreshToken}
	if err := h.RefreshTokenRepository.SaveOne(ctx, userID, token); err != nil {
		return nil, fmt.Errorf("failed to save the refresh token: %w", err)
	}
	return response, nil
}
