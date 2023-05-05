package note

import (
	"context"

	domain "github.com/nazarslota/unotes/note/internal/domain/note"
)

type GetNotesAsyncRequest struct {
	UserID string
}

type GetNotesAsyncResponse struct {
	Notes <-chan domain.Note
}

type GetNotesAsyncRequestHandler interface {
	Handle(ctx context.Context, request GetNotesAsyncRequest) (GetNotesAsyncResponse, <-chan error)
}

type getNotesAsyncRequestHandler struct {
	NoteFinder NoteFinder
}

var ErrGetNotesAsyncNotFound = func() error { return domain.ErrNoteNotFound }()

func NewGetNotesAsyncRequestHandler(noteFinder NoteFinder) GetNotesAsyncRequestHandler {
	return &getNotesAsyncRequestHandler{NoteFinder: noteFinder}
}

func (h getNotesAsyncRequestHandler) Handle(ctx context.Context, request GetNotesAsyncRequest) (GetNotesAsyncResponse, <-chan error) {
	notes, errs := h.NoteFinder.FindManyAsync(ctx, request.UserID)
	return GetNotesAsyncResponse{Notes: notes}, errs
}
