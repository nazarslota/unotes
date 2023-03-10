package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
	domain "github.com/nazarslota/unotes/auth/internal/domain/refresh"
)

// RefreshTokenRepository is a Redis repository for managing refresh tokens.
type RefreshTokenRepository struct {
	client *redis.Client
}

// NewRefreshTokenRepository creates a new RefreshTokenRepository with the provided Redis client.
//
// Returns an error if the client is nil.
func NewRefreshTokenRepository(client *redis.Client) (*RefreshTokenRepository, error) {
	if client == nil {
		return nil, fmt.Errorf("redis client is nil")
	}
	return &RefreshTokenRepository{client: client}, nil
}

// SaveRefreshToken saves the given refresh token associated with the specified user ID.
func (r RefreshTokenRepository) SaveRefreshToken(ctx context.Context, userID string, token domain.Token) error {
	return r.SaveRefreshTokens(ctx, userID, []domain.Token{token})
}

// DeleteRefreshToken removes a single refresh token for the given user ID.
//
// If the specified token cannot be found, an error of type `refresh.ErrTokenNotFound` is returned.
func (r RefreshTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, token domain.Token) error {
	return r.DeleteRefreshTokens(ctx, userID, []domain.Token{token})
}

// GetRefreshToken returns the refresh token for the specified user ID and token value.
//
// If the token is not found, an error value of `refresh.ErrTokenNotFound` is returned.
func (r RefreshTokenRepository) GetRefreshToken(ctx context.Context, userID string, token domain.Token) (domain.Token, error) {
	key := refreshTokenKeyFromUserID(userID)
	exists, err := r.client.SIsMember(ctx, key, token).Result()
	if err != nil {
		return "", fmt.Errorf("failed to check for refresh token: %w", err)
	}

	if exists {
		return token, nil
	}
	return "", domain.ErrTokenNotFound
}

// SaveRefreshTokens saves the given refresh tokens associated with the specified user ID.
func (r RefreshTokenRepository) SaveRefreshTokens(ctx context.Context, userID string, tokens []domain.Token) error {
	members := make([]any, 0, len(tokens))
	for _, token := range tokens {
		members = append(members, token)
	}

	key := refreshTokenKeyFromUserID(userID)

	pipe := r.client.TxPipeline()
	pipe.SAdd(ctx, key, members...)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to execute pipeline: %w", err)
	}
	return nil
}

// DeleteRefreshTokens removes the given refresh tokens associated with the specified user ID.
//
// If any of the tokens cannot be found, this method returns an error of type `refresh.ErrTokenNotFound`.
func (r RefreshTokenRepository) DeleteRefreshTokens(ctx context.Context, userID string, tokens []domain.Token) error {
	members := make([]any, 0, len(tokens))
	for _, token := range tokens {
		members = append(members, token)
	}

	key := refreshTokenKeyFromUserID(userID)
	result, err := r.client.SRem(ctx, key, members...).Result()
	if err != nil {
		return fmt.Errorf("failed to execute srem command: %w", err)
	} else if result == 0 {
		return domain.ErrTokenNotFound
	}
	return nil
}

// GetRefreshTokens returns a slice of all the refresh tokens associated with the given user ID.
//
// If no tokens are found, this method returns an error of type `refresh.ErrTokenNotFound`.
func (r RefreshTokenRepository) GetRefreshTokens(ctx context.Context, userID string) ([]domain.Token, error) {
	key := refreshTokenKeyFromUserID(userID)
	members, err := r.client.SMembersMap(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to execute smembers command: %w", err)
	}

	l := len(members)
	if l == 0 {
		return nil, domain.ErrTokenNotFound
	}

	tokens := make([]domain.Token, 0, l)
	for token := range members {
		tokens = append(tokens, domain.Token(token))
	}
	return tokens, nil
}

const refreshTokenPrefix = "refresh-token"

func refreshTokenKeyFromUserID(userID string) string {
	return fmt.Sprintf("%s:%s", refreshTokenPrefix, userID)
}
