package note

import (
	"context"
	"fmt"
	"testing"

	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetNoteRequestHandler_Handle(t *testing.T) {
	t.Run("should get a note", func(t *testing.T) {
		noteRepository := new(mockNoteRepository)
		note := domainnote.Note{
			ID:      "note1",
			Title:   "Test Note",
			Content: "This is a test note",
			UserID:  "user1",
		}
		noteRepository.On("FindOne", mock.Anything, mock.Anything).Return(note, nil)

		getNoteRequest := &GetNoteRequest{
			ID: "note1",
		}

		getNoteRequestHandler := NewGetNoteRequestHandler(noteRepository)
		response, err := getNoteRequestHandler.Handle(context.Background(), getNoteRequest)

		require.NoError(t, err)
		assert.Equal(t, "Test Note", response.Title)
		assert.Equal(t, "This is a test note", response.Content)
		assert.Equal(t, "user1", response.UserID)
		noteRepository.AssertExpectations(t)
	})

	t.Run("should return an error if getting a note failed", func(t *testing.T) {
		noteRepository := new(mockNoteRepository)
		noteRepository.On("FindOne", mock.Anything, mock.Anything).Return(domainnote.Note{}, fmt.Errorf("failed to find note"))

		getNoteRequest := &GetNoteRequest{
			ID: "note1",
		}

		getNoteRequestHandler := NewGetNoteRequestHandler(noteRepository)
		_, err := getNoteRequestHandler.Handle(context.Background(), getNoteRequest)

		if assert.Error(t, err) {
			assert.EqualError(t, err, "failed to find note: failed to find note")
		}
		noteRepository.AssertExpectations(t)
	})
}
