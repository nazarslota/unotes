// Package storage provides storage implementations.
package storage

import (
	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"

	domainrefresh "github.com/nazarslota/unotes/auth/internal/domain/refresh"
	domainuser "github.com/nazarslota/unotes/auth/internal/domain/user"

	storagepostgres "github.com/nazarslota/unotes/auth/internal/storage/postgres"
	storageredis "github.com/nazarslota/unotes/auth/internal/storage/redis"
)

// RepositoryProvider is a provider for the UserRepository and RefreshTokenRepository.
type RepositoryProvider struct {
	UserRepository         domainuser.Repository
	RefreshTokenRepository domainrefresh.Repository
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

// WithPostgreSQLUserRepository is a functional option that sets the UserRepository
// of the RepositoryProvider to a new instance of `postgres.UserRepository`.
func WithPostgreSQLUserRepository(db *sqlx.DB) RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.UserRepository, _ = storagepostgres.NewUserRepository(db)
	}
}

// WithRedisRefreshTokenRepository is a functional option that sets the RefreshTokenRepository
// of the RepositoryProvider to a new instance of `redis.RefreshTokenRepository`.
func WithRedisRefreshTokenRepository(db *redis.Client) RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.RefreshTokenRepository, _ = storageredis.NewRefreshTokenRepository(db)
	}
}
