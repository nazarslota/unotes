package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisClient(ctx context.Context, config *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("filed to ping redis: %w", err)
	}

	return client, nil
}
