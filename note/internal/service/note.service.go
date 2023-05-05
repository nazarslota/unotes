package service

import (
	servicenote "github.com/nazarslota/unotes/note/internal/service/note"
)

type NoteService struct {
	CreateNoteRequestHandler    servicenote.CreateNoteRequestHandler
	GetNoteRequestHandler       servicenote.GetNoteRequestHandler
	GetNotesRequestHandler      servicenote.GetNotesRequestHandler
	UpdateNoteRequestHandler    servicenote.UpdateNoteRequestHandler
	DeleteNoteRequestHandler    servicenote.DeleteNoteRequestHandler
	GetNotesAsyncRequestHandler servicenote.GetNotesAsyncRequestHandler
}

type NoteServiceOptions struct {
	NoteSaver   servicenote.NoteSaver
	NoteFinder  servicenote.NoteFinder
	NoteUpdater servicenote.NoteUpdater
	NoteDeleter servicenote.NoteDeleter
}

func NewNoteService(options NoteServiceOptions) NoteService {
	return NoteService{
		CreateNoteRequestHandler:    servicenote.NewCreateNoteRequestHandler(options.NoteSaver),
		GetNoteRequestHandler:       servicenote.NewGetNoteRequestHandler(options.NoteFinder),
		GetNotesRequestHandler:      servicenote.NewGetNotesRequestHandler(options.NoteFinder),
		UpdateNoteRequestHandler:    servicenote.NewUpdateNoteRequestHandler(options.NoteUpdater),
		DeleteNoteRequestHandler:    servicenote.NewDeleteNoteRequestHandler(options.NoteDeleter),
		GetNotesAsyncRequestHandler: servicenote.NewGetNotesAsyncRequestHandler(options.NoteFinder),
	}
}
