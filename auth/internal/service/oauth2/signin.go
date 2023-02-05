package oauth2

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	domainrefresh "github.com/nazarslota/unotes/auth/internal/domain/refresh"
	domainuser "github.com/nazarslota/unotes/auth/internal/domain/user"
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
	AccessTokenSecret      string
	RefreshTokenSecret     string
	AccessTokenExpiresIn   time.Duration
	RefreshTokenExpiresIn  time.Duration
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
	accessTokenSecret, refreshTokenSecret string,
	accessTokenExpiresIn, refreshTokenExpiresIn time.Duration,
	userRepository domainuser.Repository, refreshTokenRepository domainrefresh.Repository,
) SignInRequestHandler {
	return &signInRequestHandler{
		AccessTokenSecret:      accessTokenSecret,
		RefreshTokenSecret:     refreshTokenSecret,
		AccessTokenExpiresIn:   accessTokenExpiresIn,
		RefreshTokenExpiresIn:  refreshTokenExpiresIn,
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

func (h signInRequestHandler) Handle(ctx context.Context, request SignInRequest) (SignInResponse, error) {
	if len(request.Username) == 0 {
		return SignInResponse{}, ErrSignInInvalidUsername
	} else if len(request.Password) == 0 {
		return SignInResponse{}, ErrSignInInvalidPassword
	}

	user, err := h.UserRepository.FindUserByUsername(ctx, request.Username)
	if err != nil {
		return SignInResponse{}, fmt.Errorf("failed to find user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		err = fmt.Errorf("failed to compare passwords: %w", err)
		return SignInResponse{}, errors.Join(err, ErrSignInInvalidPassword)
	} else if err != nil {
		return SignInResponse{}, fmt.Errorf("failed to compare passwords: %w", err)
	}

	accessTokenClaims := jwt.MapClaims{
		"iss": user.Username,
		"sub": "auth",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(h.AccessTokenExpiresIn).Unix(),
		"uid": user.ID,
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).
		SignedString([]byte(h.AccessTokenSecret))
	if err != nil {
		return SignInResponse{}, fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshTokenClaims := jwt.MapClaims{
		"iss": user.Username,
		"sub": "auth",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(h.RefreshTokenExpiresIn).Unix(),
		"uid": user.ID,
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).
		SignedString([]byte(h.RefreshTokenSecret))
	if err != nil {
		return SignInResponse{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	err = h.RefreshTokenRepository.SaveRefreshToken(ctx, user.ID, domainrefresh.Token(refreshToken))
	if err != nil {
		return SignInResponse{}, fmt.Errorf("failed to save refresh token: %w", err)
	}
	return SignInResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
