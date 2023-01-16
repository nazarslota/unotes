package note

import (
	"context"
	"errors"
)

type Repository interface {
	SaveOne(ctx context.Context, note *Note) error
	FindOne(ctx context.Context, noteID string) (*Note, error)
	FindMany(ctx context.Context, userID string) ([]*Note, error)
	DeleteOne(ctx context.Context, noteID string) error
}

var (
	ErrAlreadyExist = errors.New("already exist")
	ErrNotFound     = errors.New("not found")
)
