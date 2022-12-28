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

func NewSignUpRequestHandler(userRepository user.Repository) SignUpRequestHandler {
	return &signUpRequestHandler{UserRepository: userRepository}
}

func (h *signUpRequestHandler) Handler(ctx context.Context, request *SignUpRequest) (*SignUpResponse, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context is done: %w", ctx.Err()) // failed to handle sign out request,
	default:
	}

	_, err := h.UserRepository.FindOne(ctx, request.Username)
	if !errors.Is(err, user.ErrUserNotFound) {
		return nil, fmt.Errorf("a user with this user name already exists: %w", ErrUserAlreadyExist)
	} else if err != nil && !errors.Is(err, user.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to find the user: %w", err)
	}

	u := &user.User{
		ID:       uuid.New().String(),
		Username: request.Username,
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password hash: %w", err)
	}
	u.PasswordHash = string(password)

	err = h.UserRepository.SaveOne(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("failed to save the user: %w", err)
	}
	return &SignUpResponse{}, nil
}
