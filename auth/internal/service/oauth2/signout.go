package oauth2

import (
	"context"
	"errors"
	"fmt"

	domainrefresh "github.com/nazarslota/unotes/auth/internal/domain/refresh"
)

type SignOutRequest struct {
	AccessToken string
}

type SignOutResponse struct{}

type LogOutRequestHandler interface {
	Handle(ctx context.Context, request SignOutRequest) (SignOutResponse, error)
}

type signOutRequestHandler struct {
	AccessTokenParser AccessTokenParser

	RefreshTokenDeleter RefreshTokensDeleter
	RefreshTokenGetter  RefreshTokenGetter
}

var ErrSignOutInvalidOrExpiredToken = errSignOutInvalidOrExpiredToken()

func errSignOutInvalidOrExpiredToken() error { return errors.New("invalid or expired token") }

func NewSignOutRequestHandler(
	accessTokenParser AccessTokenParser,
	refreshTokenDeleter RefreshTokensDeleter, refreshTokenGetter RefreshTokenGetter,
) LogOutRequestHandler {
	return &signOutRequestHandler{
		AccessTokenParser: accessTokenParser,

		RefreshTokenDeleter: refreshTokenDeleter,
		RefreshTokenGetter:  refreshTokenGetter,
	}
}

func (h signOutRequestHandler) Handle(ctx context.Context, request SignOutRequest) (SignOutResponse, error) {
	claims, err := h.AccessTokenParser.Parse(request.AccessToken)
	if err != nil {
		err = fmt.Errorf("failed to parse access token: %w", err)
		return SignOutResponse{}, errors.Join(err, ErrSignOutInvalidOrExpiredToken)
	}

	tokens, err := h.RefreshTokenGetter.GetRefreshTokens(ctx, claims.UserID)
	if errors.Is(err, domainrefresh.ErrTokenNotFound) {
		return SignOutResponse{}, nil
	} else if err != nil {
		return SignOutResponse{}, fmt.Errorf("failed to get refresh tokens: %w", err)
	}

	if err := h.RefreshTokenDeleter.DeleteRefreshTokens(ctx, claims.UserID, tokens); err != nil {
		return SignOutResponse{}, fmt.Errorf("failed to delete refresh tokens: %w", err)
	}
	return SignOutResponse{}, nil
}
