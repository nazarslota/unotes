package oauth2

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/udholdenhed/unotes/auth/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Username string
	Password string
}

type SignUpResponse struct {
}

type SignUpRequestHandler interface {
	Handler(ctx context.Context, request *SignUpRequest) (*SignUpResponse, error)
}

type signUpRequestHandler struct {
	UserRepository user.Repository
}

var (
	ErrSignUpUserAlreadyExist = errors.New("user already exist")
)

func NewSignUpRequestHandler(userRepository user.Repository) SignUpRequestHandler {
	return &signUpRequestHandler{UserRepository: userRepository}
}

func (h *signUpRequestHandler) Handler(ctx context.Context, request *SignUpRequest) (*SignUpResponse, error) {
	signup, err := h.isSignUp(ctx, request.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to verify the user sign up: %w", err)
	} else if signup {
		return nil, fmt.Errorf("the user is already signed up: %w", ErrSignUpUserAlreadyExist)
	}

	err = h.signUp(ctx, request.Username, request.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to sign up the user: %w", err)
	}
	return &SignUpResponse{}, nil
}

func (h *signUpRequestHandler) isSignUp(ctx context.Context, username string) (bool, error) {
	if _, err := h.UserRepository.FindOne(ctx, username); err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to fetch the user from the repository: %w", err)
	}
	return true, nil
}

func (h *signUpRequestHandler) signUp(ctx context.Context, username string, password string) error {
	u := &user.User{
		ID:       uuid.New().String(),
		Username: username,
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %w", err)
	}
	u.PasswordHash = string(hash)

	if err := h.UserRepository.SaveOne(ctx, u); err != nil {
		return fmt.Errorf("failed to save the user: %w", err)
	}
	return nil
}
