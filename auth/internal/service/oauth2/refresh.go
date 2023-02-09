package oauth2

import (
	"context"
	"errors"
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v4"
	domainrefresh "github.com/nazarslota/unotes/auth/internal/domain/refresh"
	domainuser "github.com/nazarslota/unotes/auth/internal/domain/user"
	"github.com/nazarslota/unotes/auth/pkg/jwt"
	"golang.org/x/exp/slices"
)

type RefreshRequest struct {
	RefreshToken string
}

type RefreshResponse struct {
	AccessToken  string
	RefreshToken string
}

type RefreshRequestHandler interface {
	Handle(ctx context.Context, request RefreshRequest) (RefreshResponse, error)
}

type refreshRequestHandler struct {
	AccessTokenManager   AccessTokenManager[jwt.AccessTokenClaims]
	AccessTokenExpiresIn time.Duration

	RefreshTokenManager   RefreshTokenManager[jwt.RefreshTokenClaims]
	RefreshTokenExpiresIn time.Duration

	UserRepository         domainuser.Repository
	RefreshTokenRepository domainrefresh.Repository
}

var (
	ErrRefreshInvalidOrExpiredToken = errRefreshInvalidOrExpiredToken()
)

func errRefreshInvalidOrExpiredToken() error { return errors.New("invalid or expired token") }

func NewRefreshRequestHandler(
	accessTokenManager AccessTokenManager[jwt.AccessTokenClaims], accessTokenExpiresIn time.Duration,
	refreshTokenManager RefreshTokenManager[jwt.RefreshTokenClaims], refreshTokenExpiresIn time.Duration,
	userRepository domainuser.Repository, refreshTokenRepository domainrefresh.Repository,
) RefreshRequestHandler {
	return &refreshRequestHandler{
		AccessTokenManager:     accessTokenManager,
		AccessTokenExpiresIn:   accessTokenExpiresIn,
		RefreshTokenManager:    refreshTokenManager,
		RefreshTokenExpiresIn:  refreshTokenExpiresIn,
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

func (h refreshRequestHandler) Handle(ctx context.Context, request RefreshRequest) (RefreshResponse, error) {
	claims, err := h.RefreshTokenManager.Parse(request.RefreshToken)
	if err != nil {
		err = fmt.Errorf("failed to parse refresh accessToken: %w", err)
		return RefreshResponse{}, errors.Join(err, ErrRefreshInvalidOrExpiredToken)
	}

	tokens, err := h.RefreshTokenRepository.GetRefreshTokens(ctx, claims.UserID)
	if errors.Is(err, domainrefresh.ErrTokenNotFound) {
		err = fmt.Errorf("failed to get refresh tokens: %w", err)
		return RefreshResponse{}, errors.Join(err, ErrRefreshInvalidOrExpiredToken)
	} else if err != nil {
		return RefreshResponse{}, fmt.Errorf("failed to get refresh tokens: %w", err)
	} else if !slices.Contains(tokens, domainrefresh.Token(request.RefreshToken)) {
		return RefreshResponse{}, ErrRefreshInvalidOrExpiredToken
	}

	accessToken, err := h.AccessTokenManager.New(jwt.AccessTokenClaims{
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(h.AccessTokenExpiresIn)),
		},
		UserID: claims.UserID,
	})
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("failed to create new access token: %w", err)
	}
	refreshToken, err := h.RefreshTokenManager.New(jwt.RefreshTokenClaims{
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(h.RefreshTokenExpiresIn)),
		},
		UserID: claims.UserID,
	})
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("failed to create new refresh token: %w", err)
	}

	err = h.RefreshTokenRepository.DeleteRefreshToken(ctx, claims.UserID, domainrefresh.Token(request.RefreshToken))
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("failed to delete refresh token: %w", err)
	}
	err = h.RefreshTokenRepository.SaveRefreshToken(ctx, claims.UserID, domainrefresh.Token(refreshToken))
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("failed to save refresh token: %w", err)
	}
	return RefreshResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
