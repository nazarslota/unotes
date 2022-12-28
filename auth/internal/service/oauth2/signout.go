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

func NewSignOutRequestHandler(
	accessTokenSecret string, refreshTokenRepository refreshtoken.Repository,
) LogOutRequestHandler {
	return &signOutRequestHandler{
		AccessTokenSecret:      accessTokenSecret,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

func (h *signOutRequestHandler) Handle(ctx context.Context, request *SignOutRequest) (*SignOutResponse, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context is done: %w", ctx.Err()) // failed to handle sign out request,
	default:
	}

	claims, err := tokens{}.ParseHS256(request.AccessToken, h.AccessTokenSecret)
	if err != nil {
		return nil, err
	}

	userID := ""
	if _, ok := claims["user_id"]; !ok {
		return nil, fmt.Errorf("filed to get user id from token")
	} else if userID, ok = claims["user_id"].(string); !ok {
		return nil, fmt.Errorf("failed to convert user id to string")
	}

	tokens, err := h.RefreshTokenRepository.FindMany(ctx, userID)
	if err != nil {
		if errors.Is(err, refreshtoken.ErrTokensNotFound) {
			return nil, fmt.Errorf("failed to find user tokens: %w", ErrInvalidOrExpiredToken)
		}
		return nil, fmt.Errorf("failed to find user tokens: %w", err)
	}

	for _, token := range tokens {
		err := h.RefreshTokenRepository.DeleteOne(ctx, userID, &token)
		if err != nil {
			return nil, fmt.Errorf("failed to delete refresh token: %w", err)
		}
	}
	return &SignOutResponse{}, nil
}
