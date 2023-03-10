package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPostgreSQL(t *testing.T) {
	t.Run("should create new postgres client", func(t *testing.T) {
		db, err := NewPostgreSQL(context.Background(), Config{
			Host:     "localhost",
			Port:     "5432",
			Username: "postgres",
			Password: "postgres",
			DBName:   "postgres",
			SSLMode:  "disable",
		})
		assert.NoError(t, err)
		assert.NotNil(t, db)
	})

	t.Run("should return error if context is invalid", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		db, err := NewPostgreSQL(ctx, Config{
			Host:     "localhost",
			Port:     "5432",
			Username: "postgres",
			Password: "postgres",
			DBName:   "postgres",
			SSLMode:  "disable",
		})
		assert.ErrorIs(t, err, context.Canceled)
		assert.Nil(t, db)
	})

	t.Run("should return error if config is invalid", func(t *testing.T) {
		db, err := NewPostgreSQL(context.Background(), Config{
			Host:     "invalid-host",
			Port:     "5432",
			Username: "postgres",
			Password: "postgres",
			DBName:   "postgres",
			SSLMode:  "disable",
		})
		assert.Error(t, err)
		assert.Nil(t, db)
	})
}
