package oauth2

import (
	"context"

	"github.com/golang-jwt/jwt"
	"github.com/udholdenhed/unotes/auth-service/internal/domain/refreshtoken"
	"github.com/udholdenhed/unotes/auth-service/pkg/errors"
)

type LogOutRequest struct {
	AccessToken string `json:"access_token"`
}

type LogOutRequestHandler interface {
	Handle(ctx context.Context, r LogOutRequest) error
}

type signOutRequestHandler struct {
	AccessTokenSecret      string
	RefreshTokenRepository refreshtoken.Repository
}

func NewSignOutRequestHandler(accessTokenSecret string, refreshTokenRepository refreshtoken.Repository) LogOutRequestHandler {
	return &signOutRequestHandler{
		AccessTokenSecret:      accessTokenSecret,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

func (h *signOutRequestHandler) Handle(ctx context.Context, r LogOutRequest) error {
	claims := make(jwt.MapClaims)
	if _, err := jwt.ParseWithClaims(r.AccessToken, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidOrExpiredToken
		}

		return []byte(h.AccessTokenSecret), nil
	}); err != nil {
		return ErrInvalidOrExpiredToken.SetInternal(err)
	}

	var userID string
	if _, ok := claims["user_id"]; !ok {
		return ErrInvalidOrExpiredToken
	} else if userID, ok = claims["user_id"].(string); !ok {
		return ErrInvalidOrExpiredToken
	}

	allRefreshTokens, err := h.RefreshTokenRepository.GetAllRefreshTokens(ctx, userID)
	if err != nil {
		return errors.ErrInternalServerError.SetInternal(err)
	}

	for _, token := range allRefreshTokens {
		if err := h.RefreshTokenRepository.DeleteRefreshToken(ctx, userID, token); err != nil {
			return errors.ErrInternalServerError.SetInternal(err)
		}
	}
	return nil
}
