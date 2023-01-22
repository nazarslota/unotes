package note

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteNoteRequestHandler_Handle(t *testing.T) {
	t.Run("should delete a note", func(t *testing.T) {
		noteRepository := new(mockNoteRepository)
		noteRepository.On("DeleteOne", mock.Anything, mock.Anything).Return(nil)

		deleteNoteRequest := &DeleteNoteRequest{
			ID: "note1",
		}

		deleteNoteRequestHandler := NewDeleteNoteRequestHandler(noteRepository)
		response, err := deleteNoteRequestHandler.Handle(context.Background(), deleteNoteRequest)

		require.NoError(t, err)
		assert.NotNil(t, response)
		noteRepository.AssertExpectations(t)
	})

	t.Run("should return an error if deleting a note failed", func(t *testing.T) {
		noteRepository := new(mockNoteRepository)
		noteRepository.On("DeleteOne", mock.Anything, mock.Anything).Return(fmt.Errorf("failed to delete note"))

		deleteNoteRequest := &DeleteNoteRequest{
			ID: "note1",
		}

		deleteNoteRequestHandler := NewDeleteNoteRequestHandler(noteRepository)
		_, err := deleteNoteRequestHandler.Handle(context.Background(), deleteNoteRequest)

		if assert.Error(t, err) {
			assert.EqualError(t, err, "failed to delete note: failed to delete note")
		}
		noteRepository.AssertExpectations(t)
	})
}
