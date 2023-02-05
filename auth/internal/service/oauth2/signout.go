package oauth2

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
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
	AccessTokenSecret      string
	RefreshTokenRepository domainrefresh.Repository
}

var ErrSignOutInvalidOrExpiredToken = errSignOutInvalidOrExpiredToken()

func errSignOutInvalidOrExpiredToken() error { return errors.New("invalid or expired token") }

func NewSignOutRequestHandler(accessTokenSecret string, refreshTokenRepository domainrefresh.Repository) LogOutRequestHandler {
	return &signOutRequestHandler{
		AccessTokenSecret:      accessTokenSecret,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

func (h signOutRequestHandler) Handle(ctx context.Context, request SignOutRequest) (SignOutResponse, error) {
	claims := make(jwt.MapClaims)
	if _, err := jwt.ParseWithClaims(request.AccessToken, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(h.AccessTokenSecret), nil
	}); err != nil {
		err = fmt.Errorf("failed to parse token: %w", err)
		return SignOutResponse{}, errors.Join(err, ErrSignOutInvalidOrExpiredToken)
	}

	var userID string
	if _, ok := claims["uid"]; !ok {
		return SignOutResponse{}, ErrSignOutInvalidOrExpiredToken
	} else if userID, ok = claims["uid"].(string); !ok {
		return SignOutResponse{}, ErrSignOutInvalidOrExpiredToken
	}

	tokens, err := h.RefreshTokenRepository.GetRefreshTokens(ctx, userID)
	if errors.Is(err, domainrefresh.ErrTokenNotFound) {
		return SignOutResponse{}, nil
	} else if err != nil {
		return SignOutResponse{}, fmt.Errorf("failed to get refresh tokens: %w", err)
	}

	if err := h.RefreshTokenRepository.DeleteRefreshTokens(ctx, userID, tokens); err != nil {
		return SignOutResponse{}, fmt.Errorf("failed to delete refresh tokens: %w", err)
	}
	return SignOutResponse{}, nil
}
