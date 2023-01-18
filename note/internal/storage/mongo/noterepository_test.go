package mongo

import (
	"context"
	"testing"

	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestNewNoteRepository(t *testing.T) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)

	defer func() { _ = client.Disconnect(context.Background()) }()

	database := client.Database("test")
	defer func() { _ = database.Drop(context.Background()) }()

	notes := database.Collection("notes")
	defer func() { _ = notes.Drop(context.Background()) }()

	repository := NewNoteRepository(database)
	assert.NotNil(t, repository)

	repository = NewNoteRepository(database, "notes")
	assert.NotNil(t, repository)
}

func TestNoteRepository_SaveOne(t *testing.T) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)

	database := client.Database("test")
	defer func() { _ = database.Drop(context.Background()) }()

	notes := database.Collection("notes")
	defer func() { _ = notes.Drop(context.Background()) }()

	repository := NewNoteRepository(database, "notes")
	note := &domainnote.Note{ID: "123", Title: "Test Note", Content: "This is a test note"}

	err = repository.SaveOne(context.Background(), note)
	assert.NoError(t, err)

	result := new(domainnote.Note)
	err = notes.FindOne(context.Background(), bson.M{"_id": note.ID}).Decode(&result)
	if assert.NoError(t, err) {
		assert.Equal(t, note, result)
	}

	err = repository.SaveOne(context.Background(), note)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, domainnote.ErrAlreadyExist)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = repository.SaveOne(ctx, note)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, context.Canceled)
	}
}

func TestNoteRepository_FindOne(t *testing.T) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)

	database := client.Database("test")
	defer func() { _ = database.Drop(context.Background()) }()

	notes := database.Collection("notes")
	defer func() { _ = notes.Drop(context.Background()) }()

	repository := NewNoteRepository(database, "notes")
	note := &domainnote.Note{ID: "123", Title: "Test Note", Content: "This is a test note"}

	_, err = repository.FindOne(context.Background(), note.ID)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, domainnote.ErrNotFound)
	}

	_, err = notes.InsertOne(context.Background(), note)
	require.NoError(t, err)

	result, err := repository.FindOne(context.Background(), note.ID)
	if assert.NoError(t, err) {
		assert.Equal(t, note, result)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err = repository.FindOne(ctx, note.ID)
	if assert.Error(t, err) {
		assert.Nil(t, result)
		assert.ErrorIs(t, err, context.Canceled)
	}
}

func TestNoteRepository_FindMany(t *testing.T) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)

	database := client.Database("test")
	defer func() { _ = database.Drop(context.Background()) }()

	notes := database.Collection("notes")
	defer func() { _ = notes.Drop(context.Background()) }()

	repository := NewNoteRepository(database, "notes")
	note1 := &domainnote.Note{ID: "123", UserID: "user1", Title: "Test Note 1", Content: "This is a test note"}
	note2 := &domainnote.Note{ID: "456", UserID: "user1", Title: "Test Note 2", Content: "This is a test note"}
	note3 := &domainnote.Note{ID: "789", UserID: "user2", Title: "Test Note 3", Content: "This is a test note"}

	result, err := repository.FindMany(context.Background(), "user1")
	if assert.Error(t, err) {
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domainnote.ErrNotFound)
	}

	_, err = notes.InsertMany(context.Background(), []any{note1, note2, note3})
	require.NoError(t, err)

	result, err = repository.FindMany(context.Background(), "user1")
	if assert.NoError(t, err) {
		assert.ElementsMatch(t, []*domainnote.Note{note1, note2}, result)
	}

	result, err = repository.FindMany(context.Background(), "user2")
	if assert.NoError(t, err) {
		assert.ElementsMatch(t, []*domainnote.Note{note3}, result)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err = repository.FindMany(ctx, "user2")
	if assert.Error(t, err) {
		assert.Nil(t, result)
		assert.ErrorIs(t, err, context.Canceled)
	}
}

func TestNoteRepository_DeleteOne(t *testing.T) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)

	defer func() { _ = client.Disconnect(context.Background()) }()

	database := client.Database("test")
	defer func() { _ = database.Drop(context.Background()) }()

	notes := database.Collection("notes")
	defer func() { _ = notes.Drop(context.Background()) }()

	repository := NewNoteRepository(database, "notes")
	note := &domainnote.Note{Title: "Test Note", Content: "This is a test note"}

	err = repository.DeleteOne(context.Background(), note.ID)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, domainnote.ErrNotFound)
	}

	_, err = notes.InsertOne(context.Background(), note)
	require.NoError(t, err)

	err = repository.DeleteOne(context.Background(), note.ID)
	assert.NoError(t, err)

	result := new(domainnote.Note)
	err = notes.FindOne(context.Background(), bson.M{"_id": note.ID}).Decode(&result)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, mongo.ErrNoDocuments)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = repository.DeleteOne(ctx, note.ID)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, context.Canceled)
	}
}