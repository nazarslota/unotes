package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v9"
	"github.com/nazarslota/unotes/auth/internal/domain/refreshtoken"
)

// RefreshTokenRepository is a struct that implements the refreshtoken.Repository interface using a Redis client.
type RefreshTokenRepository struct {
	client *redis.Client
}

// refreshTokensPrefix is a string prefix used for the keys in Redis where the refresh tokens are stored.
const refreshTokensPrefix = "REFRESH_TOKENS_"

// NewRefreshTokenRepository returns a new RefreshTokenRepository with the given Redis client.
func NewRefreshTokenRepository(client *redis.Client) *RefreshTokenRepository {
	return &RefreshTokenRepository{client: client}
}

// SaveOne saves the given refresh token to Redis. If an error occurs, it is returned.
func (r *RefreshTokenRepository) SaveOne(ctx context.Context, userID string, token *refreshtoken.Token) error {
	// Get all tokens belonging to the user.
	tokens, err := r.FindMany(ctx, userID)
	if errors.Is(err, refreshtoken.ErrTokensNotFound) {
		tokens = make([]refreshtoken.Token, 0)
	} else if err != nil {
		return fmt.Errorf("failed to receive all tokens: %w", err)
	}

	// Add the new token to the list.
	tokens = append(tokens, *token)

	// Marshal the list of tokens.
	b, err := json.Marshal(tokens)
	if err != nil {
		return fmt.Errorf("failed to marshal the received tokens: %w", err)
	}

	// Save the list of tokens to Redis.
	if err := r.client.Set(ctx, refreshTokensPrefix+userID, string(b), 0).Err(); err != nil {
		return fmt.Errorf("failed to save the token: %w", err)
	}
	return nil
}

// FindOne finds the given refresh token in Redis. If the token is not found, it returns an error.
func (r *RefreshTokenRepository) FindOne(ctx context.Context, userID string, token *refreshtoken.Token) (*refreshtoken.Token, error) {
	// Get all tokens belonging to the user.
	tokens, err := r.FindMany(ctx, userID)
	if errors.Is(err, refreshtoken.ErrTokensNotFound) {
		return nil, refreshtoken.ErrTokenNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to receive all tokens: %w", err)
	}

	// Find the token in the list.
	for _, t := range tokens {
		if t == *token {
			return &t, nil
		}
	}
	return nil, refreshtoken.ErrTokenNotFound
}

// DeleteOne deletes the given refresh token from Redis. If an error occurs, it is returned.
func (r *RefreshTokenRepository) DeleteOne(
	ctx context.Context, userID string, token *refreshtoken.Token,
) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to delete the refresh token: %w", ctx.Err())
	default:
	}

	// Get all tokens belonging to the user.
	tokens, err := r.FindMany(ctx, userID)
	if errors.Is(err, refreshtoken.ErrTokensNotFound) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to receive all tokens: %w", err)
	}

	// Remove the token from the list.
	for i, t := range tokens {
		if t == *token {
			tokens = append(tokens[:i], tokens[i+1:]...)
			break
		}
	}

	// Marshal the updated list of tokens.
	b, err := json.Marshal(tokens)
	if err != nil {
		return fmt.Errorf("failed to marshal the received tokens: %w", err)
	}

	if err := r.client.Set(ctx, refreshTokensPrefix+userID, string(b), 0).Err(); err != nil {
		return fmt.Errorf("failed to save the token: %w", err)
	}

	return nil
}

// FindMany returns all refresh tokens belonging to the given user from Redis. If an error occurs, it is returned.
func (r *RefreshTokenRepository) FindMany(ctx context.Context, userID string) ([]refreshtoken.Token, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("failed to find the refresh tokens: %w", ctx.Err())
	default:
	}

	// Get the list of tokens from Redis.
	val, err := r.client.Get(ctx, refreshTokensPrefix+userID).Result()
	if err == redis.Nil {
		return nil, refreshtoken.ErrTokensNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to get the tokens: %w", err)
	}

	// Unmarshal the list of tokens.
	var tokens []refreshtoken.Token
	if err := json.Unmarshal([]byte(val), &tokens); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the received tokens: %w", err)
	}
	return tokens, nil
}
