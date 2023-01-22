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

func NewGetNoteRequestHandler(noteRepository domainnote.Repository) GetNoteRequestHandler {
	return &getNoteRequestHandler{NoteRepository: noteRepository}
}

func (c getNoteRequestHandler) Handle(ctx context.Context, request *GetNoteRequest) (*GetNoteResponse, error) {
	note, err := c.NoteRepository.FindOne(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find notes: %w", err)
	}

	response := &GetNoteResponse{
		Title:   note.Title,
		Content: note.Content,
		UserID:  note.UserID,
	}
	return response, nil
}
