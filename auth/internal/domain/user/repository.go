package user

import (
	"context"
	"errors"
)

type Repository interface {
	SaveUser(ctx context.Context, user User) error
	FindUserByUserID(ctx context.Context, userID string) (User, error)
	FindUserByUsername(ctx context.Context, username string) (User, error)
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
