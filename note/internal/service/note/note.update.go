package note

import (
	"context"
	"fmt"
	"time"

	domain "github.com/nazarslota/unotes/note/internal/domain/note"
)

type UpdateNoteRequest struct {
	ID                string
	NewTitle          string
	NewContent        string
	NewPriority       *string
	NewCompletionTime *time.Time
}

type UpdateNoteResponse struct {
}

type UpdateNoteRequestHandler interface {
	Handle(ctx context.Context, request UpdateNoteRequest) (UpdateNoteResponse, error)
}

type updateNoteRequestHandler struct {
	NoteUpdater NoteUpdater
}

var ErrUpdateNoteNotFound = func() error { return domain.ErrNoteNotFound }()

func NewUpdateNoteRequestHandler(noteUpdater NoteUpdater) UpdateNoteRequestHandler {
	return &updateNoteRequestHandler{NoteUpdater: noteUpdater}
}

func (h updateNoteRequestHandler) Handle(ctx context.Context, request UpdateNoteRequest) (UpdateNoteResponse, error) {
	note := domain.Note{
		ID:             request.ID,
		Title:          request.NewTitle,
		Content:        request.NewContent,
		Priority:       request.NewPriority,
		CompletionTime: request.NewCompletionTime,
	}

	if err := h.NoteUpdater.UpdateOne(ctx, note); err != nil {
		return UpdateNoteResponse{}, fmt.Errorf("failed to update note: %w", err)
	}
	return UpdateNoteResponse{}, nil
}
