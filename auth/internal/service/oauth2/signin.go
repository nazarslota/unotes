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

// SignInRequest represents a sign in request.
type SignInRequest struct {
	Username string
	Password string
}

// SignInResponse represents a sign in response.
type SignInResponse struct {
	AccessToken  string
	RefreshToken string
}

// SignInRequestHandler is an interface that defines a sign in request handler.
type SignInRequestHandler interface {
	Handle(ctx context.Context, request *SignInRequest) (*SignInResponse, error)
}

// signInRequestHandler is a sign in request handler that verifies the sign in request, generates
// access and refresh tokens, and saves the refresh token.
type signInRequestHandler struct {
	AccessTokenSecret      string
	RefreshTokenSecret     string
	AccessTokenExpiresIn   time.Duration
	RefreshTokenExpiresIn  time.Duration
	UserRepository         user.Repository
	RefreshTokenRepository refreshtoken.Repository
}

// ErrSignInUserNotFound is returned when the user is not signed up.
var ErrSignInUserNotFound = errors.New("user not found")

// ErrSignInInvalidPassword is returned when the password is invalid.
var ErrSignInInvalidPassword = errors.New("invalid user password")

// NewSignInRequestHandler creates a new sign in request handler.
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

// Handle handles a sign in request and returns a response.
//
// It can return the following errors:
//   - ErrSignInUserNotFound: if the user is not signed up
//   - ErrSignInInvalidPassword: if the password is invalid
//   - other errors: if an error occurred while comparing passwords or creating the access or refresh tokens
func (h *signInRequestHandler) Handle(ctx context.Context, request *SignInRequest) (*SignInResponse, error) {
	// Check if the user is signed up.
	u, err := h.UserRepository.FindOne(ctx, request.Username)
	if errors.Is(err, user.ErrUserNotFound) {
		return nil, fmt.Errorf("the user is not signed up: %w", ErrSignInUserNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("failed to verify the user sign up: %w", err)
	}

	// Check if the password is correct.
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

	// Save the refresh token.
	if err := h.RefreshTokenRepository.SaveOne(ctx, u.ID, &refreshtoken.Token{Token: response.RefreshToken}); err != nil {
		return nil, fmt.Errorf("failed to save the refresh token: %w", err)
	}

	return response, nil
}
