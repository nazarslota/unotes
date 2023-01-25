package note

import (
	"context"
	"errors"
	"fmt"

	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
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
	Handle(ctx context.Context, request *GetNoteRequest) (*GetNoteResponse, error)
}

type getNoteRequestHandler struct {
	NoteRepository domainnote.Repository
}

var ErrGetNoteNotFound = func() error { return domainnote.ErrNoteNotFound }()

func NewGetNoteRequestHandler(noteRepository domainnote.Repository) GetNoteRequestHandler {
	return &getNoteRequestHandler{NoteRepository: noteRepository}
}

func (h getNoteRequestHandler) Handle(ctx context.Context, request *GetNoteRequest) (*GetNoteResponse, error) {
	note, err := h.NoteRepository.FindOne(ctx, request.ID)
	if errors.Is(err, domainnote.ErrNoteNotFound) {
		return nil, fmt.Errorf("failed to find note: %w", ErrGetNoteNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("failed to find note: %w", err)
	}

	response := &GetNoteResponse{
		Title:   note.Title,
		Content: note.Content,
		UserID:  note.UserID,
	}
	return response, nil
}
