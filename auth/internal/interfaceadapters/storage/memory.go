package storage

import "github.com/udholdenhed/unotes/auth/internal/interfaceadapters/storage/memory"

func WithMemoryUserRepository() RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.UserRepository = memory.NewUserRepository()
	}
}

func WithMemoryRefreshTokenRepository() RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.RefreshTokenRepository = memory.NewRefreshTokenRepository()
	}
}
