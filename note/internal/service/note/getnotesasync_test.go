package note

import (
	"context"
	"fmt"
	"testing"

	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetNotesAsyncRequestHandler_Handle(t *testing.T) {
	t.Run("should get notes", func(t *testing.T) {
		noteRepository := new(mockNoteRepository)
		notes := []domainnote.Note{
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

		getNotesAsyncRequestHandler := NewGetNotesAsyncRequestHandler(noteRepository)

		getNotesAsyncRequest := GetNotesAsyncRequest{UserID: "user1"}
		response, errs := getNotesAsyncRequestHandler.Handle(context.Background(), getNotesAsyncRequest)

		receivedNotes := make([]domainnote.Note, 0)
		for note := range response.Notes {
			receivedNotes = append(receivedNotes, note)
		}
		assert.Equal(t, notes, receivedNotes)

		receivedErrs := make([]error, 0)
		for err := range errs {
			receivedErrs = append(receivedErrs, err)
		}
		assert.Empty(t, receivedErrs)

		noteRepository.AssertExpectations(t)
	})

	t.Run("should return an error if getting notes failed", func(t *testing.T) {
		noteRepository := new(mockNoteRepository)
		noteRepository.On("FindMany", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to find notes"))

		getNotesAsyncRequestHandler := NewGetNotesAsyncRequestHandler(noteRepository)

		getNotesAsyncRequest := GetNotesAsyncRequest{UserID: "user1"}
		response, errs := getNotesAsyncRequestHandler.Handle(context.Background(), getNotesAsyncRequest)

		receivedNotes := make([]domainnote.Note, 0)
		for note := range response.Notes {
			receivedNotes = append(receivedNotes, note)
		}
		assert.Empty(t, receivedNotes)

		receivedErrs := make([]error, 0)
		for err := range errs {
			receivedErrs = append(receivedErrs, err)
		}

		if assert.Len(t, receivedErrs, 1) {
			assert.EqualError(t, receivedErrs[0], "failed to find notes")
		}

		noteRepository.AssertExpectations(t)
	})
}
