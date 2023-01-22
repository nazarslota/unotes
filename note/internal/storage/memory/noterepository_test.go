package memory

import (
	"context"
	"testing"

	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNoteRepository(t *testing.T) {
	repository := NewNoteRepository()
	assert.NotNil(t, repository)
}

func TestNoteRepository_SaveOne(t *testing.T) {
	repository := NewNoteRepository()
	require.NotNil(t, repository)

	note := &domainnote.Note{ID: "1", Title: "Test Note", Content: "This is a test note"}

	err := repository.SaveOne(context.Background(), note)
	assert.NoError(t, err)

	err = repository.SaveOne(context.Background(), note)
	assert.ErrorIs(t, err, domainnote.ErrAlreadyExist)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = repository.SaveOne(ctx, note)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestNoteRepository_FindOne(t *testing.T) {
	repository := NewNoteRepository()
	require.NotNil(t, repository)

	note := &domainnote.Note{ID: "1", Title: "Test Note", Content: "This is a test note"}

	_, err := repository.FindOne(context.Background(), note.ID)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, domainnote.ErrNotFound)
	}

	err = repository.SaveOne(context.Background(), note)
	require.NoError(t, err)

	result, err := repository.FindOne(context.Background(), note.ID)
	if assert.NoError(t, err) {
		assert.Equal(t, note, result)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err = repository.FindOne(ctx, note.ID)
	if assert.ErrorIs(t, err, context.Canceled) {
		assert.Nil(t, result)
	}
}

func TestNoteRepository_FindMany(t *testing.T) {
	repository := NewNoteRepository()
	require.NotNil(t, repository)

	note1 := &domainnote.Note{ID: "123", UserID: "user1", Title: "Test Note 1", Content: "This is a test note"}
	note2 := &domainnote.Note{ID: "456", UserID: "user1", Title: "Test Note 2", Content: "This is another test note"}
	note3 := &domainnote.Note{ID: "789", UserID: "user2", Title: "Test Note 3", Content: "This is yet another test note"}

	_, err := repository.FindMany(context.Background(), "user1")
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, domainnote.ErrNotFound)
	}

	err = repository.SaveOne(context.Background(), note1)
	require.NoError(t, err)

	err = repository.SaveOne(context.Background(), note2)
	require.NoError(t, err)

	err = repository.SaveOne(context.Background(), note3)
	require.NoError(t, err)

	notes, err := repository.FindMany(context.Background(), "user1")
	if assert.NoError(t, err) {
		assert.ElementsMatch(t, []*domainnote.Note{note1, note2}, notes)
	}

	notes, err = repository.FindMany(context.Background(), "user2")
	if assert.NoError(t, err) {
		assert.ElementsMatch(t, []*domainnote.Note{note3}, notes)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = repository.FindMany(ctx, "user1")
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, context.Canceled)
	}
}

func TestNoteRepository_DeleteOne(t *testing.T) {
	repository := NewNoteRepository()
	require.NotNil(t, repository)

	n := &domainnote.Note{ID: "1", Title: "Test Note", Content: "This is a test note"}

	err := repository.DeleteOne(context.Background(), n.ID)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, domainnote.ErrNotFound)
	}

	err = repository.SaveOne(context.Background(), n)
	require.NoError(t, err)

	err = repository.DeleteOne(context.Background(), n.ID)
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = repository.DeleteOne(ctx, n.ID)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, context.Canceled)
	}
}
