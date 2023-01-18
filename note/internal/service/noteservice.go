package service

import (
	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
	servicenote "github.com/nazarslota/unotes/note/internal/service/note"
)

type NoteService struct {
	CreateNoteRequestHandler servicenote.CreateNoteRequestHandler
	DeleteNoteRequestHandler servicenote.DeleteNoteRequestHandler
	GetNoteRequestHandler    servicenote.GetNoteRequestHandler
	GetNotesRequestHandler   servicenote.GetNotesRequestHandler
}

type NoteServiceOptions struct {
	NoteRepository domainnote.Repository
}

func NewNoteService(options NoteServiceOptions) NoteService {
	return NoteService{
		CreateNoteRequestHandler: servicenote.NewCreateNoteRequestHandler(options.NoteRepository),
		DeleteNoteRequestHandler: servicenote.NewDeleteNoteRequestHandler(options.NoteRepository),
		GetNoteRequestHandler:    servicenote.NewGetNoteRequestHandler(options.NoteRepository),
		GetNotesRequestHandler:   servicenote.NewGetNotesRequestHandler(options.NoteRepository),
	}
}
