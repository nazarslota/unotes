package storage

import (
	"context"
	"testing"

	storagememory "github.com/nazarslota/unotes/note/internal/storage/memory"
	storagemongo "github.com/nazarslota/unotes/note/internal/storage/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestNewRepositoryProvider(t *testing.T) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)
	defer func() { _ = client.Disconnect(context.Background()) }()

	database := client.Database("test")

	// Test with memory note repository
	rp := NewRepositoryProvider(WithMemoryNoteRepository())
	assert.NotNil(t, rp.NoteRepository)
	_, ok := rp.NoteRepository.(*storagememory.NoteRepository)
	assert.True(t, ok)

	// Test with MongoDB note repository
	rp = NewRepositoryProvider(WithMongoNoteRepository(database))
	assert.NotNil(t, rp.NoteRepository)
	_, ok = rp.NoteRepository.(*storagemongo.NoteRepository)
	assert.True(t, ok)
}
