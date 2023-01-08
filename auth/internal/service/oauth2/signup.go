package oauth2

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nazarslota/unotes/auth/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

// SignUpRequest represents a request to sign up a new user.
type SignUpRequest struct {
	// Username is the user's desired username.
	Username string
	// Password is the user's desired password.
	Password string
}

// SignUpResponse represents a response to a sign up request.
type SignUpResponse struct {
}

// SignUpRequestHandler is an interface for handling sign up requests.
type SignUpRequestHandler interface {
	// Handler handles a sign up request and returns a response.
	Handler(ctx context.Context, request *SignUpRequest) (*SignUpResponse, error)
}

// signUpRequestHandler is an implementation of the SignUpRequestHandler interface.
type signUpRequestHandler struct {
	// UserRepository is the repository used to store and retrieve users.
	UserRepository user.Repository
}

// ErrSignUpUserAlreadyExist is returned when a user with the same username already exists.
var ErrSignUpUserAlreadyExist = errors.New("user already exists")

// NewSignUpRequestHandler returns a new SignUpRequestHandler.
func NewSignUpRequestHandler(userRepository user.Repository) SignUpRequestHandler {
	return &signUpRequestHandler{UserRepository: userRepository}
}

// Handler handles a sign-up request and returns a response.
//
// It can return the following errors:
//   - ErrSignUpUserAlreadyExist: if a user with the same username already exists
//   - other errors: if an error occurred while generating the password hash or saving the user to the repository
func (h *signUpRequestHandler) Handler(ctx context.Context, request *SignUpRequest) (*SignUpResponse, error) {
	// Check if the user is already signed up.
	signedUp, err := h.isSignUp(ctx, request.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to verify the user sign up: %w", err)
	} else if signedUp {
		return nil, fmt.Errorf("the user is already signed up: %w", ErrSignUpUserAlreadyExist)
	}

	// Sign up the user.
	if err := h.signUp(ctx, request.Username, request.Password); err != nil {
		return nil, fmt.Errorf("failed to sign up the user: %w", err)
	}
	return &SignUpResponse{}, nil
}

// isSignUp checks if the user is already signed up.
func (h *signUpRequestHandler) isSignUp(ctx context.Context, username string) (bool, error) {
	// Check if the user exists in the repository.
	if _, err := h.UserRepository.FindOne(ctx, username); err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to fetch the user from the repository: %w", err)
	}
	return true, nil
}

// signUp creates and saves a new user to the repository.
func (h *signUpRequestHandler) signUp(ctx context.Context, username string, password string) error {
	// Create a new user.
	u := &user.User{
		ID:       uuid.New().String(),
		Username: username,
	}

	// Generate a password hash.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %w", err)
	}
	u.PasswordHash = string(hash)

	// Save the user to the repository.
	if err := h.UserRepository.SaveOne(ctx, u); err != nil {
		return fmt.Errorf("failed to save the user: %w", err)
	}
	return nil
}
