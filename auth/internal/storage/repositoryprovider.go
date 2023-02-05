package storage

import (
	redisdriver "github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	"github.com/nazarslota/unotes/auth/internal/domain/refresh"
	"github.com/nazarslota/unotes/auth/internal/domain/user"
	"github.com/nazarslota/unotes/auth/internal/storage/postgres"
	"github.com/nazarslota/unotes/auth/internal/storage/redis"
)

type RepositoryProvider struct {
	UserRepository         user.Repository
	RefreshTokenRepository refresh.Repository
}

type RepositoryProviderOption func(rp *RepositoryProvider)

func NewRepositoryProvider(options ...RepositoryProviderOption) *RepositoryProvider {
	s := new(RepositoryProvider)
	for _, option := range options {
		option(s)
	}
	return s
}

func WithPostgreSQLUserRepository(db *sqlx.DB) RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.UserRepository = postgres.NewUserRepository(db)
	}
}

func WithRedisRefreshTokenRepository(client *redisdriver.Client) RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.RefreshTokenRepository = redis.NewRefreshTokenRepository(client)
	}
}
