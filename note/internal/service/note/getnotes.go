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
	Notes []struct {
		ID      string
		Title   string
		Content string
	}
}

type GetNotesRequestHandler interface {
	Handle(ctx context.Context, request *GetNotesRequest) (*GetNotesResponse, error)
}

type getNotesRequestHandler struct {
	NoteRepository domainnote.Repository
}

func NewGetNotesRequestHandler(noteRepository domainnote.Repository) GetNotesRequestHandler {
	return &getNotesRequestHandler{NoteRepository: noteRepository}
}

func (c getNotesRequestHandler) Handle(ctx context.Context, request *GetNotesRequest) (*GetNotesResponse, error) {
	notes, err := c.NoteRepository.FindMany(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find notes: %w", err)
	}

	response := &GetNotesResponse{
		Notes: make([]struct {
			ID      string
			Title   string
			Content string
		}, 0),
	}

	for _, note := range notes {
		response.Notes = append(response.Notes, struct {
			ID      string
			Title   string
			Content string
		}{ID: note.ID, Title: note.Title, Content: note.Content})
	}
	return response, nil
}
