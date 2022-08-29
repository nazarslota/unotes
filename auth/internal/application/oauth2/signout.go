package oauth2

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/udholdenhed/unotes/auth/internal/domain/refreshtoken"
	"github.com/udholdenhed/unotes/auth/pkg/errors"
)

type LogOutRequest struct {
	AccessToken string `json:"access_token"`
}

type LogOutRequestHandler interface {
	Handle(r LogOutRequest) error
}

type signOutRequestHandler struct {
	AccessTokenSecret      string
	RefreshTokenRepository refreshtoken.Repository
}

func NewSignOutRequestHandler(ats string, rtr refreshtoken.Repository) LogOutRequestHandler {
	return &signOutRequestHandler{
		AccessTokenSecret:      ats,
		RefreshTokenRepository: rtr,
	}
}

func (h *signOutRequestHandler) Handle(r LogOutRequest) error {
	claims := make(jwt.MapClaims)
	if _, err := jwt.ParseWithClaims(r.AccessToken, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidOrExpiredRefreshToken
		}
		return []byte(h.AccessTokenSecret), nil
	}); err != nil {
		return ErrInvalidOrExpiredRefreshToken.SetInternal(err)
	}

	var uid string
	if _, ok := claims["user_id"]; !ok {
		return ErrInvalidOrExpiredRefreshToken
	} else if uid, ok = claims["user_id"].(string); !ok {
		return ErrInvalidOrExpiredRefreshToken
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tokens, err := h.RefreshTokenRepository.GetAllRefreshTokens(ctx, uid)
	if err != nil {
		return errors.ErrInternalServerError.SetInternal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, token := range tokens {
		if err := h.RefreshTokenRepository.DeleteRefreshToken(ctx, uid, token); err != nil {
			return errors.ErrInternalServerError.SetInternal(err)
		}
	}
	return nil
}
