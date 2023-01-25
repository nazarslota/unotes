package note

import (
	"context"
	"fmt"

	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
)

type GetNotesRequest struct {
	UserID string
}

type GetNotesResponse struct {
	Notes []domainnote.Note
}

type GetNotesRequestHandler interface {
	Handle(ctx context.Context, request GetNotesRequest) (GetNotesResponse, error)
}

type getNotesRequestHandler struct {
	NoteRepository domainnote.Repository
}

var ErrGetNotesNotFound = func() error { return domainnote.ErrNoteNotFound }()

func NewGetNotesRequestHandler(noteRepository domainnote.Repository) GetNotesRequestHandler {
	return &getNotesRequestHandler{NoteRepository: noteRepository}
}

func (h getNotesRequestHandler) Handle(ctx context.Context, request GetNotesRequest) (GetNotesResponse, error) {
	notes, err := h.NoteRepository.FindMany(ctx, request.UserID)
	if err != nil {
		return GetNotesResponse{}, fmt.Errorf("failed to find notes: %w", err)
	}
	return GetNotesResponse{Notes: notes}, nil
}
