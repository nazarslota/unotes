package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMongoDB(t *testing.T) {
	database, err := NewMongoDB(context.Background(), Config{
		Host:     "localhost",
		Port:     "27017",
		Database: "testdb",
	})
	assert.NoError(t, err)

	err = database.Drop(context.Background())
	require.NoError(t, err)

	database, err = NewMongoDB(context.Background(), Config{
		Host:     "",
		Port:     "",
		Database: "testdb",
	})
	assert.Error(t, err)

	database, err = NewMongoDB(context.Background(), Config{
		Host:     "localhost",
		Port:     "-1",
		Database: "testdb",
	})
	assert.Error(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	database, err = NewMongoDB(ctx, Config{
		Host:     "localhost",
		Port:     "27017",
		Database: "testdb",
	})
	assert.ErrorIs(t, err, context.Canceled)
}
