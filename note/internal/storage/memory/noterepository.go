package memory

import (
	"context"
	"fmt"
	"sync"

	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
)

type NoteRepository struct {
	notes *sync.Map
}

var _ domainnote.Repository = (*NoteRepository)(nil)

func NewNoteRepository() *NoteRepository {
	return &NoteRepository{notes: new(sync.Map)}
}

func (r NoteRepository) SaveOne(ctx context.Context, note domainnote.Note) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("saving note failed: %w", ctx.Err())
	default:
	}

	if _, loaded := r.notes.LoadOrStore(note.ID, note); loaded {
		return fmt.Errorf("saving note failed: %w", domainnote.ErrNoteAlreadyExist)
	}
	return nil
}

func (r NoteRepository) FindOne(ctx context.Context, noteID string) (domainnote.Note, error) {
	select {
	case <-ctx.Done():
		return domainnote.Note{}, fmt.Errorf("finding note failed: %w", ctx.Err())
	default:
	}

	value, ok := r.notes.Load(noteID)
	if !ok {
		return domainnote.Note{}, fmt.Errorf("finding note failed: %w", domainnote.ErrNoteNotFound)
	}
	return value.(domainnote.Note), nil
}

func (r NoteRepository) FindMany(ctx context.Context, userID string) ([]domainnote.Note, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("finding notes failed: %w", ctx.Err())
	default:
	}

	notes := make([]domainnote.Note, 0)
	r.notes.Range(func(_, value any) bool {
		if note := value.(domainnote.Note); note.UserID == userID {
			notes = append(notes, note)
		}
		return true
	})

	if len(notes) == 0 {
		return nil, fmt.Errorf("finding notes failed: %w", domainnote.ErrNoteNotFound)
	}
	return notes, nil
}

func (r NoteRepository) UpdateOne(ctx context.Context, note domainnote.Note) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("updating note failed: %w", ctx.Err())
	default:
	}

	if _, loaded := r.notes.LoadAndDelete(note.ID); !loaded {
		return fmt.Errorf("updatind note failed: %w", domainnote.ErrNoteNotFound)
	}

	r.notes.Store(note.ID, note)
	return nil
}

func (r NoteRepository) DeleteOne(ctx context.Context, noteID string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("deleting note failed: %w", ctx.Err())
	default:
	}

	if _, loaded := r.notes.LoadAndDelete(noteID); !loaded {
		return fmt.Errorf("deleting note failed: %w", domainnote.ErrNoteNotFound)
	}
	return nil
}

func (r NoteRepository) FindManyAsync(ctx context.Context, userID string) (<-chan domainnote.Note, <-chan error) {
	nts := make(chan domainnote.Note)
	errs := make(chan error)

	go func() {
		defer close(nts)
		defer close(errs)

		count := 0
		r.notes.Range(func(_, value any) bool {
			select {
			case <-ctx.Done():
				errs <- fmt.Errorf("finding notes failed: %w", ctx.Err())
				return false
			default:
			}

			if note := value.(domainnote.Note); note.UserID == userID {
				nts <- note
				count++
			}
			return true
		})

		if count == 0 {
			errs <- fmt.Errorf("finding notes failed: %w", domainnote.ErrNoteNotFound)
		}
	}()
	return nts, errs
}
