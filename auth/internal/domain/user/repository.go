package user

import (
	"context"
	"errors"
)

// Repository is the interface for a user repository.
type Repository interface {
	// SaveOne saves a single user to the repository.
	SaveOne(ctx context.Context, user *User) error
	// FindOne retrieves a single user from the repository by username.
	FindOne(ctx context.Context, username string) (*User, error)
}

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = errors.New("user not found")
