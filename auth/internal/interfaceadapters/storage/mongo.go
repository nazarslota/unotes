package storage

import (
	"github.com/udholdenhed/unotes/auth/internal/interfaceadapters/storage/mongo"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func WithMongoUserRepository(users *mongodriver.Collection) RepositoryProviderOption {
	return func(services *RepositoryProvider) {
		services.UserRepository = mongo.NewUserRepository(users)
	}
}
