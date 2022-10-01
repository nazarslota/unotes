package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/udholdenhed/unotes/user-service/internal/domain/user"
	"github.com/udholdenhed/unotes/user-service/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Username string
	Email    string
	Password string
}

type CreateUserRequestHandler interface {
	Handle(ctx context.Context, r CreateUserRequest) error
}

type createUserRequestHandler struct {
	UserRepository user.Repository
}

func NewCreateUserRequestHandler(userRepository user.Repository) CreateUserRequestHandler {
	return &createUserRequestHandler{UserRepository: userRepository}
}

func (h *createUserRequestHandler) Handle(ctx context.Context, r CreateUserRequest) error {
	u, err := h.UserRepository.FindByUsername(ctx, r.Username)
	if err != nil {
		return errors.ErrInternalServerError.SetInternal(err)
	} else if u != nil {
		return errors.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("user with this username already exists"),
		)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.ErrInternalServerError.SetInternal(err)
	}

	u = &user.User{
		ID:           uuid.New().String(),
		Username:     r.Username,
		Email:        r.Email,
		PasswordHash: string(passwordHash),
	}

	if err := h.UserRepository.Create(ctx, *u); err != nil {
		return errors.ErrInternalServerError.SetInternal(err)
	}

	return nil
}
