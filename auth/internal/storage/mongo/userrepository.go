package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/udholdenhed/unotes/auth/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

var _ user.Repository = (*userRepository)(nil)

func NewUserRepository(db *mongo.Database, collection ...string) user.Repository {
	if len(collection) > 0 {
		return &userRepository{collection: db.Collection(collection[0])}
	}
	return &userRepository{collection: db.Collection("users")}
}

func (r *userRepository) SaveOne(ctx context.Context, user *user.User) error {

	if _, err := r.collection.InsertOne(ctx, user); err != nil {
		return fmt.Errorf("failed to save the user: %w", err)
	}
	return nil
}

func (r *userRepository) FindOne(ctx context.Context, username string) (*user.User, error) {
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
