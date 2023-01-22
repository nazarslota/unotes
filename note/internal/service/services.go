package service

type Services struct {
	NoteService NoteService
}

func NewServices(noteServiceOptions NoteServiceOptions) Services {
	return Services{
		NoteService: NewNoteService(noteServiceOptions),
	}
}
