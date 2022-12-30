package user

import (
	"context"
	"errors"
)

type Repository interface {
	SaveOne(ctx context.Context, user *User) error
	FindOne(ctx context.Context, username string) (*User, error)
}

var ErrUserNotFound = errors.New("user not found")
