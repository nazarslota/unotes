package oauth2

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/udholdenhed/unotes/auth/internal/domain/refreshtoken"
	"github.com/udholdenhed/unotes/auth/pkg/errors"
)

type RefreshRequest struct {
	RefreshToken string
}

type RefreshResult struct {
	AccessToken  string
	RefreshToken string
}

type RefreshRequestHandler interface {
	Handle(r RefreshRequest) (RefreshResult, error)
}

type refreshRequestHandler struct {
	AccessTokenSecret      string
	RefreshTokenSecret     string
	AccessTokenExpiresIn   time.Duration
	RefreshTokenExpiresIn  time.Duration
	RefreshTokenRepository refreshtoken.Repository
}

func NewRefreshRequestHandler(ats, rts string, ate, rte time.Duration, rtr refreshtoken.Repository) RefreshRequestHandler {
	return &refreshRequestHandler{
		AccessTokenSecret:      ats,
		RefreshTokenSecret:     rts,
		AccessTokenExpiresIn:   ate,
		RefreshTokenExpiresIn:  rte,
		RefreshTokenRepository: rtr,
	}
}

func (h *refreshRequestHandler) Handle(r RefreshRequest) (RefreshResult, error) {
	claims := make(jwt.MapClaims)
	if _, err := jwt.ParseWithClaims(r.RefreshToken, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidOrExpiredRefreshToken
		}
		return []byte(h.RefreshTokenSecret), nil
	}); err != nil {
		return RefreshResult{}, ErrInvalidOrExpiredRefreshToken.SetInternal(err)
	}

	var uid string
	if _, ok := claims["user_id"]; !ok {
		return RefreshResult{}, ErrInvalidOrExpiredRefreshToken
	} else if uid, ok = claims["user_id"].(string); !ok {
		return RefreshResult{}, ErrInvalidOrExpiredRefreshToken
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token, err := h.RefreshTokenRepository.GetRefreshToken(ctx, uid, refreshtoken.RefreshToken{Token: r.RefreshToken})
	if err != nil {
		return RefreshResult{}, ErrInvalidOrExpiredRefreshToken.SetInternal(err)
	} else if token == nil {
		return RefreshResult{}, ErrInvalidOrExpiredRefreshToken
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.RefreshTokenRepository.DeleteRefreshToken(ctx, uid, refreshtoken.RefreshToken{Token: r.RefreshToken}); err != nil {
		return RefreshResult{}, errors.ErrInternalServerError.SetInternal(err)
	}

	at, err := func() (string, error) {
		now := time.Now()
		claims := jwt.MapClaims{
			"exp": now.Add(h.AccessTokenExpiresIn).Unix(),
			"iat": now.Unix(),
			"nbf": now.Unix(),
		}
		claims["user_id"] = uid
		return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.AccessTokenSecret))
	}()
	if err != nil {
		return RefreshResult{}, errors.ErrInternalServerError.SetInternal(err)
	}

	rt, err := func() (string, error) {
		now := time.Now()
		claims := jwt.MapClaims{
			"exp": now.Add(h.RefreshTokenExpiresIn).Unix(),
			"iat": now.Unix(),
			"nbf": now.Unix(),
		}
		claims["user_id"] = uid
		return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.RefreshTokenSecret))
	}()
	if err != nil {
		return RefreshResult{}, errors.ErrInternalServerError.SetInternal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.RefreshTokenRepository.SetRefreshToken(ctx, uid, refreshtoken.RefreshToken{Token: rt}); err != nil {
		return RefreshResult{}, errors.ErrInternalServerError.SetInternal(err)
	}
	return RefreshResult{AccessToken: at, RefreshToken: rt}, nil
}
