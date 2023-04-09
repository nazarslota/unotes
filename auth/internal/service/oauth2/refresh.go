package oauth2

import (
	"context"
	"errors"
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v4"
	domainrefresh "github.com/nazarslota/unotes/auth/internal/domain/refresh"
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
	AccessTokenCreator   AccessTokenCreator
	AccessTokenParser    AccessTokenParser
	AccessTokenExpiresIn time.Duration

	RefreshTokenCreator   RefreshTokenCreator
	RefreshTokenParser    RefreshTokenParser
	RefreshTokenExpiresIn time.Duration

	RefreshTokenSaver   RefreshTokenSaver
	RefreshTokenDeleter RefreshTokenDeleter
	RefreshTokenGetter  RefreshTokenGetter
}

var (
	ErrRefreshInvalidOrExpiredToken = errRefreshInvalidOrExpiredToken()
)

func errRefreshInvalidOrExpiredToken() error { return errors.New("invalid or expired token") }

func NewRefreshRequestHandler(
	accessTokenCreator AccessTokenCreator, accessTokenParser AccessTokenParser, accessTokenExpiresIn time.Duration,
	refreshTokenCreator RefreshTokenCreator, refreshTokenParser RefreshTokenParser, refreshTokenExpiresIn time.Duration,
	refreshTokenSaver RefreshTokenSaver, refreshTokenDeleter RefreshTokenDeleter, refreshTokenGetter RefreshTokenGetter,
) RefreshRequestHandler {
	return &refreshRequestHandler{
		AccessTokenCreator:   accessTokenCreator,
		AccessTokenParser:    accessTokenParser,
		AccessTokenExpiresIn: accessTokenExpiresIn,

		RefreshTokenCreator:   refreshTokenCreator,
		RefreshTokenParser:    refreshTokenParser,
		RefreshTokenExpiresIn: refreshTokenExpiresIn,

		RefreshTokenSaver:   refreshTokenSaver,
		RefreshTokenDeleter: refreshTokenDeleter,
		RefreshTokenGetter:  refreshTokenGetter,
	}
}

func (h refreshRequestHandler) Handle(ctx context.Context, request RefreshRequest) (RefreshResponse, error) {
	claims, err := h.RefreshTokenParser.Parse(request.RefreshToken)
	if err != nil {
		err = fmt.Errorf("failed to parse refresh accessToken: %w", err)
		return RefreshResponse{}, errors.Join(err, ErrRefreshInvalidOrExpiredToken)
	}

	tokens, err := h.RefreshTokenGetter.GetRefreshTokens(ctx, claims.UserID)
	if errors.Is(err, domainrefresh.ErrTokenNotFound) {
		err = fmt.Errorf("failed to get refresh tokens: %w", err)
		return RefreshResponse{}, errors.Join(err, ErrRefreshInvalidOrExpiredToken)
	} else if err != nil {
		return RefreshResponse{}, fmt.Errorf("failed to get refresh tokens: %w", err)
	} else if !slices.Contains(tokens, domainrefresh.Token(request.RefreshToken)) {
		return RefreshResponse{}, ErrRefreshInvalidOrExpiredToken
	}

	accessToken, err := h.AccessTokenCreator.New(jwt.AccessTokenClaims{
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(h.AccessTokenExpiresIn)),
		},
		UserID: claims.UserID,
	})
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("failed to create new access token: %w", err)
	}
	refreshToken, err := h.RefreshTokenCreator.New(jwt.RefreshTokenClaims{
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(h.RefreshTokenExpiresIn)),
		},
		UserID: claims.UserID,
	})
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("failed to create new refresh token: %w", err)
	}

	err = h.RefreshTokenDeleter.DeleteRefreshToken(ctx, claims.UserID, domainrefresh.Token(request.RefreshToken))
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("failed to delete refresh token: %w", err)
	}
	err = h.RefreshTokenSaver.SaveRefreshToken(ctx, claims.UserID, domainrefresh.Token(refreshToken))
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("failed to save refresh token: %w", err)
	}
	return RefreshResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
