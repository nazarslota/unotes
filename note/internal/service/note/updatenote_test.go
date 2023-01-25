package note

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateNoteRequestHandler_Handle(t *testing.T) {
	t.Run("should update a note", func(t *testing.T) {
		noteRepository := new(mockNoteRepository)
		noteRepository.On("UpdateOne", mock.Anything, mock.Anything).Return(nil)

		updateNoteRequest := UpdateNoteRequest{
			ID:         "Test Note ID",
			NewTitle:   "Test Note Updated",
			NewContent: "This is an updated test note",
		}

		updateNoteRequestHandler := NewUpdateNoteRequestHandler(noteRepository)
		_, err := updateNoteRequestHandler.Handle(context.Background(), updateNoteRequest)
		assert.NoError(t, err)

		noteRepository.AssertExpectations(t)
	})

	t.Run("should return an error if updating a note failed", func(t *testing.T) {
		noteRepository := new(mockNoteRepository)
		noteRepository.On("UpdateOne", mock.Anything, mock.Anything).Return(fmt.Errorf("failed to update note"))

		updateNoteRequest := UpdateNoteRequest{
			ID:         "Test Note ID",
			NewTitle:   "Test Note Updated",
			NewContent: "This is an updated test note",
		}

		updateNoteRequestHandler := NewUpdateNoteRequestHandler(noteRepository)
		_, err := updateNoteRequestHandler.Handle(context.Background(), updateNoteRequest)

		if assert.Error(t, err) {
			assert.EqualError(t, err, "failed to update note: failed to update note")
		}
		noteRepository.AssertExpectations(t)
	})
	t.Run("should return an error if updating a note not found", func(t *testing.T) {
		noteRepository := new(mockNoteRepository)
		noteRepository.On("UpdateOne", mock.Anything, mock.Anything).Return(ErrUpdateNoteNotFound)

		updateNoteRequest := UpdateNoteRequest{
			ID:         "Test Note ID",
			NewTitle:   "Test Note Updated",
			NewContent: "This is an updated test note",
		}

		updateNoteRequestHandler := NewUpdateNoteRequestHandler(noteRepository)
		_, err := updateNoteRequestHandler.Handle(context.Background(), updateNoteRequest)

		if assert.Error(t, err) {
			assert.EqualError(t, err, "failed to update note: not found")
		}
		noteRepository.AssertExpectations(t)
	})
}
