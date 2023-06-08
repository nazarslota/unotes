package note

import (
	"context"
	"fmt"

	domain "github.com/nazarslota/unotes/note/internal/domain/note"
)

type DeleteNoteRequest struct {
	ID string
}

type DeleteNoteResponse struct {
}

type DeleteNoteRequestHandler interface {
	Handle(ctx context.Context, request DeleteNoteRequest) (DeleteNoteResponse, error)
}

type deleteNoteRequestHandler struct {
	NoteDeleter NoteDeleter
}

var ErrDeleteNoteNotFound = func() error { return domain.ErrNoteNotFound }()

func NewDeleteNoteRequestHandler(noteDeleter NoteDeleter) DeleteNoteRequestHandler {
	return &deleteNoteRequestHandler{NoteDeleter: noteDeleter}
}

func (h deleteNoteRequestHandler) Handle(ctx context.Context, request DeleteNoteRequest) (DeleteNoteResponse, error) {
	if err := h.NoteDeleter.DeleteOne(ctx, request.ID); err != nil {
		return DeleteNoteResponse{}, fmt.Errorf("failed to delete note: %w", err)
	}
	return DeleteNoteResponse{}, nil
}
