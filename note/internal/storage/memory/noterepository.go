package memory

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

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

	value, loaded := r.notes.LoadAndDelete(note.ID)
	if !loaded {
		return fmt.Errorf("updatind note failed: %w", domainnote.ErrNoteNotFound)
	}

	note.UserID = value.(domainnote.Note).UserID
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
	notes, errs := make(chan domainnote.Note), make(chan error)
	select {
	case <-ctx.Done():
		go func() {
			errs <- fmt.Errorf("finding notes failed: %w", ctx.Err())
			close(errs)
		}()
		close(notes)
		return notes, errs
	default:
	}

	go func() {
		var found atomic.Bool
		var wgNotes, wgErrs sync.WaitGroup
		r.notes.Range(func(_, value any) bool {
			if note := value.(domainnote.Note); userID == note.UserID {
				wgNotes.Add(1)
				if !found.Load() {
					found.Store(true)
				}

				go func() {
					notes <- note
					wgNotes.Done()
				}()
			}
			return true
		})

		if !found.Load() {
			wgErrs.Add(1)
			go func() {
				errs <- fmt.Errorf("finding notes failed: %w", domainnote.ErrNoteNotFound)
				wgErrs.Done()
			}()
		}

		go func() {
			wgNotes.Wait()
			close(notes)
		}()
		go func() {
			wgErrs.Wait()
			close(errs)
		}()
	}()
	return notes, errs
}
