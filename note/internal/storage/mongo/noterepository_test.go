package mongo

import (
	"context"
	"testing"

	domain "github.com/nazarslota/unotes/note/internal/domain/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	noteAA = domain.Note{
		ID:      "note-a-id",
		Title:   "note-a-title",
		Content: "note-a-content",
		UserID:  "user-a-id",
	}

	noteAB = domain.Note{
		ID:      "note-b-id",
		Title:   "note-b-title",
		Content: "note-b-content",
		UserID:  "user-a-id",
	}
)

var repository *NoteRepository

func init() {
	db, err := NewMongoDB(context.Background(), Config{
		Host:     "localhost",
		Port:     "27017",
		Username: "",
		Password: "",
		Database: "test",
	})
	if err != nil {
		panic(err)
	}

	repository, err = NewNoteRepository(db, "test")
	if err != nil {
		panic(err)
	}
}

func TestNewNoteRepository(t *testing.T) {
	t.Run("should create new note repository", func(t *testing.T) {
		db, err := NewMongoDB(context.Background(), Config{
			Host:     "localhost",
			Port:     "27017",
			Username: "",
			Password: "",
			Database: "test",
		})
		require.NoError(t, err)
		require.NotNil(t, db)

		repository, err := NewNoteRepository(db, "test")
		assert.NoError(t, err)
		assert.NotNil(t, repository)
	})

	t.Run("should return error if db is nil", func(t *testing.T) {
		repository, err := NewNoteRepository(nil, "test")
		assert.Error(t, err)
		assert.Nil(t, repository)
	})
}

func TestNoteRepository_SaveOne(t *testing.T) {
	t.Run("should successfully save a note", func(t *testing.T) {
		err := repository.SaveOne(context.Background(), noteAA)
		assert.NoError(t, err)

		filter := bson.M{"_id": noteAA.ID}
		result := repository.collection.FindOne(context.Background(), filter)

		var note domain.Note
		err = result.Decode(&note)
		assert.NoError(t, err)
		assert.Equal(t, noteAA, note)

		t.Cleanup(func() {
			_ = repository.collection.Drop(context.Background())
		})
	})

	t.Run("should return an error if context is invalid", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := repository.SaveOne(ctx, noteAA)
		assert.ErrorIs(t, err, context.Canceled)

		filter := bson.M{"_id": noteAA.ID}
		err = repository.collection.FindOne(context.Background(), filter).Err()
		require.ErrorIs(t, err, mongo.ErrNoDocuments)

		t.Cleanup(func() {
			_ = repository.collection.Drop(context.Background())
		})
	})

	t.Run("should return an error if note already exists", func(t *testing.T) {
		_, err := repository.collection.InsertOne(context.Background(), noteAA)
		require.NoError(t, err)

		err = repository.SaveOne(context.Background(), noteAA)
		assert.ErrorIs(t, err, domain.ErrNoteAlreadyExist)

		t.Cleanup(func() {
			_ = repository.collection.Drop(context.Background())
		})
	})

	t.Cleanup(func() {
		_ = repository.collection.Database().Drop(context.Background())
	})
}

func TestNoteRepository_FindOne(t *testing.T) {
	t.Run("should return note", func(t *testing.T) {
		_, err := repository.collection.InsertOne(context.Background(), noteAA)
		require.NoError(t, err)

		result, err := repository.FindOne(context.Background(), noteAA.ID)
		assert.NoError(t, err)
		assert.Equal(t, noteAA, result)

		t.Cleanup(func() {
			_ = repository.collection.Drop(context.Background())
		})
	})

	t.Run("should return an error if context is invalid", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		result, err := repository.FindOne(ctx, noteAA.ID)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Empty(t, result)

		t.Cleanup(func() {
			_ = repository.collection.Drop(context.Background())
		})
	})

	t.Run("should return an error if note does not exist", func(t *testing.T) {
		result, err := repository.FindOne(context.Background(), "invalid-note-id")
		assert.Error(t, err)
		assert.Empty(t, result)

		t.Cleanup(func() {
			_ = repository.collection.Drop(context.Background())
		})
	})

	t.Cleanup(func() {
		_ = repository.collection.Database().Drop(context.Background())
	})
}

func TestNoteRepository_FindMany(t *testing.T) {
	t.Run("should return notes", func(t *testing.T) {
		_, err := repository.collection.InsertMany(context.Background(), []any{noteAA, noteAB})
		require.NoError(t, err)

		result, err := repository.FindMany(context.Background(), noteAA.UserID)
		assert.NoError(t, err)
		assert.Contains(t, result, noteAA)
		assert.Contains(t, result, noteAB)

		t.Cleanup(func() {
			_ = repository.collection.Drop(context.Background())
		})
	})

	t.Run("should return an error if context is invalid", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := repository.collection.InsertMany(context.Background(), []any{noteAA, noteAB})
		require.NoError(t, err)

		result, err := repository.FindMany(ctx, noteAA.ID)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Nil(t, result)

		t.Cleanup(func() {
			_ = repository.collection.Drop(context.Background())
		})
	})

	t.Run("should return an error if notes does not exist", func(t *testing.T) {
		result, err := repository.FindMany(context.Background(), "invalid-user-id")
		assert.ErrorIs(t, err, domain.ErrNoteNotFound)
		assert.Nil(t, result)

		t.Cleanup(func() {
			_ = repository.collection.Drop(context.Background())
		})
	})

	t.Cleanup(func() {
		_ = repository.collection.Database().Drop(context.Background())
	})
}

func TestNoteRepository_UpdateOne(t *testing.T) {
	t.Run("should update note", func(t *testing.T) {
		_, err := repository.collection.InsertOne(context.Background(), noteAA)
		require.NoError(t, err)

		updated := noteAA
		updated.Content = "updated-note-content"

		err = repository.UpdateOne(context.Background(), updated)
		assert.NoError(t, err)

		filter := bson.M{"_id": noteAA.ID}
		result := repository.collection.FindOne(context.Background(), filter)

		var note domain.Note
		err = result.Decode(&note)
		assert.NoError(t, err)
		assert.Equal(t, updated, note)

		t.Cleanup(func() {
			_ = repository.collection.Drop(context.Background())
		})
	})

	t.Run("should return an error if context is invalid", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := repository.collection.InsertOne(context.Background(), noteAA)
		require.NoError(t, err)

		updated := noteAA
		updated.Content = "updated-note-content"

		err = repository.UpdateOne(ctx, updated)
		assert.ErrorIs(t, err, context.Canceled)

		filter := bson.M{"_id": noteAA.ID}
		result := repository.collection.FindOne(context.Background(), filter)

		var note domain.Note
		err = result.Decode(&note)
		assert.NoError(t, err)
		assert.NotEqual(t, updated, noteAA)

		t.Cleanup(func() {
			_ = repository.collection.Drop(context.Background())
		})
	})

	t.Run("should return en error if note does not exist", func(t *testing.T) {
		err := repository.UpdateOne(context.Background(), noteAA)
		assert.ErrorIs(t, err, domain.ErrNoteNotFound)

		filter := bson.M{"_id": noteAA.ID}
		err = repository.collection.FindOne(context.Background(), filter).Err()
		require.ErrorIs(t, err, mongo.ErrNoDocuments)
	})
}

func TestNoteRepository_DeleteOne(t *testing.T) {
	t.Run("should successfully delete note", func(t *testing.T) {
		_, err := repository.collection.InsertOne(context.Background(), noteAA)
		require.NoError(t, err)

		err = repository.DeleteOne(context.Background(), noteAA.ID)
		assert.NoError(t, err)

		t.Cleanup(func() {
			_ = repository.collection.Drop(context.Background())
		})
	})

	t.Run("should return an error if context is invalid", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := repository.collection.InsertOne(context.Background(), noteAA)
		require.NoError(t, err)

		err = repository.DeleteOne(ctx, noteAA.ID)
		assert.ErrorIs(t, err, context.Canceled)

		t.Cleanup(func() {
			_ = repository.collection.Drop(context.Background())
		})
	})

	t.Run("should return an error if note does not exist", func(t *testing.T) {
		err := repository.DeleteOne(context.Background(), "invalid-note-id")
		assert.ErrorIs(t, err, domain.ErrNoteNotFound)
	})
}

func TestNoteRepository_FindManyAsync(t *testing.T) {
	t.Run("should successfully find notes", func(t *testing.T) {
		_, err := repository.collection.InsertMany(context.Background(), []any{noteAA, noteAB})
		require.NoError(t, err)

		notes, errs := repository.FindManyAsync(context.Background(), noteAA.UserID)

		note, ok := <-notes
		assert.True(t, ok)
		assert.Contains(t, []any{noteAA, noteAB}, note)

		note, ok = <-notes
		assert.True(t, ok)
		assert.Contains(t, []any{noteAA, noteAB}, note)

		note, ok = <-notes
		assert.False(t, ok)
		assert.Empty(t, note)

		err, ok = <-errs
		assert.False(t, ok)
		assert.NoError(t, err)

		t.Cleanup(func() {
			_ = repository.collection.Drop(context.Background())
		})
	})

	t.Run("should return an error if context is invalid", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := repository.collection.InsertMany(context.Background(), []any{noteAA, noteAB})
		require.NoError(t, err)

		notes, errs := repository.FindManyAsync(ctx, "invalid-user-id")
		note, ok := <-notes
		assert.False(t, ok)
		assert.Empty(t, note)

		err, ok = <-errs
		assert.True(t, ok)
		assert.ErrorIs(t, err, context.Canceled)

		err, ok = <-errs
		assert.False(t, ok)
		assert.NoError(t, err)

		t.Cleanup(func() {
			_ = repository.collection.Drop(context.Background())
		})
	})

	t.Run("should return an error if note does not exist", func(t *testing.T) {
		notes, errs := repository.FindManyAsync(context.Background(), "invalid-user-id")

		note, ok := <-notes
		assert.False(t, ok)
		assert.Empty(t, note)

		err, ok := <-errs
		assert.True(t, ok)
		assert.ErrorIs(t, err, domain.ErrNoteNotFound)

		t.Cleanup(func() {
			_ = repository.collection.Drop(context.Background())
		})
	})
}
