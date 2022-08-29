package storage

import (
	"github.com/udholdenhed/unotes/auth/internal/domain/refreshtoken"
	"github.com/udholdenhed/unotes/auth/internal/domain/user"
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
