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

func TestGetNotesRequestHandler_Handle(t *testing.T) {
	t.Run("should get notes", func(t *testing.T) {
		noteRepository := new(mockNoteRepository)
		notes := []*domainnote.Note{
			{
				ID:      "note1",
				Title:   "Test Note 1",
				Content: "This is a test note 1",
				UserID:  "user1",
			},
			{
				ID:      "note2",
				Title:   "Test Note 2",
				Content: "This is a test note 2",
				UserID:  "user1",
			},
		}
		noteRepository.On("FindMany", mock.Anything, mock.Anything).Return(notes, nil)

		getNotesRequest := &GetNotesRequest{
			UserID: "user1",
		}

		getNotesRequestHandler := NewGetNotesRequestHandler(noteRepository)
		response, err := getNotesRequestHandler.Handle(context.Background(), getNotesRequest)

		require.NoError(t, err)
		assert.Len(t, response.Notes, 2)
		assert.Equal(t, "note1", response.Notes[0].ID)
		assert.Equal(t, "Test Note 1", response.Notes[0].Title)
		assert.Equal(t, "This is a test note 1", response.Notes[0].Content)
		assert.Equal(t, "note2", response.Notes[1].ID)
		assert.Equal(t, "Test Note 2", response.Notes[1].Title)
		assert.Equal(t, "This is a test note 2", response.Notes[1].Content)
		noteRepository.AssertExpectations(t)
	})

	t.Run("should return an error if getting notes failed", func(t *testing.T) {
		noteRepository := new(mockNoteRepository)
		noteRepository.On("FindMany", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to find notes"))

		getNotesRequest := &GetNotesRequest{
			UserID: "user1",
		}

		getNotesRequestHandler := NewGetNotesRequestHandler(noteRepository)
		_, err := getNotesRequestHandler.Handle(context.Background(), getNotesRequest)

		if assert.Error(t, err) {
			assert.EqualError(t, err, "failed to find notes: failed to find notes")
		}
		noteRepository.AssertExpectations(t)
	})
}
