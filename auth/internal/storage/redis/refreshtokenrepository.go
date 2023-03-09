package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
	"github.com/nazarslota/unotes/auth/internal/domain/refresh"
)

// RefreshTokenRepository is a struct that contains a Redis client for interacting with the Redis server.
type RefreshTokenRepository struct {
	client *redis.Client
}

// NewRefreshTokenRepository creates a new instance of RefreshTokenRepository.
func NewRefreshTokenRepository(client *redis.Client) (*RefreshTokenRepository, error) {
	if client == nil {
		return nil, fmt.Errorf("redis client is nil")
	}
	return &RefreshTokenRepository{client: client}, nil
}

// SaveRefreshToken saves a single refresh token for a given user ID.
func (r RefreshTokenRepository) SaveRefreshToken(ctx context.Context, userID string, token refresh.Token) error {
	return r.SaveRefreshTokens(ctx, userID, []refresh.Token{token})
}

// DeleteRefreshToken deletes a single refresh token for a given user ID.
func (r RefreshTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, token refresh.Token) error {
	return r.DeleteRefreshTokens(ctx, userID, []refresh.Token{token})
}

// GetRefreshToken retrieves a single refresh token for a given user ID.
func (r RefreshTokenRepository) GetRefreshToken(ctx context.Context, userID string, token refresh.Token) (refresh.Token, error) {
	key := refreshTokenKeyFromUserID(userID)
	exists, err := r.client.SIsMember(ctx, key, token).Result()
	if err != nil {
		return "", fmt.Errorf("failed to check for refresh token: %w", err)
	}

	if exists {
		return token, nil
	}
	return "", refresh.ErrTokenNotFound
}

// SaveRefreshTokens saves multiple refresh tokens for a given user ID.
func (r RefreshTokenRepository) SaveRefreshTokens(ctx context.Context, userID string, tokens []refresh.Token) error {
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

// DeleteRefreshTokens deletes multiple refresh tokens for a given user ID.
func (r RefreshTokenRepository) DeleteRefreshTokens(ctx context.Context, userID string, tokens []refresh.Token) error {
	members := make([]any, 0, len(tokens))
	for _, token := range tokens {
		members = append(members, token)
	}

	key := refreshTokenKeyFromUserID(userID)
	result, err := r.client.SRem(ctx, key, members...).Result()
	if err != nil {
		return fmt.Errorf("failed to execute srem command: %w", err)
	} else if result == 0 {
		return refresh.ErrTokenNotFound
	}
	return nil
}

// GetRefreshTokens retrieves all refresh tokens for a given user ID.
func (r RefreshTokenRepository) GetRefreshTokens(ctx context.Context, userID string) ([]refresh.Token, error) {
	key := refreshTokenKeyFromUserID(userID)
	members, err := r.client.SMembersMap(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to execute smembers command: %w", err)
	}

	l := len(members)
	if l == 0 {
		return nil, refresh.ErrTokenNotFound
	}

	tokens := make([]refresh.Token, 0, l)
	for token := range members {
		tokens = append(tokens, refresh.Token(token))
	}
	return tokens, nil
}

const refreshTokenPrefix = "refresh-token"

func refreshTokenKeyFromUserID(userID string) string {
	return fmt.Sprintf("%s:%s", refreshTokenPrefix, userID)
}
