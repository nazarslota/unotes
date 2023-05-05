package note

import "errors"

type Note struct {
	ID      string `bson:"_id"`
	Title   string `bson:"title"`
	Content string `bson:"content"`
	UserID  string `bson:"user_id"`
}

var (
	ErrNoteAlreadyExist = errors.New("already exist")
	ErrNoteNotFound     = errors.New("not found")
)
