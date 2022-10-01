package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/udholdenhed/unotes/user-service/internal/domain/user"
	"github.com/udholdenhed/unotes/user-service/pkg/errors"
)

type FindUserByUsernameRequest struct {
	Username string
}

type FindUserByUsernameResult struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
}

type FindUserByUsernameRequestHandler interface {
	Handle(ctx context.Context, r FindUserByUsernameRequest) (FindUserByUsernameResult, error)
}

type findUserByUsernameRequestHandler struct {
	UserRepository user.Repository
}

func NewFindByUsernameRequestHandler(userRepository user.Repository) FindUserByUsernameRequestHandler {
	return &findUserByUsernameRequestHandler{UserRepository: userRepository}
}

func (h *findUserByUsernameRequestHandler) Handle(ctx context.Context, r FindUserByUsernameRequest) (FindUserByUsernameResult, error) {
	u, err := h.UserRepository.FindByUsername(ctx, r.Username)
	if err != nil {
		return FindUserByUsernameResult{}, errors.ErrInternalServerError.SetInternal(err)
	} else if u == nil {
		return FindUserByUsernameResult{}, errors.NewHTTPError(
			http.StatusNotFound,
			fmt.Sprintf("user with this username does not exist"),
		)
	}

	res := FindUserByUsernameResult{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	}
	return res, nil
}
