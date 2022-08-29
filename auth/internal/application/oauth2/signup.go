package oauth2

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/udholdenhed/unotes/auth/internal/domain/user"
	"github.com/udholdenhed/unotes/auth/pkg/errors"
	"github.com/udholdenhed/unotes/auth/pkg/utils"
	"github.com/udholdenhed/unotes/auth/pkg/validator/password"
)

type SignUpRequest struct {
	Username string
	Password string
}

type SignUpRequestHandler interface {
	Handler(r SignUpRequest) error
}

type signUpRequestHandler struct {
	UserRepository user.Repository
}

func NewSignUpRequestHandler(r user.Repository) SignUpRequestHandler {
	return &signUpRequestHandler{UserRepository: r}
}

func (h *signUpRequestHandler) Handler(r SignUpRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u, err := h.UserRepository.FindOne(ctx, r.Username)
	if err != nil {
		return err
	} else if u != nil {
		return errors.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("user with this username already exists"),
		)
	}

	vf, _ := password.NewDefaultValidationFuncs(8, 64)
	if err := password.NewValidator(vf.ValidatePasswordLength).Validate(r.Password); err != nil {
		return errors.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("incorrect password length"),
		).SetInternal(err)
	}

	u = &user.User{
		ID:           utils.NewUUID(),
		Username:     r.Username,
		PasswordHash: utils.HashPassword(r.Password),
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.UserRepository.Create(ctx, *u); err != nil {
		return errors.ErrInternalServerError.SetInternal(err)
	}
	return nil
}
