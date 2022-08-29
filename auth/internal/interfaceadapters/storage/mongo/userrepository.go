package mongo

import (
	"context"
	"errors"

	"github.com/udholdenhed/unotes/auth/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	UsersCollection *mongo.Collection
}

var _ user.Repository = (*UserRepository)(nil)

func NewUserRepository(db *mongo.Collection) user.Repository {
	return &UserRepository{UsersCollection: db}
}

func (ur *UserRepository) Create(ctx context.Context, u user.User) error {
	if _, err := ur.UsersCollection.InsertOne(ctx, u); err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) FindOne(ctx context.Context, username string) (*user.User, error) {
	result := ur.UsersCollection.FindOne(ctx, bson.M{"username": username})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, result.Err()
	}

	u := new(user.User)
	if err := result.Decode(&u); err != nil {
		return nil, err
	}
	return u, nil
}
