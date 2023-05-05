package note

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domain "github.com/nazarslota/unotes/note/internal/domain/note"
)

type CreateNoteRequest struct {
	Title   string
	Content string
	UserID  string
}

type CreateNoteResponse struct {
	ID     string
	UserID string
}

type CreateNoteRequestHandler interface {
	Handle(ctx context.Context, request CreateNoteRequest) (CreateNoteResponse, error)
}

type createNoteRequestHandler struct {
	NoteSaver NoteSaver
}

var ErrCreateNoteAlreadyExist = func() error { return domain.ErrNoteAlreadyExist }()

func NewCreateNoteRequestHandler(noteSaver NoteSaver) CreateNoteRequestHandler {
	return &createNoteRequestHandler{NoteSaver: noteSaver}
}

func (h createNoteRequestHandler) Handle(ctx context.Context, request CreateNoteRequest) (CreateNoteResponse, error) {
	note := domain.Note{
		ID:      uuid.New().String(),
		Title:   request.Title,
		Content: request.Content,
		UserID:  request.UserID,
	}

	if err := h.NoteSaver.SaveOne(ctx, note); err != nil {
		return CreateNoteResponse{}, fmt.Errorf("failed to save note: %w", err)
	}
	return CreateNoteResponse{ID: note.ID, UserID: request.UserID}, nil
}
