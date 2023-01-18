package service

import (
	"testing"

	servicenote "github.com/nazarslota/unotes/note/internal/service/note"
	"github.com/stretchr/testify/assert"
)

func TestNewNoteService(t *testing.T) {
	service := NewNoteService(NoteServiceOptions{NoteRepository: nil})

	assert.IsType(t, servicenote.NewCreateNoteRequestHandler(nil), service.CreateNoteRequestHandler)
	assert.IsType(t, servicenote.NewDeleteNoteRequestHandler(nil), service.DeleteNoteRequestHandler)
	assert.IsType(t, servicenote.NewGetNoteRequestHandler(nil), service.GetNoteRequestHandler)
	assert.IsType(t, servicenote.NewGetNotesRequestHandler(nil), service.GetNotesRequestHandler)
}
