package service

type Services struct {
	NoteService NoteService
	JWTService  JWTService
}

func NewServices(jwtServiceOptions JWTServiceOptions, noteServiceOptions NoteServiceOptions) Services {
	return Services{
		JWTService:  NewJWTService(jwtServiceOptions),
		NoteService: NewNoteService(noteServiceOptions),
	}
}
