package oauth2

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/udholdenhed/unotes/auth-service/internal/domain/refreshtoken"
	"github.com/udholdenhed/unotes/auth-service/pkg/errors"
)

type RefreshRequest struct {
	RefreshToken string
}

type RefreshResult struct {
	AccessToken  string
	RefreshToken string
}

type RefreshRequestHandler interface {
	Handle(ctx context.Context, r RefreshRequest) (RefreshResult, error)
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

func (h *refreshRequestHandler) Handle(ctx context.Context, r RefreshRequest) (RefreshResult, error) {
	claims := make(jwt.MapClaims)
	if _, err := jwt.ParseWithClaims(r.RefreshToken, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidOrExpiredToken
		}

		return []byte(h.RefreshTokenSecret), nil
	}); err != nil {
		return RefreshResult{}, ErrInvalidOrExpiredToken.SetInternal(err)
	}

	var userID string
	if _, ok := claims["user_id"]; !ok {
		return RefreshResult{}, ErrInvalidOrExpiredToken
	} else if userID, ok = claims["user_id"].(string); !ok {
		return RefreshResult{}, ErrInvalidOrExpiredToken
	}

	token, err := h.RefreshTokenRepository.GetRefreshToken(ctx, userID, refreshtoken.Token{Token: r.RefreshToken})
	if err != nil {
		return RefreshResult{}, ErrInvalidOrExpiredToken.SetInternal(err)
	} else if token == nil {
		return RefreshResult{}, ErrInvalidOrExpiredToken
	}

	if err := h.RefreshTokenRepository.DeleteRefreshToken(ctx, userID, refreshtoken.Token{Token: r.RefreshToken}); err != nil {
		return RefreshResult{}, errors.ErrInternalServerError.SetInternal(err)
	}

	accessToken, err := func() (string, error) {
		now := time.Now()
		claims := jwt.MapClaims{
			"exp": now.Add(h.AccessTokenExpiresIn).Unix(),
			"iat": now.Unix(),
			"nbf": now.Unix(),
		}
		claims["user_id"] = userID
		return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.AccessTokenSecret))
	}()
	if err != nil {
		return RefreshResult{}, errors.ErrInternalServerError.SetInternal(err)
	}

	refreshToken, err := func() (string, error) {
		now := time.Now()
		claims := jwt.MapClaims{
			"exp": now.Add(h.RefreshTokenExpiresIn).Unix(),
			"iat": now.Unix(),
			"nbf": now.Unix(),
		}
		claims["user_id"] = userID
		return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.RefreshTokenSecret))
	}()
	if err != nil {
		return RefreshResult{}, errors.ErrInternalServerError.SetInternal(err)
	}

	if err := h.RefreshTokenRepository.SetRefreshToken(ctx, userID, refreshtoken.Token{Token: refreshToken}); err != nil {
		return RefreshResult{}, errors.ErrInternalServerError.SetInternal(err)
	}

	return RefreshResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
