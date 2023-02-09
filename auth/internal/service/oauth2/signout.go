package oauth2

import (
	"context"
	"errors"
	"fmt"

	domainrefresh "github.com/nazarslota/unotes/auth/internal/domain/refresh"
	"github.com/nazarslota/unotes/auth/pkg/jwt"
)

type SignOutRequest struct {
	AccessToken string
}

type SignOutResponse struct{}

type LogOutRequestHandler interface {
	Handle(ctx context.Context, request SignOutRequest) (SignOutResponse, error)
}

type signOutRequestHandler struct {
	AccessTokenManager     AccessTokenManager[jwt.AccessTokenClaims]
	RefreshTokenRepository domainrefresh.Repository
}

var ErrSignOutInvalidOrExpiredToken = errSignOutInvalidOrExpiredToken()

func errSignOutInvalidOrExpiredToken() error { return errors.New("invalid or expired token") }

func NewSignOutRequestHandler(accessTokenManager AccessTokenManager[jwt.AccessTokenClaims], refreshTokenRepository domainrefresh.Repository) LogOutRequestHandler {
	return &signOutRequestHandler{
		AccessTokenManager:     accessTokenManager,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

func (h signOutRequestHandler) Handle(ctx context.Context, request SignOutRequest) (SignOutResponse, error) {
	claims, err := h.AccessTokenManager.Parse(request.AccessToken)
	if err != nil {
		err = fmt.Errorf("failed to parse access token: %w", err)
		return SignOutResponse{}, errors.Join(err, ErrSignOutInvalidOrExpiredToken)
	}

	tokens, err := h.RefreshTokenRepository.GetRefreshTokens(ctx, claims.UserID)
	if errors.Is(err, domainrefresh.ErrTokenNotFound) {
		return SignOutResponse{}, nil
	} else if err != nil {
		return SignOutResponse{}, fmt.Errorf("failed to get refresh tokens: %w", err)
	}

	if err := h.RefreshTokenRepository.DeleteRefreshTokens(ctx, claims.UserID, tokens); err != nil {
		return SignOutResponse{}, fmt.Errorf("failed to delete refresh tokens: %w", err)
	}
	return SignOutResponse{}, nil
}
