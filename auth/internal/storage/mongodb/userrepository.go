package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/nazarslota/unotes/auth/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database, collection ...string) *UserRepository {
	if len(collection) > 0 {
		return &UserRepository{collection: db.Collection(collection[0])}
	}
	return &UserRepository{collection: db.Collection("users")}
}

func (r *UserRepository) SaveOne(ctx context.Context, user *user.User) error {
	if _, err := r.collection.InsertOne(ctx, user); err != nil {
		return fmt.Errorf("failed to save the user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindOne(ctx context.Context, username string) (*user.User, error) {
	result := r.collection.FindOne(ctx, bson.M{"username": username})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("user could not be found: %w", user.ErrUserNotFound)
		}
		return nil, fmt.Errorf("user could not be found: %w", result.Err())
	}

	u := new(user.User)
	if err := result.Decode(&u); err != nil {
		return nil, fmt.Errorf("unable to decode the user: %w", err)
	}
	return u, nil
}
