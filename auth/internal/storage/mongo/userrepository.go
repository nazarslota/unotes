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

func NewUserRepository(collection *mongo.Collection) user.Repository {
	return &userRepository{collection: collection}
}

func (r *userRepository) Create(ctx context.Context, u user.User) error {
	if _, err := r.collection.InsertOne(ctx, u); err != nil {
		return fmt.Errorf("failed to create a user: %w", err)
	}
	return nil
}

func (r *userRepository) FindOne(ctx context.Context, username string) (*user.User, error) {
	result := r.collection.FindOne(ctx, bson.M{"username": username})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find the user: %w", result.Err())
	}

	u := new(user.User)
	if err := result.Decode(&u); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}
	return u, nil
}
