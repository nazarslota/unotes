package service

import (
	"testing"

	servicenote "github.com/nazarslota/unotes/note/internal/service/note"
	"github.com/stretchr/testify/assert"
)

func TestNewServices(t *testing.T) {
	noteServiceOptions := NoteServiceOptions{NoteRepository: nil}
	services := NewServices(noteServiceOptions)

	assert.IsType(t, NoteService{}, services.NoteService)
	assert.IsType(t, servicenote.NewCreateNoteRequestHandler(nil), services.NoteService.CreateNoteRequestHandler)
	assert.IsType(t, servicenote.NewGetNoteRequestHandler(nil), services.NoteService.GetNoteRequestHandler)
	assert.IsType(t, servicenote.NewGetNotesRequestHandler(nil), services.NoteService.GetNotesRequestHandler)
	assert.IsType(t, servicenote.NewUpdateNoteRequestHandler(nil), services.NoteService.UpdateNoteRequestHandler)
	assert.IsType(t, servicenote.NewDeleteNoteRequestHandler(nil), services.NoteService.DeleteNoteRequestHandler)
}
