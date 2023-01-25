package note

import (
	"context"
	"fmt"

	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
)

type CreateNoteRequest struct {
	ID      string
	Title   string
	Content string
	UserID  string
}

type CreateNoteResponse struct {
}

type CreateNoteRequestHandler interface {
	Handle(ctx context.Context, request CreateNoteRequest) (CreateNoteResponse, error)
}

type createNoteRequestHandler struct {
	NoteRepository domainnote.Repository
}

var ErrCreateNoteAlreadyExist = func() error { return domainnote.ErrNoteAlreadyExist }()

func NewCreateNoteRequestHandler(noteRepository domainnote.Repository) CreateNoteRequestHandler {
	return &createNoteRequestHandler{NoteRepository: noteRepository}
}

func (h createNoteRequestHandler) Handle(ctx context.Context, request CreateNoteRequest) (CreateNoteResponse, error) {
	note := domainnote.Note{
		ID:      request.ID,
		Title:   request.Title,
		Content: request.Content,
		UserID:  request.UserID,
	}

	if err := h.NoteRepository.SaveOne(ctx, note); err != nil {
		return CreateNoteResponse{}, fmt.Errorf("failed to create note: %w", err)
	}
	return CreateNoteResponse{}, nil
}
