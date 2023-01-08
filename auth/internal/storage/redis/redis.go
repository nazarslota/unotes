package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

// Config is a struct that holds the configuration for a Redis client.
type Config struct {
	// Addr is the address of the Redis server.
	Addr string
	// Password is the password to use when connecting to the Redis server.
	Password string
	// DB is the Redis database to use.
	DB int
}

// NewRedis creates and returns a new Redis client with the given configuration.
// If the connection to the Redis server fails, an error is returned.
func NewRedis(ctx context.Context, config *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return client, nil
}
