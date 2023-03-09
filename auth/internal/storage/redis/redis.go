package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

// Config represents the configuration options for the Redis client.
type Config struct {
	// Addr is the Redis server address (e.g. "localhost:6379").
	Addr string

	// Password is the Redis server password (leave blank if none).
	Password string

	// DB is the Redis database number to use (default is 0).
	DB int
}

// NewRedis creates a new Redis client with the given configuration options.
//
// If the context is canceled before the client is created, an error is returned.
//
// If "PING" command fails returns an error.
//
// When finished with the Redis client, it is important to close it to release any resources it may be holding.
// This can be done by calling the Close method on the client.
func NewRedis(ctx context.Context, config Config) (*redis.Client, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("invalid context: %w", ctx.Err())
	default:
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis client: %w", err)
	}
	return client, nil
}
