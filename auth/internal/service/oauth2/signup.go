package oauth2

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	domainuser "github.com/nazarslota/unotes/auth/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Username string
	Password string
}

type SignUpResponse struct{}

type SignUpRequestHandler interface {
	Handler(ctx context.Context, request SignUpRequest) (SignUpResponse, error)
}

type signUpRequestHandler struct {
	UserRepository domainuser.Repository
}

var (
	ErrSignUpInvalidUsername   = errSignUpInvalidUsername()
	ErrSignUpInvalidPassword   = errSignUpInvalidPassword()
	ErrSignUpUserAlreadyExists = errSignUpUserAlreadyExists()
)

func errSignUpInvalidUsername() error   { return errors.New("invalid username") }
func errSignUpInvalidPassword() error   { return errors.New("invalid password") }
func errSignUpUserAlreadyExists() error { return domainuser.ErrUserAlreadyExists }

func NewSignUpRequestHandler(userRepository domainuser.Repository) SignUpRequestHandler {
	return &signUpRequestHandler{UserRepository: userRepository}
}

func (h signUpRequestHandler) Handler(ctx context.Context, request SignUpRequest) (SignUpResponse, error) {
	if len(request.Username) == 0 {
		return SignUpResponse{}, ErrSignUpInvalidUsername
	} else if len(request.Password) == 0 {
		return SignUpResponse{}, ErrSignUpInvalidPassword
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.MinCost)
	if err != nil {
		return SignUpResponse{}, fmt.Errorf("failed to generate password hash: %w", err)
	}

	if err := h.UserRepository.SaveUser(ctx, domainuser.User{
		ID:           uuid.New().String(),
		Username:     request.Username,
		PasswordHash: string(passwordHash),
	}); err != nil {
		return SignUpResponse{}, fmt.Errorf("failed to save user: %w", err)
	}
	return SignUpResponse{}, nil
}
