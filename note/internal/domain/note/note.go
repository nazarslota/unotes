package note

import (
	"errors"
	"time"
)

type Note struct {
	ID             string     `json:"id" bson:"_id"`
	Title          string     `json:"title" bson:"title"`
	Content        string     `json:"content" bson:"content"`
	UserID         string     `json:"user_id" bson:"user_id"`
	CreatedAt      time.Time  `json:"created_at" bson:"created_at"`
	Priority       *string    `json:"priority,omitempty" bson:"priority,omitempty"`
	CompletionTime *time.Time `json:"completion_time,omitempty" bson:"completion_time,omitempty"`
}

var (
	ErrNoteAlreadyExist = errors.New("already exist")
	ErrNoteNotFound     = errors.New("not found")
)
