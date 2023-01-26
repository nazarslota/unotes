package note

import (
	"context"
	"errors"
)

type Repository interface {
	SaveOne(ctx context.Context, note Note) error
	FindOne(ctx context.Context, noteID string) (Note, error)
	FindMany(ctx context.Context, userID string) ([]Note, error)
	UpdateOne(ctx context.Context, note Note) error
	DeleteOne(ctx context.Context, noteID string) error
	FindManyAsync(ctx context.Context, userID string) (<-chan Note, <-chan error)
}

var (
	ErrNoteAlreadyExist = errors.New("already exist")
	ErrNoteNotFound     = errors.New("not found")
)
