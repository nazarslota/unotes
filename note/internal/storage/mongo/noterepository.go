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

func (r NoteRepository) SaveOne(ctx context.Context, note domainnote.Note) error {
	if _, err := r.notes.InsertOne(ctx, note); mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("saving note failed: %w", domainnote.ErrNoteAlreadyExist)
	} else if err != nil {
		return fmt.Errorf("saving note failed: %w", err)
	}
	return nil
}

func (r NoteRepository) FindOne(ctx context.Context, noteID string) (domainnote.Note, error) {
	res := r.notes.FindOne(ctx, bson.M{"_id": noteID})
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domainnote.Note{}, fmt.Errorf("finding note failed: %w", domainnote.ErrNoteNotFound)
		}
		return domainnote.Note{}, fmt.Errorf("finding note failed: %w", err)
	}

	var note domainnote.Note
	if err := res.Decode(&note); err != nil {
		return domainnote.Note{}, fmt.Errorf("finding note failed: %w", err)
	}
	return note, nil
}

func (r NoteRepository) FindMany(ctx context.Context, userID string) ([]domainnote.Note, error) {
	cursor, err := r.notes.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("finding notes failed: %w", err)
	}
	defer func() { _ = cursor.Close(ctx) }()

	notes := make([]domainnote.Note, 0, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		var note domainnote.Note
		if err := cursor.Decode(&note); err != nil {
			return nil, fmt.Errorf("finding notes failed: %w", err)
		}
		notes = append(notes, note)
	}

	if len(notes) == 0 {
		return nil, fmt.Errorf("finding notes failed: %w", domainnote.ErrNoteNotFound)
	}
	return notes, nil
}

func (r NoteRepository) UpdateOne(ctx context.Context, note domainnote.Note) error {
	update := bson.M{"$set": bson.M{
		"title":   note.Title,
		"content": note.Content,
	}}

	if result, err := r.notes.UpdateByID(ctx, note.ID, update); err != nil {
		return fmt.Errorf("updating note failed: %w", err)
	} else if result.MatchedCount == 0 {
		return fmt.Errorf("updating note failed: %w", domainnote.ErrNoteNotFound)
	}
	return nil
}

func (r NoteRepository) DeleteOne(ctx context.Context, noteID string) error {
	if result, err := r.notes.DeleteOne(ctx, bson.M{"_id": noteID}); err != nil {
		return fmt.Errorf("deleting note failed: %w", err)
	} else if result.DeletedCount == 0 {
		return fmt.Errorf("deleting note failed: %w", domainnote.ErrNoteNotFound)
	}
	return nil
}

func (r NoteRepository) FindManyAsync(ctx context.Context, userID string) (<-chan domainnote.Note, <-chan error) {
	nts := make(chan domainnote.Note)
	errs := make(chan error)

	go func() {
		defer close(nts)
		defer close(errs)

		cursor, err := r.notes.Find(ctx, bson.M{"user_id": userID})
		if err != nil {
			errs <- fmt.Errorf("finding nts failed: %w", err)
			return
		}
		defer func() { _ = cursor.Close(ctx) }()

		for cursor.Next(ctx) {
			var note domainnote.Note
			if err := cursor.Decode(&note); err != nil {
				errs <- fmt.Errorf("finding nts failed: %w", err)
				return
			}
			nts <- note
		}
	}()
	return nts, errs
}
