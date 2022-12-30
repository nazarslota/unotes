package oauth2

import (
	"context"
	"errors"
	"fmt"

	"github.com/udholdenhed/unotes/auth/internal/domain/refreshtoken"
)

type SignOutRequest struct {
	AccessToken string `json:"access_token"`
}

type SignOutResponse struct {
}

type LogOutRequestHandler interface {
	Handle(ctx context.Context, request *SignOutRequest) (*SignOutResponse, error)
}

type signOutRequestHandler struct {
	AccessTokenSecret      string
	RefreshTokenRepository refreshtoken.Repository
}

var (
	ErrSignOutInvalidOrExpiredToken = errors.New("invalid or expired token")
)

func NewSignOutRequestHandler(
	accessTokenSecret string, refreshTokenRepository refreshtoken.Repository,
) LogOutRequestHandler {
	return &signOutRequestHandler{
		AccessTokenSecret:      accessTokenSecret,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

func (h *signOutRequestHandler) Handle(ctx context.Context, request *SignOutRequest) (*SignOutResponse, error) {
	claims, err := parseHS256(request.AccessToken, h.AccessTokenSecret)
	if errors.Is(err, ErrJWTInvalidOrExpiredToken) {
		return nil, fmt.Errorf("failed to parse the access token: %w", ErrSignOutInvalidOrExpiredToken)
	} else if err != nil {
		return nil, fmt.Errorf("failed to parse the access token: %w", err)
	}

	userID := ""
	if _, ok := claims["user_id"]; !ok {
		return nil, fmt.Errorf("failed to get the user ID from the token: %w", ErrSignOutInvalidOrExpiredToken)
	} else if userID, ok = claims["user_id"].(string); !ok {
		return nil, fmt.Errorf("failed to convert user ID to string: %w", ErrSignOutInvalidOrExpiredToken)
	}

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
