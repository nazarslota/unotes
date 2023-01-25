package note

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
)

type CreateNoteRequest struct {
	Title   string
	Content string
	UserID  string
}

type CreateNoteResponse struct {
	ID string
}

type CreateNoteRequestHandler interface {
	Handle(ctx context.Context, request *CreateNoteRequest) (*CreateNoteResponse, error)
}

type createNoteRequestHandler struct {
	NoteRepository domainnote.Repository
}

var ErrCreateNoteAlreadyExist = func() error { return domainnote.ErrNoteAlreadyExist }()

func NewCreateNoteRequestHandler(noteRepository domainnote.Repository) CreateNoteRequestHandler {
	return &createNoteRequestHandler{NoteRepository: noteRepository}
}

func (h createNoteRequestHandler) Handle(ctx context.Context, request *CreateNoteRequest) (*CreateNoteResponse, error) {
	note := domainnote.Note{
		ID:      uuid.New().String(),
		Title:   request.Title,
		Content: request.Content,
		UserID:  request.UserID,
	}

	if err := h.NoteRepository.SaveOne(ctx, note); errors.Is(err, domainnote.ErrNoteAlreadyExist) {
		return nil, fmt.Errorf("failed to create note: %w", ErrCreateNoteAlreadyExist)
	} else if err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}
	return &CreateNoteResponse{ID: note.ID}, nil
}
