// Package redis provides a Redis repository implementation for storing and managing refresh tokens.
package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

// Config represents the configuration options for the Redis db.
type Config struct {
	Addr     string // Addr is the Redis server address (e.g. "localhost:6379").
	Password string // Password is the Redis server password (leave blank if none).
	DB       int    // DB is the Redis database number to use (default is 0).
}

// NewRedis creates a new Redis db with the given configuration options.
//
// If the context is canceled before the db is created, an error is returned.
// If "PING" command fails returns an error.
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
		return nil, fmt.Errorf("failed to ping redis db: %w", err)
	}
	return client, nil
}
