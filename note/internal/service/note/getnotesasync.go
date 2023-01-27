package note

import (
	"context"

	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
)

type GetNotesAsyncRequest struct {
	UserID string
}

type GetNotesAsyncResponse struct {
	Notes <-chan domainnote.Note
}

type GetNotesAsyncRequestHandler interface {
	Handle(ctx context.Context, request GetNotesAsyncRequest) (GetNotesAsyncResponse, <-chan error)
}

type getNotesAsyncRequestHandler struct {
	NoteRepository domainnote.Repository
}

var ErrGetNotesAsyncNotFound = func() error { return domainnote.ErrNoteNotFound }()

func NewGetNotesAsyncRequestHandler(noteRepository domainnote.Repository) GetNotesAsyncRequestHandler {
	return &getNotesAsyncRequestHandler{NoteRepository: noteRepository}
}

func (h getNotesAsyncRequestHandler) Handle(ctx context.Context, request GetNotesAsyncRequest) (GetNotesAsyncResponse, <-chan error) {
	notes, errs := h.NoteRepository.FindManyAsync(ctx, request.UserID)
	return GetNotesAsyncResponse{Notes: notes}, errs
}
