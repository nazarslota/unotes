package oauth2

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/udholdenhed/unotes/auth/internal/domain/refreshtoken"
	"github.com/udholdenhed/unotes/auth/internal/domain/user"
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
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context is done: %w", ctx.Err()) // failed to handle sign in request,
	default:
	}

	u, err := h.UserRepository.FindOne(ctx, request.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to find the user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(request.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, fmt.Errorf("invalid password: %w", ErrInvalidPassword)
		}
		return nil, fmt.Errorf("failed to compare passwords: %w", err)
	}

	response := new(SignInResponse)
	access, err := tokens{}.NewHS256(h.AccessTokenSecret, h.AccessTokenExpiresIn, u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create an access token: %w", err)
	}
	response.AccessToken = access

	refresh, err := tokens{}.NewHS256(h.RefreshTokenSecret, h.RefreshTokenExpiresIn, u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create a refresh token: %w", err)
	}
	response.RefreshToken = refresh

	err = h.RefreshTokenRepository.SaveOne(ctx, u.ID, &refreshtoken.Token{Token: refresh})
	if err != nil {
		return nil, fmt.Errorf("failed to save the refresh token: %w", err)
	}
	return response, nil
}
