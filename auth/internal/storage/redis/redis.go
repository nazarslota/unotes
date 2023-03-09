// Package redis provides a Redis repository implementation for storing and managing refresh tokens.
package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

// Config contains the Redis server address, password, and database index.
type Config struct {
	Addr     string
	Password string
	DB       int
}

// NewRedis creates a new Redis client using the provided context and configuration, and returns an error if the client
// fails to connect to the Redis server.
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
