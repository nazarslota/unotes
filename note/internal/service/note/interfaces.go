package note

import (
	"context"

	domain "github.com/nazarslota/unotes/note/internal/domain/note"
)

type NoteSaver interface {
	SaveOne(ctx context.Context, note domain.Note) error
}

type NoteFinder interface {
	FindOne(ctx context.Context, noteID string) (domain.Note, error)
	FindMany(ctx context.Context, userID string) ([]domain.Note, error)
	FindManyAsync(ctx context.Context, userID string) (<-chan domain.Note, <-chan error)
}

type NoteUpdater interface {
	UpdateOne(ctx context.Context, note domain.Note) error
}

type NoteDeleter interface {
	DeleteOne(ctx context.Context, noteID string) error
}
