package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRedis(t *testing.T) {
	t.Run("should create a new redis client", func(t *testing.T) {
		redis, err := NewRedis(context.Background(), Config{})
		assert.NoError(t, err)
		assert.NotNil(t, redis)

		err = redis.Ping(context.Background()).Err()
		assert.NoError(t, err)
	})

	t.Run("should return an error if the address is invalid", func(t *testing.T) {
		redis, err := NewRedis(context.Background(), Config{Addr: "invalid:address"})
		assert.Error(t, err)
		assert.Nil(t, redis)
	})

	t.Run("should return an error if the context is invalid", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		redis, err := NewRedis(ctx, Config{})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, context.Canceled)
		}
		assert.Nil(t, redis)

		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		time.Sleep(100 * time.Millisecond)

		redis, err = NewRedis(ctx, Config{})
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, context.DeadlineExceeded)
		}
		assert.Nil(t, redis)
	})
}
