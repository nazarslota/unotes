package memory

import (
	"context"

	"github.com/udholdenhed/unotes/auth/internal/domain/user"
)

type UserRepository struct {
	Users map[string]user.User
}

var _ user.Repository = (*UserRepository)(nil)

func NewUserRepository() user.Repository {
	return &UserRepository{Users: make(map[string]user.User)}
}

func (ur *UserRepository) Create(_ context.Context, u user.User) error {
	ur.Users[u.ID] = u
	return nil
}

func (ur *UserRepository) FindOne(_ context.Context, username string) (*user.User, error) {
	for _, u := range ur.Users {
		if u.Username == username {
			return &u, nil
		}
	}
	return nil, nil
}
