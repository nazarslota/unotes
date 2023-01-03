package oauth2

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nazarslota/unotes/auth/internal/domain/refreshtoken"
	"github.com/nazarslota/unotes/auth/internal/domain/user"
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
	Handle(ctx context.Context, request *SignInRequest) (*SignInResponse, error)
}

type signInRequestHandler struct {
	AccessTokenSecret      string
	RefreshTokenSecret     string
	AccessTokenExpiresIn   time.Duration
	RefreshTokenExpiresIn  time.Duration
	UserRepository         user.Repository
	RefreshTokenRepository refreshtoken.Repository
}

var (
	ErrSignInUserNotFound    = errors.New("user not found")
	ErrSignInInvalidPassword = errors.New("invalid user password")
)

func NewSignInRequestHandler(
	accessTokenSecret, refreshTokenSecret string,
	accessTokenExpiresIn, refreshTokenExpiresIn time.Duration,
	userRepository user.Repository, refreshTokenRepository refreshtoken.Repository,
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

func (h *signInRequestHandler) Handle(ctx context.Context, request *SignInRequest) (*SignInResponse, error) {
	u, err := h.UserRepository.FindOne(ctx, request.Username)
	if errors.Is(err, user.ErrUserNotFound) {
		return nil, fmt.Errorf("the user is not signed up: %w", ErrSignInUserNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("failed to verify the user sign up: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(request.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, fmt.Errorf("invalid user password: %w", ErrSignInInvalidPassword)
		}
		return nil, fmt.Errorf("failed to compare passwords: %w", err)
	}

	response := new(SignInResponse)
	response.AccessToken, err = newHS256(h.AccessTokenSecret, h.AccessTokenExpiresIn, u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create an access token: %w", err)
	}

	response.RefreshToken, err = newHS256(h.RefreshTokenSecret, h.RefreshTokenExpiresIn, u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create a refresh token: %w", err)
	}

	err = h.RefreshTokenRepository.SaveOne(ctx, u.ID, &refreshtoken.Token{Token: response.RefreshToken})
	if err != nil {
		return nil, fmt.Errorf("failed to save the refresh token: %w", err)
	}
	return response, nil
}
