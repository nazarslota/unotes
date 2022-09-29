package storage

import (
	redisdriver "github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	"github.com/udholdenhed/unotes/auth-service/internal/domain/refreshtoken"
	"github.com/udholdenhed/unotes/auth-service/internal/domain/user"
	"github.com/udholdenhed/unotes/auth-service/internal/storage/memory"
	"github.com/udholdenhed/unotes/auth-service/internal/storage/mongo"
	"github.com/udholdenhed/unotes/auth-service/internal/storage/postgres"
	"github.com/udholdenhed/unotes/auth-service/internal/storage/redis"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type RepositoryProvider struct {
	UserRepository         user.Repository
	RefreshTokenRepository refreshtoken.Repository
}

type RepositoryProviderOption func(rp *RepositoryProvider)

func NewRepositoryProvider(options ...RepositoryProviderOption) *RepositoryProvider {
	s := &RepositoryProvider{}
	for _, option := range options {
		option(s)
	}
	return s
}

// user.Repository

func WithMemoryUserRepository() RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.UserRepository = memory.NewUserRepository()
	}
}

func WithMongoDBUserRepository(db *mongodriver.Database) RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.UserRepository = mongo.NewUserRepository(db)
	}
}

func WithPostgreSQLUserRepository(db *sqlx.DB) RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.UserRepository = postgres.NewUserRepository(db)
	}
}

// refreshtoken.Repository

func WithMemoryRefreshTokenRepository() RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.RefreshTokenRepository = memory.NewRefreshTokenRepository()
	}
}

func WithRedisRefreshTokenRepository(client *redisdriver.Client) RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.RefreshTokenRepository = redis.NewRefreshTokenRepository(client)
	}
}
