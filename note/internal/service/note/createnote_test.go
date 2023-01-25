package note

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateNoteRequestHandler_Handle(t *testing.T) {
	t.Run("should save a note", func(t *testing.T) {
		noteRepository := new(mockNoteRepository)
		noteRepository.On("SaveOne", mock.Anything, mock.Anything).Return(nil)

		createNoteRequest := CreateNoteRequest{
			ID:      "Test Note ID",
			Title:   "Test Note",
			Content: "This is a test note",
			UserID:  "user1",
		}

		createNoteRequestHandler := NewCreateNoteRequestHandler(noteRepository)
		_, err := createNoteRequestHandler.Handle(context.Background(), createNoteRequest)
		assert.NoError(t, err)

		noteRepository.AssertExpectations(t)
	})

	t.Run("should return an error if saving a note failed", func(t *testing.T) {
		noteRepository := new(mockNoteRepository)
		noteRepository.On("SaveOne", mock.Anything, mock.Anything).Return(fmt.Errorf("failed to save note"))

		createNoteRequest := CreateNoteRequest{
			ID:      "Test Note ID",
			Title:   "Test Note",
			Content: "This is a test note",
			UserID:  "user1",
		}

		createNoteRequestHandler := NewCreateNoteRequestHandler(noteRepository)
		_, err := createNoteRequestHandler.Handle(context.Background(), createNoteRequest)

		if assert.Error(t, err) {
			assert.EqualError(t, err, "failed to create note: failed to save note")
		}
		noteRepository.AssertExpectations(t)
	})
}
