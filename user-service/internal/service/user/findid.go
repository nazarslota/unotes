package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/udholdenhed/unotes/user-service/internal/domain/user"
	"github.com/udholdenhed/unotes/user-service/pkg/errors"
)

type FindUserByIDRequest struct {
	ID string
}

type FindUserByIDResult struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
}

type FindUserByIDRequestHandler interface {
	Handle(ctx context.Context, r FindUserByIDRequest) (FindUserByIDResult, error)
}

type findUserByIDRequestHandler struct {
	UserRepository user.Repository
}

func NewFindUserByIDRequestHandler(userRepository user.Repository) FindUserByIDRequestHandler {
	return &findUserByIDRequestHandler{UserRepository: userRepository}
}

func (h *findUserByIDRequestHandler) Handle(ctx context.Context, r FindUserByIDRequest) (FindUserByIDResult, error) {
	u, err := h.UserRepository.FindByID(ctx, r.ID)
	if err != nil {
		return FindUserByIDResult{}, errors.ErrInternalServerError.SetInternal(err)
	} else if u == nil {
		return FindUserByIDResult{}, errors.NewHTTPError(
			http.StatusNotFound,
			fmt.Sprintf("user with this username does not exist"),
		)
	}

	res := FindUserByIDResult{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	}
	return res, nil
}
