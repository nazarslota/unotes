package note

import (
	"context"
	"fmt"

	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
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
	NoteRepository domainnote.Repository
}

var ErrDeleteNoteNotFound = func() error { return domainnote.ErrNoteNotFound }()

func NewDeleteNoteRequestHandler(noteRepository domainnote.Repository) DeleteNoteRequestHandler {
	return &deleteNoteRequestHandler{NoteRepository: noteRepository}
}

func (h deleteNoteRequestHandler) Handle(ctx context.Context, request DeleteNoteRequest) (DeleteNoteResponse, error) {
	if err := h.NoteRepository.DeleteOne(ctx, request.ID); err != nil {
		return DeleteNoteResponse{}, fmt.Errorf("failed to delete note: %w", err)
	}
	return DeleteNoteResponse{}, nil
}
