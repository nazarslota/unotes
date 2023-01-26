package note

import (
	"context"

	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
	"github.com/stretchr/testify/mock"
)

type mockNoteRepository struct {
	mock.Mock
}

func (r *mockNoteRepository) SaveOne(ctx context.Context, note domainnote.Note) error {
	args := r.Called(ctx, note)
	return args.Error(0)
}

func (r *mockNoteRepository) FindOne(ctx context.Context, noteID string) (domainnote.Note, error) {
	args := r.Called(ctx, noteID)
	return args.Get(0).(domainnote.Note), args.Error(1)
}

func (r *mockNoteRepository) FindMany(ctx context.Context, userID string) ([]domainnote.Note, error) {
	args := r.Called(ctx, userID)
	notes, ok := args.Get(0).([]domainnote.Note)
	if ok {
		return notes, args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *mockNoteRepository) UpdateOne(ctx context.Context, note domainnote.Note) error {
	args := r.Called(ctx, note)
	return args.Error(0)
}

func (r *mockNoteRepository) DeleteOne(ctx context.Context, noteID string) error {
	args := r.Called(ctx, noteID)
	return args.Error(0)
}

func (r *mockNoteRepository) FindManyAsync(ctx context.Context, userID string) (<-chan domainnote.Note, <-chan error) {
	nts := make(chan domainnote.Note)
	errs := make(chan error)

	go func() {
		defer close(nts)
		defer close(errs)

		notes, err := r.FindMany(ctx, userID)
		if err != nil {
			errs <- err
			return
		}

		for _, note := range notes {
			nts <- note
		}
	}()
	return nts, errs
}
