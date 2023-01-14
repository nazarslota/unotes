package storage

import (
	redisdriver "github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	"github.com/nazarslota/unotes/auth/internal/domain/refreshtoken"
	"github.com/nazarslota/unotes/auth/internal/domain/user"
	"github.com/nazarslota/unotes/auth/internal/storage/memory"
	"github.com/nazarslota/unotes/auth/internal/storage/mongodb"
	"github.com/nazarslota/unotes/auth/internal/storage/postgresql"
	"github.com/nazarslota/unotes/auth/internal/storage/redis"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

// RepositoryProvider is a struct that holds various repository implementations.
type RepositoryProvider struct {
	// UserRepository is a repository for storing and retrieving users.
	UserRepository user.Repository
	// RefreshTokenRepository is a repository for storing and retrieving refresh tokens.
	RefreshTokenRepository refreshtoken.Repository
}

// RepositoryProviderOption is a function type that can be used to configure a RepositoryProvider.
type RepositoryProviderOption func(rp *RepositoryProvider)

// NewRepositoryProvider creates and returns a new RepositoryProvider with the given options.
func NewRepositoryProvider(options ...RepositoryProviderOption) *RepositoryProvider {
	s := &RepositoryProvider{}
	for _, option := range options {
		option(s)
	}
	return s
}

// WithMemoryUserRepository is a RepositoryProviderOption that sets the UserRepository field to a new memory-based user.Repository.
func WithMemoryUserRepository() RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.UserRepository = memory.NewUserRepository()
	}
}

// WithMongoDBUserRepository is a RepositoryProviderOption that sets the UserRepository field to a new MongoDB-based user.Repository.
func WithMongoDBUserRepository(db *mongodriver.Database) RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.UserRepository = mongodb.NewUserRepository(db)
	}
}

// WithPostgreSQLUserRepository is a RepositoryProviderOption that sets the UserRepository field to a new PostgreSQL-based user.Repository.
func WithPostgreSQLUserRepository(db *sqlx.DB) RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.UserRepository = postgresql.NewUserRepository(db)
	}
}

// WithMemoryRefreshTokenRepository is a RepositoryProviderOption that sets the RefreshTokenRepository field to a new memory-based refreshtoken.Repository.
func WithMemoryRefreshTokenRepository() RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.RefreshTokenRepository = memory.NewRefreshTokenRepository()
	}
}

// WithRedisRefreshTokenRepository is a RepositoryProviderOption that sets the RefreshTokenRepository field to a new Redis-based refreshtoken.Repository.
func WithRedisRefreshTokenRepository(client *redisdriver.Client) RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.RefreshTokenRepository = redis.NewRefreshTokenRepository(client)
	}
}
