package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/udholdenhed/unotes/user-service/internal/domain/user"
	"github.com/udholdenhed/unotes/user-service/internal/storage/memory"
	"github.com/udholdenhed/unotes/user-service/internal/storage/postgres"
)

type RepositoryProvider struct {
	UserRepository user.Repository
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

func WithPostgreSQLUserRepository(db *sqlx.DB) RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.UserRepository = postgres.NewUserRepository(db)
	}
}
