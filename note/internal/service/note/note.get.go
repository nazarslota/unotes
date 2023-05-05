package note

import (
	"context"
	"fmt"

	domain "github.com/nazarslota/unotes/note/internal/domain/note"
)

type GetNoteRequest struct {
	ID string
}

type GetNoteResponse struct {
	Title   string
	Content string
	UserID  string
}

type GetNoteRequestHandler interface {
	Handle(ctx context.Context, request GetNoteRequest) (GetNoteResponse, error)
}

type getNoteRequestHandler struct {
	NoteFinder NoteFinder
}

var ErrGetNoteNotFound = func() error { return domain.ErrNoteNotFound }()

func NewGetNoteRequestHandler(noteFinder NoteFinder) GetNoteRequestHandler {
	return &getNoteRequestHandler{NoteFinder: noteFinder}
}

func (h getNoteRequestHandler) Handle(ctx context.Context, request GetNoteRequest) (GetNoteResponse, error) {
	note, err := h.NoteFinder.FindOne(ctx, request.ID)
	if err != nil {
		return GetNoteResponse{}, fmt.Errorf("failed to find note: %w", err)
	}
	return GetNoteResponse{Title: note.Title, Content: note.Content, UserID: note.UserID}, nil
}
