package oauth2

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/udholdenhed/unotes/auth/internal/domain/refreshtoken"
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
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context is done: %w", ctx.Err()) // failed to handle refresh request,
	default:
	}

	claims, err := tokens{}.ParseHS256(request.RefreshToken, h.RefreshTokenSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the refresh token: %w", ErrInvalidOrExpiredToken)
	}

	userID := ""
	if _, ok := claims["user_id"]; !ok {
		return nil, fmt.Errorf("filed to get user id from token: %w", ErrInvalidOrExpiredToken)
	} else if userID, ok = claims["user_id"].(string); !ok {
		return nil, fmt.Errorf("failed to convert user id to string: %w", ErrInvalidOrExpiredToken)
	}

	token := &refreshtoken.Token{Token: request.RefreshToken}
	_, err = h.RefreshTokenRepository.FindOne(ctx, userID, token)
	if err != nil {
		if errors.Is(err, refreshtoken.ErrTokenNotFound) {
			return nil, fmt.Errorf("the refresh token could not be found: %w", ErrInvalidOrExpiredToken)
		}
		return nil, fmt.Errorf("the refresh token could not be found: %w", err)
	}

	err = h.RefreshTokenRepository.DeleteOne(ctx, userID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to delete refresh token: %w", err)
	}

	response := new(RefreshResponse)
	access, err := tokens{}.NewHS256(h.AccessTokenSecret, h.AccessTokenExpiresIn, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create an access token: %w", err)
	}
	response.AccessToken = access

	refresh, err := tokens{}.NewHS256(h.RefreshTokenSecret, h.RefreshTokenExpiresIn, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create a refresh token: %w", err)
	}
	response.RefreshToken = refresh

	err = h.RefreshTokenRepository.SaveOne(ctx, userID, &refreshtoken.Token{Token: refresh})
	if err != nil {
		return nil, fmt.Errorf("failed to save the refresh token: %w", err)
	}
	return response, nil
}
