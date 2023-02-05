package oauth2

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	domainrefresh "github.com/nazarslota/unotes/auth/internal/domain/refresh"
	domainuser "github.com/nazarslota/unotes/auth/internal/domain/user"
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
	AccessTokenSecret      string
	RefreshTokenSecret     string
	AccessTokenExpiresIn   time.Duration
	RefreshTokenExpiresIn  time.Duration
	UserRepository         domainuser.Repository
	RefreshTokenRepository domainrefresh.Repository
}

var (
	ErrRefreshInvalidOrExpiredToken = errRefreshInvalidOrExpiredToken()
	ErrRefreshUserNotFound          = errRefreshUserNotFound()
)

func errRefreshInvalidOrExpiredToken() error { return errors.New("invalid or expired token") }
func errRefreshUserNotFound() error          { return domainuser.ErrUserNotFound }

func NewRefreshRequestHandler(
	accessTokenSecret, refreshTokenSecret string,
	accessTokenExpiresIn, refreshTokenExpiresIn time.Duration,
	userRepository domainuser.Repository, refreshTokenRepository domainrefresh.Repository,
) RefreshRequestHandler {
	return &refreshRequestHandler{
		AccessTokenSecret:      accessTokenSecret,
		RefreshTokenSecret:     refreshTokenSecret,
		AccessTokenExpiresIn:   accessTokenExpiresIn,
		RefreshTokenExpiresIn:  refreshTokenExpiresIn,
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

func (h refreshRequestHandler) Handle(ctx context.Context, request RefreshRequest) (RefreshResponse, error) {
	userID, err := h.parseRefreshToken(request.RefreshToken)
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("failed parse refresh token: %w", err)
	}
	if err := h.validateRefreshToken(ctx, userID, request.RefreshToken); err != nil {
		return RefreshResponse{}, fmt.Errorf("failed to validate refresh token: %w", err)
	}

	user, err := h.UserRepository.FindUserByUserID(ctx, userID)
	if errors.Is(err, domainuser.ErrUserNotFound) {
		err = fmt.Errorf("failed to find user: %w", err)
		return RefreshResponse{}, errors.Join(err, ErrRefreshInvalidOrExpiredToken)
	} else if err != nil {
		return RefreshResponse{}, fmt.Errorf("failed to find user: %w", err)
	}

	accessToken, err := h.newAccessToken(userID, user.Username)
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("failed to create access token: %w", err)
	}
	refreshToken, err := h.newRefreshToken(userID, user.Username)
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("failed to create refresh token: %w", err)
	}

	err = h.RefreshTokenRepository.DeleteRefreshToken(ctx, userID, domainrefresh.Token(request.RefreshToken))
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("failed to delete refresh token: %w", err)
	}
	err = h.RefreshTokenRepository.SaveRefreshToken(ctx, userID, domainrefresh.Token(refreshToken))
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("failed to save refresh token: %w", err)
	}
	return RefreshResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (h refreshRequestHandler) parseRefreshToken(token string) (userID string, err error) {
	claims := make(jwt.MapClaims)
	if _, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(h.RefreshTokenSecret), nil
	}); err != nil {
		err = fmt.Errorf("failed to parse token: %w", err)
		return "", errors.Join(err, ErrRefreshInvalidOrExpiredToken)
	}

	if _, ok := claims["uid"]; !ok {
		return "", ErrRefreshInvalidOrExpiredToken
	} else if userID, ok = claims["uid"].(string); !ok {
		return "", ErrRefreshInvalidOrExpiredToken
	}
	return userID, nil
}

func (h refreshRequestHandler) validateRefreshToken(ctx context.Context, userID string, token string) error {
	tokens, err := h.RefreshTokenRepository.GetRefreshTokens(ctx, userID)
	if errors.Is(err, domainrefresh.ErrTokenNotFound) {
		return ErrRefreshInvalidOrExpiredToken
	} else if err != nil {
		return fmt.Errorf("failed to get refresh tokens: %w", err)
	}

	if !slices.Contains(tokens, domainrefresh.Token(token)) {
		return ErrRefreshInvalidOrExpiredToken
	}
	return nil
}

func (h refreshRequestHandler) newAccessToken(userID, username string) (string, error) {
	return h.newToken(userID, username, []byte(h.AccessTokenSecret), h.AccessTokenExpiresIn)
}

func (h refreshRequestHandler) newRefreshToken(userID, username string) (string, error) {
	return h.newToken(userID, username, []byte(h.RefreshTokenSecret), h.RefreshTokenExpiresIn)
}

func (h refreshRequestHandler) newToken(userID, username string, secret []byte, expiresIn time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"iss": username,
		"sub": "auth",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(expiresIn).Unix(),
		"uid": userID,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}
	return token, nil
}
