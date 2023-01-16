package storage

import (
	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
	storagememory "github.com/nazarslota/unotes/note/internal/storage/memory"
	storagemongo "github.com/nazarslota/unotes/note/internal/storage/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoryProvider struct {
	NoteRepository domainnote.Repository
}

type RepositoryProviderOption func(rp *RepositoryProvider)

func NewRepositoryProvider(options ...RepositoryProviderOption) *RepositoryProvider {
	rp := &RepositoryProvider{}
	for _, option := range options {
		option(rp)
	}
	return rp
}

func WithMemoryNoteRepository() RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.NoteRepository = storagememory.NewNoteRepository()
	}
}

func WithMongoDBNoteRepository(database *mongo.Database) RepositoryProviderOption {
	return func(rp *RepositoryProvider) {
		rp.NoteRepository = storagemongo.NewNoteRepository(database)
	}
}
