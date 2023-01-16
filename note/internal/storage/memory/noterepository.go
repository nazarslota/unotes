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

func (r NoteRepository) SaveOne(ctx context.Context, note *domainnote.Note) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("saving note failed: %w", ctx.Err())
	default:
	}

	if _, loaded := r.notes.LoadOrStore(note.ID, note); loaded {
		return fmt.Errorf("saving note failed: %w", domainnote.ErrAlreadyExist)
	}
	return nil
}

func (r NoteRepository) FindOne(ctx context.Context, noteID string) (*domainnote.Note, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("finding note failed: %w", ctx.Err())
	default:
	}

	value, ok := r.notes.Load(noteID)
	if !ok {
		return nil, fmt.Errorf("finding note failed: %w", domainnote.ErrNotFound)
	}
	return value.(*domainnote.Note), nil
}

func (r NoteRepository) FindMany(ctx context.Context, userID string) ([]*domainnote.Note, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("finding notes failed: %w", ctx.Err())
	default:
	}

	notes := make([]*domainnote.Note, 0)
	f := func(_, value any) bool {
		note := value.(*domainnote.Note)
		if note.UserID == userID {
			notes = append(notes, note)
		}
		return true
	}
	r.notes.Range(f)

	if len(notes) == 0 {
		return nil, fmt.Errorf("finding notes failed: %w", domainnote.ErrNotFound)
	}
	return notes, nil
}

func (r NoteRepository) DeleteOne(ctx context.Context, noteID string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("deleting note failed: %w", ctx.Err())
	default:
	}

	if _, loaded := r.notes.LoadAndDelete(noteID); !loaded {
		return fmt.Errorf("deleting note failed: %w", domainnote.ErrNotFound)
	}
	return nil
}
