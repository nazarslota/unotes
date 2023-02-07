package note

import (
	"context"
	"fmt"

	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
)

type GetNoteRequest struct {
	ID string
}

type GetNoteResponse struct {
	ID      string
	Title   string
	Content string
	UserID  string
}

type GetNoteRequestHandler interface {
	Handle(ctx context.Context, request GetNoteRequest) (GetNoteResponse, error)
}

type getNoteRequestHandler struct {
	NoteRepository domainnote.Repository
}

var ErrGetNoteNotFound = func() error { return domainnote.ErrNoteNotFound }()

func NewGetNoteRequestHandler(noteRepository domainnote.Repository) GetNoteRequestHandler {
	return &getNoteRequestHandler{NoteRepository: noteRepository}
}

func (h getNoteRequestHandler) Handle(ctx context.Context, request GetNoteRequest) (GetNoteResponse, error) {
	note, err := h.NoteRepository.FindOne(ctx, request.ID)
	if err != nil {
		return GetNoteResponse{}, fmt.Errorf("failed to find note: %w", err)
	}
	return GetNoteResponse{ID: note.ID, Title: note.Title, Content: note.Content, UserID: note.UserID}, nil
}
