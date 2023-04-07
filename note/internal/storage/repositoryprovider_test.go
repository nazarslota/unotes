package storage

import (
	"context"
	"testing"

	"github.com/nazarslota/unotes/note/internal/storage/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRepositoryProvider(t *testing.T) {
	t.Run("should create a new repository provider", func(t *testing.T) {
		provider := NewRepositoryProvider()
		assert.NotNil(t, provider)
		assert.Nil(t, provider.NoteRepository)
	})

	t.Run("should not create a new repository provider with given options", func(t *testing.T) {
		db, err := mongo.NewMongoDB(context.Background(), mongo.Config{
			Host:     "localhost",
			Port:     "27017",
			Username: "",
			Password: "",
			Database: "test",
		})
		require.NoError(t, err)
		require.NotNil(t, db)

		provider := NewRepositoryProvider(WithMongoNoteRepository(db))
		assert.NotNil(t, provider)
		assert.NotNil(t, provider.NoteRepository)

		t.Cleanup(func() { _ = db.Drop(context.Background()) })
	})
}
