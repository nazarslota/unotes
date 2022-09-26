package oauth2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/udholdenhed/unotes/auth-service/internal/domain/user"
	"github.com/udholdenhed/unotes/auth-service/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Username string
	Password string
}

type SignUpRequestHandler interface {
	Handler(ctx context.Context, r SignUpRequest) error
}

type signUpRequestHandler struct {
	UserRepository user.Repository
}

func NewSignUpRequestHandler(r user.Repository) SignUpRequestHandler {
	return &signUpRequestHandler{UserRepository: r}
}

func (h *signUpRequestHandler) Handler(ctx context.Context, r SignUpRequest) error {
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
		PasswordHash: string(passwordHash),
	}

	if err := h.UserRepository.Create(ctx, *u); err != nil {
		return errors.ErrInternalServerError.SetInternal(err)
	}

	return nil
}
