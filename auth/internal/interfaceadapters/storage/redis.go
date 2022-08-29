package storage

import (
	redisdriver "github.com/go-redis/redis/v9"
	"github.com/udholdenhed/unotes/auth/internal/interfaceadapters/storage/redis"
)

func WithRedisRefreshTokenRepository(client *redisdriver.Client) RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.RefreshTokenRepository = redis.NewRefreshTokenRepository(client)
	}
}
