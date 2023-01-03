package oauth2

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nazarslota/unotes/auth/internal/domain/refreshtoken"
)

type RefreshRequest struct {
	RefreshToken string
}

type RefreshResponse struct {
	AccessToken  string
	RefreshToken string
}

type RefreshRequestHandler interface {
	Handle(ctx context.Context, request *RefreshRequest) (*RefreshResponse, error)
}

type refreshRequestHandler struct {
	AccessTokenSecret      string
	RefreshTokenSecret     string
	AccessTokenExpiresIn   time.Duration
	RefreshTokenExpiresIn  time.Duration
	RefreshTokenRepository refreshtoken.Repository
}

var (
	ErrRefreshInvalidOrExpiredToken = errors.New("invalid or expired token")
)

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

func (h *refreshRequestHandler) Handle(ctx context.Context, request *RefreshRequest) (*RefreshResponse, error) {
	claims, err := parseHS256(request.RefreshToken, h.RefreshTokenSecret)
	if errors.Is(err, ErrJWTInvalidOrExpiredToken) {
		return nil, fmt.Errorf("failed to parse the access token: %w", ErrRefreshInvalidOrExpiredToken)
	} else if err != nil {
		return nil, fmt.Errorf("failed to parse the access token: %w", err)
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
