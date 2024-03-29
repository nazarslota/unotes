// Package storage provides storage implementations.
package storage

import (
	storagemongo "github.com/nazarslota/unotes/note/internal/storage/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

// RepositoryProvider is a provider for the note repository.
type RepositoryProvider struct {
	MongoNoteRepository *storagemongo.NoteRepository
}

// RepositoryProviderOption is a functional option for the RepositoryProvider.
type RepositoryProviderOption func(rp *RepositoryProvider)

// NewRepositoryProvider creates a new instance of the RepositoryProvider.
// It takes one or more options that can be used to configure the provider.
func NewRepositoryProvider(options ...RepositoryProviderOption) *RepositoryProvider {
	rp := &RepositoryProvider{}
	for _, option := range options {
		option(rp)
	}
	return rp
}

// WithMongoNoteRepository is a functional option that sets the MongoNoteRepository
// of the RepositoryProvider to a new instance of `mongo.NoteRepository`.
func WithMongoNoteRepository(db *mongo.Database) RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.MongoNoteRepository, _ = storagemongo.NewNoteRepository(db)
	}
}
