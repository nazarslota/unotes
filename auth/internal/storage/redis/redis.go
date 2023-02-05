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

func NewRedis(ctx context.Context, config Config) (*redis.Client, error) {
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
