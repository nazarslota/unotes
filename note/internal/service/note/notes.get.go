package note

import (
	"context"
	"fmt"

	domain "github.com/nazarslota/unotes/note/internal/domain/note"
)

type GetNotesRequest struct {
	UserID string
}

type GetNotesResponse struct {
	Notes []domain.Note
}

type GetNotesRequestHandler interface {
	Handle(ctx context.Context, request GetNotesRequest) (GetNotesResponse, error)
}

type getNotesRequestHandler struct {
	NoteFinder NoteFinder
}

var ErrGetNotesNotFound = func() error { return domain.ErrNoteNotFound }()

func NewGetNotesRequestHandler(noteFinder NoteFinder) GetNotesRequestHandler {
	return &getNotesRequestHandler{NoteFinder: noteFinder}
}

func (h getNotesRequestHandler) Handle(ctx context.Context, request GetNotesRequest) (GetNotesResponse, error) {
	notes, err := h.NoteFinder.FindMany(ctx, request.UserID)
	if err != nil {
		return GetNotesResponse{}, fmt.Errorf("failed to find notes: %w", err)
	}
	return GetNotesResponse{Notes: notes}, nil
}
