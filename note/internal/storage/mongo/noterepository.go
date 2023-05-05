package mongo

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	domain "github.com/nazarslota/unotes/note/internal/domain/note"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// NoteRepository is a struct that provides methods for interacting with the MongoDB database.
type NoteRepository struct {
	collection *mongo.Collection
}

// NewNoteRepository creates a new NoteRepository instance with a MongoDB collection.
func NewNoteRepository(db *mongo.Database, collection ...string) (*NoteRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	if len(collection) > 0 {
		return &NoteRepository{collection: db.Collection(collection[0])}, nil
	}
	return &NoteRepository{collection: db.Collection("notes")}, nil
}

// SaveOne saves a note to the MongoDB collection.
// If a note with the same ID already exists, returns an error.
func (r NoteRepository) SaveOne(ctx context.Context, note domain.Note) error {
	if _, err := r.collection.InsertOne(ctx, note); mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("saving note failed: %w", domain.ErrNoteAlreadyExist)
	} else if err != nil {
		return fmt.Errorf("saving note failed: %w", err)
	}
	return nil
}

// FindOne finds a note with a specific ID in the MongoDB collection.
// If no note is found, returns an error.
func (r NoteRepository) FindOne(ctx context.Context, noteID string) (domain.Note, error) {
	res := r.collection.FindOne(ctx, bson.M{"_id": noteID})
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Note{}, fmt.Errorf("finding note failed: %w", domain.ErrNoteNotFound)
		}
		return domain.Note{}, fmt.Errorf("finding note failed: %w", err)
	}

	var note domain.Note
	if err := res.Decode(&note); err != nil {
		return domain.Note{}, fmt.Errorf("finding note failed: %w", err)
	}
	return note, nil
}

// FindMany finds all notes associated with a specific user in the MongoDB collection.
// If no notes are found, returns an error.
func (r NoteRepository) FindMany(ctx context.Context, userID string) ([]domain.Note, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("finding collection failed: %w", err)
	}
	defer func() { _ = cursor.Close(ctx) }()

	notes := make([]domain.Note, 0, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		var note domain.Note
		if err := cursor.Decode(&note); err != nil {
			return nil, fmt.Errorf("finding collection failed: %w", err)
		}
		notes = append(notes, note)
	}

	if len(notes) == 0 {
		return nil, fmt.Errorf("finding collection failed: %w", domain.ErrNoteNotFound)
	}
	return notes, nil
}

// UpdateOne updates a note in the MongoDB collection.
// If no note with the specified ID is found, returns an error.
func (r NoteRepository) UpdateOne(ctx context.Context, note domain.Note) error {
	update := bson.M{"$set": bson.M{
		"title":   note.Title,
		"content": note.Content,
	}}

	if result, err := r.collection.UpdateByID(ctx, note.ID, update); err != nil {
		return fmt.Errorf("updating note failed: %w", err)
	} else if result.MatchedCount == 0 {
		return fmt.Errorf("updating note failed: %w", domain.ErrNoteNotFound)
	}
	return nil
}

// DeleteOne deletes a note from the MongoDB collection.
// If no note with the specified ID is found, returns an error.
func (r NoteRepository) DeleteOne(ctx context.Context, noteID string) error {
	if result, err := r.collection.DeleteOne(ctx, bson.M{"_id": noteID}); err != nil {
		return fmt.Errorf("deleting note failed: %w", err)
	} else if result.DeletedCount == 0 {
		return fmt.Errorf("deleting note failed: %w", domain.ErrNoteNotFound)
	}
	return nil
}

func (r NoteRepository) FindManyAsync(ctx context.Context, userID string) (<-chan domain.Note, <-chan error) { // TODO: Documentation.
	notes, errs := make(chan domain.Note), make(chan error)
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		go func() {
			errs <- fmt.Errorf("finding collection failed: %w", err)
			close(errs)
		}()
		close(notes)
		return notes, errs
	}
	go r.notesFromCursor(ctx, cursor, notes, errs)
	return notes, errs
}

func (r NoteRepository) notesFromCursor(ctx context.Context, cursor *mongo.Cursor, notes chan<- domain.Note, errs chan<- error) {
	var found atomic.Bool
	var wgNotes, wgErrs sync.WaitGroup
	for cursor.Next(ctx) {
		var note domain.Note
		err := cursor.Decode(&note)
		if err != nil {
			wgErrs.Add(1)
			go func() {
				errs <- fmt.Errorf("finding collection failed: %w", err)
				wgErrs.Done()
			}()
		} else {
			wgNotes.Add(1)
			if !found.Load() {
				found.Store(true)
			}

			go func() {
				notes <- note
				wgNotes.Done()
			}()
		}
	}

	if !found.Load() {
		wgErrs.Add(1)
		go func() {
			errs <- fmt.Errorf("finding collection failed: %w", domain.ErrNoteNotFound)
			wgErrs.Done()
		}()
	}

	go func() {
		wgNotes.Wait()
		close(notes)
	}()
	go func() {
		wgErrs.Wait()
		close(errs)
	}()
}
