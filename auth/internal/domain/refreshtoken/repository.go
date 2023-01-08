package refreshtoken

import (
	"context"
	"errors"
)

// Repository is the interface for a token repository.
type Repository interface {
	// SaveOne saves a single token to the repository.
	SaveOne(ctx context.Context, userID string, token *Token) error
	// FindOne retrieves a single token from the repository.
	FindOne(ctx context.Context, userID string, token *Token) (*Token, error)
	// DeleteOne deletes a single token from the repository.
	DeleteOne(ctx context.Context, userID string, token *Token) error
	// FindMany retrieves multiple tokens from the repository.
	FindMany(ctx context.Context, userID string) ([]Token, error)
}

// ErrTokenNotFound is returned when a token is not found.
var ErrTokenNotFound = errors.New("token not found")

// ErrTokensNotFound is returned when tokens are not found.
var ErrTokensNotFound = errors.New("tokens not found")
