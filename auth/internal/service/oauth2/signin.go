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
	"golang.org/x/crypto/bcrypt"
)

type SignInRequest struct {
	Username string
	Password string
}

type SignInResponse struct {
	AccessToken  string
	RefreshToken string
}

type SignInRequestHandler interface {
	Handle(ctx context.Context, request SignInRequest) (SignInResponse, error)
}

type signInRequestHandler struct {
	AccessTokenManager   AccessTokenManager[jwt.AccessTokenClaims]
	AccessTokenExpiresIn time.Duration

	RefreshTokenManager   RefreshTokenManager[jwt.RefreshTokenClaims]
	RefreshTokenExpiresIn time.Duration

	UserRepository         domainuser.Repository
	RefreshTokenRepository domainrefresh.Repository
}

var (
	ErrSignInInvalidUsername = errSignInInvalidUsername()
	ErrSignInInvalidPassword = errSignInInvalidPassword()
	ErrSignInUserNotFound    = errSignInUserNotFound()
)

func errSignInInvalidUsername() error { return errors.New("invalid username") }
func errSignInInvalidPassword() error { return errors.New("invalid password") }
func errSignInUserNotFound() error    { return domainuser.ErrUserNotFound }

func NewSignInRequestHandler(
	accessTokenManager AccessTokenManager[jwt.AccessTokenClaims], accessTokenExpiresIn time.Duration,
	refreshTokenManager RefreshTokenManager[jwt.RefreshTokenClaims], refreshTokenExpiresIn time.Duration,
	userRepository domainuser.Repository, refreshTokenRepository domainrefresh.Repository,
) SignInRequestHandler {
	return &signInRequestHandler{
		AccessTokenManager:     accessTokenManager,
		AccessTokenExpiresIn:   accessTokenExpiresIn,
		RefreshTokenManager:    refreshTokenManager,
		RefreshTokenExpiresIn:  refreshTokenExpiresIn,
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

func (h signInRequestHandler) Handle(ctx context.Context, request SignInRequest) (SignInResponse, error) {
	user, err := h.UserRepository.FindUserByUsername(ctx, request.Username)
	if err != nil {
		err = fmt.Errorf("failed to find user: %w", err)
		return SignInResponse{}, errors.Join(err, ErrSignInInvalidUsername)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		err = fmt.Errorf("failed to compare passwords: %w", err)
		return SignInResponse{}, errors.Join(err, ErrSignInInvalidPassword)
	} else if err != nil {
		return SignInResponse{}, fmt.Errorf("failed to compare passwords: %w", err)
	}

	accessToken, err := h.AccessTokenManager.New(jwt.AccessTokenClaims{
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(h.AccessTokenExpiresIn)),
		},
		UserID: user.ID,
	})
	if err != nil {
		return SignInResponse{}, fmt.Errorf("failed to create access token: %w", err)
	}
	refreshToken, err := h.RefreshTokenManager.New(jwt.RefreshTokenClaims{
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(h.RefreshTokenExpiresIn)),
		},
		UserID: user.ID,
	})
	if err != nil {
		return SignInResponse{}, fmt.Errorf("failed to create refresh token: %w", err)
	}

	err = h.RefreshTokenRepository.SaveRefreshToken(ctx, user.ID, domainrefresh.Token(refreshToken))
	if err != nil {
		return SignInResponse{}, fmt.Errorf("failed to save refresh token: %w", err)
	}
	return SignInResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
