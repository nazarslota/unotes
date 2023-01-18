package mongo

import (
	"context"
	"errors"
	"fmt"

	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NoteRepository struct {
	notes *mongo.Collection
}

var _ domainnote.Repository = (*NoteRepository)(nil)

func NewNoteRepository(database *mongo.Database, collection ...string) *NoteRepository {
	if len(collection) > 0 {
		return &NoteRepository{notes: database.Collection(collection[0])}
	}
	return &NoteRepository{notes: database.Collection("notes")}
}

func (r NoteRepository) SaveOne(ctx context.Context, note *domainnote.Note) error {
	res := r.notes.FindOne(ctx, bson.M{"_id": note.ID})
	if err := res.Err(); err == nil {
		return fmt.Errorf("saving note failed: %w", domainnote.ErrAlreadyExist)
	} else if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("saving note failed: %w", err)
	}

	if _, err := r.notes.InsertOne(ctx, note); err != nil {
		return fmt.Errorf("saving note failed: %w", err)
	}
	return nil
}

func (r NoteRepository) FindOne(ctx context.Context, noteID string) (*domainnote.Note, error) {
	res := r.notes.FindOne(ctx, bson.M{"_id": noteID})
	if err := res.Err(); errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("finding note failed: %w", domainnote.ErrNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("finding note failed: %w", err)
	}

	note := new(domainnote.Note)
	if err := res.Decode(&note); err != nil {
		return nil, fmt.Errorf("finding note failed: %w", err)
	}
	return note, nil
}

func (r NoteRepository) FindMany(ctx context.Context, userID string) ([]*domainnote.Note, error) {
	cursor, err := r.notes.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("finding notes failed: %w", err)
	}

	notes := make([]*domainnote.Note, 0)
	for cursor.Next(ctx) {
		note := new(domainnote.Note)
		if err := cursor.Decode(&note); err != nil {
			return nil, fmt.Errorf("finding notes failed: %w", err)
		}
		notes = append(notes, note)
	}

	if len(notes) == 0 {
		return nil, fmt.Errorf("finding notes failed: %w", domainnote.ErrNotFound)
	}
	return notes, nil
}

func (r NoteRepository) DeleteOne(ctx context.Context, noteID string) error {
	err := r.notes.FindOne(ctx, bson.M{"_id": noteID}).Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("deleting note failed: %w", domainnote.ErrNotFound)
	} else if err != nil {
		return fmt.Errorf("deleting note failed: %w", err)
	}

	if _, err := r.notes.DeleteOne(ctx, bson.M{"_id": noteID}); err != nil {
		return fmt.Errorf("deleting note failed: %w", err)
	}
	return nil
}