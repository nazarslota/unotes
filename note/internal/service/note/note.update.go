package note

import (
	"context"
	"fmt"

	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
)

type UpdateNoteRequest struct {
	ID         string
	NewTitle   string
	NewContent string
}

type UpdateNoteResponse struct {
}

type UpdateNoteRequestHandler interface {
	Handle(ctx context.Context, request UpdateNoteRequest) (UpdateNoteResponse, error)
}

type updateNoteRequestHandler struct {
	NoteRepository domainnote.Repository
}

var ErrUpdateNoteNotFound = func() error { return domainnote.ErrNoteNotFound }()

func NewUpdateNoteRequestHandler(noteRepository domainnote.Repository) UpdateNoteRequestHandler {
	return &updateNoteRequestHandler{
		NoteRepository: noteRepository,
	}
}

func (h updateNoteRequestHandler) Handle(ctx context.Context, request UpdateNoteRequest) (UpdateNoteResponse, error) {
	note := domainnote.Note{
		ID:      request.ID,
		Title:   request.NewTitle,
		Content: request.NewContent,
	}

	if err := h.NoteRepository.UpdateOne(ctx, note); err != nil {
		return UpdateNoteResponse{}, fmt.Errorf("failed to update note: %w", err)
	}
	return UpdateNoteResponse{}, nil
}
