package note

import (
	"context"
	"errors"
	"fmt"

	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
)

type DeleteNoteRequest struct {
	ID string
}

type DeleteNoteResponse struct {
}

type DeleteNoteRequestHandler interface {
	Handle(ctx context.Context, request *DeleteNoteRequest) (*DeleteNoteResponse, error)
}

type deleteNoteRequestHandler struct {
	NoteRepository domainnote.Repository
}

var ErrDeleteNoteNotFound = func() error { return domainnote.ErrNotFound }()

func NewDeleteNoteRequestHandler(noteRepository domainnote.Repository) DeleteNoteRequestHandler {
	return &deleteNoteRequestHandler{NoteRepository: noteRepository}
}

func (c deleteNoteRequestHandler) Handle(ctx context.Context, request *DeleteNoteRequest) (*DeleteNoteResponse, error) {
	if err := c.NoteRepository.DeleteOne(ctx, request.ID); errors.Is(err, domainnote.ErrNotFound) {
		return nil, fmt.Errorf("failed to delete note: %w", ErrDeleteNoteNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("failed to delete note: %w", err)
	}
	return &DeleteNoteResponse{}, nil
}
