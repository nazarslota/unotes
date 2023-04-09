// Package storage provides storage implementations.
package storage

import (
	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"

	storagepostgres "github.com/nazarslota/unotes/auth/internal/storage/postgres"
	storageredis "github.com/nazarslota/unotes/auth/internal/storage/redis"
)

// RepositoryProvider is a provider for the PostgresUserRepository and RedisRefreshTokenRepository.
type RepositoryProvider struct {
	PostgresUserRepository      *storagepostgres.UserRepository
	RedisRefreshTokenRepository *storageredis.RefreshTokenRepository
}

// RepositoryProviderOption is a functional option for the RepositoryProvider.
type RepositoryProviderOption func(rp *RepositoryProvider)

// NewRepositoryProvider creates a new instance of the RepositoryProvider.
// It takes one or more options that can be used to configure the provider.
func NewRepositoryProvider(options ...RepositoryProviderOption) *RepositoryProvider {
	s := new(RepositoryProvider)
	for _, option := range options {
		option(s)
	}
	return s
}

// WithPostgreSQLUserRepository is a functional option that sets the PostgresUserRepository
// of the RepositoryProvider to a new instance of `postgres.UserRepository`.
func WithPostgreSQLUserRepository(db *sqlx.DB) RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.PostgresUserRepository, _ = storagepostgres.NewUserRepository(db)
	}
}

// WithRedisRefreshTokenRepository is a functional option that sets the RedisRefreshTokenRepository
// of the RepositoryProvider to a new instance of `redis.RefreshTokenRepository`.
func WithRedisRefreshTokenRepository(db *redis.Client) RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.RedisRefreshTokenRepository, _ = storageredis.NewRefreshTokenRepository(db)
	}
}
