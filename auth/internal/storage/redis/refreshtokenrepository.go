package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
	"github.com/nazarslota/unotes/auth/internal/domain/refresh"
)

// RefreshTokenRepository is a Redis repository for managing refresh tokens.
type RefreshTokenRepository struct {
	client *redis.Client
}

// NewRefreshTokenRepository creates a new RefreshTokenRepository with the provided Redis client, and returns an error
// if the client is nil.
func NewRefreshTokenRepository(client *redis.Client) (*RefreshTokenRepository, error) {
	if client == nil {
		return nil, fmt.Errorf("redis client is nil")
	}
	return &RefreshTokenRepository{client: client}, nil
}

// SaveRefreshToken saves the given refresh token associated with the specified user ID.
//
// Args:
//   - ctx: a `context.Context` used to execute Redis commands.
//   - userID: a string representing the ID of the user to save tokens for.
//   - token: a `refresh.Token` values to save.
//
// Returns:
//   - error: an error value indicating the success or failure of the operation. If an error
//     occurs while executing the Redis command, it is returned.
func (r RefreshTokenRepository) SaveRefreshToken(ctx context.Context, userID string, token refresh.Token) error {
	return r.SaveRefreshTokens(ctx, userID, []refresh.Token{token})
}

// DeleteRefreshToken removes a single refresh token for the given user ID.
// Args:
// - ctx: a `context.Context` used to execute Redis commands.
// - userID: a string representing the ID of the user to delete the token for.
// - token: a `refresh.Token` value to delete for.
//
// Returns:
//   - error: an error value indicating the success or failure of the operation. If the
//     specified token cannot be found, an error of type `refresh.ErrTokenNotFound` is returned.
func (r RefreshTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, token refresh.Token) error {
	return r.DeleteRefreshTokens(ctx, userID, []refresh.Token{token})
}

// GetRefreshToken returns the refresh token for the specified user ID and token value.
//
// Args:
//   - ctx: a `context.Context` used to execute Redis commands.
//   - userID: a string representing the ID of the user to get the token for.
//   - token: a `refresh.Token` value to search for.
//
// Returns:
//   - refresh.Token: the refresh token for the given user ID and token value, if it exists.
//   - error: an error value indicating the success or failure of the operation. If the token is
//     not found, an error value of `refresh.ErrTokenNotFound` is returned. If an error occurs while
//     executing the Redis command, it is also returned.
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

// SaveRefreshTokens saves the given refresh tokens associated with the specified user ID.
//
// Args:
//   - ctx: a `context.Context` used to execute Redis commands.
//   - userID: a string representing the ID of the user to save tokens for.
//   - tokens: a slice of `refresh.Token` values to save.
//
// Returns:
//   - error: an error value indicating the success or failure of the operation. If an error
//     occurs while executing the Redis command, it is returned.
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

// DeleteRefreshTokens removes the given refresh tokens associated with the specified user ID.
// If any of the tokens cannot be found, this method returns an error of type `refresh.ErrTokenNotFound`.
//
// Args:
//   - ctx: a `context.Context` used to execute Redis commands.
//   - userID: a string representing the ID of the user to delete tokens for.
//   - tokens: a slice of `refresh.Token` values to delete.
//
// Returns:
//   - error: an error value indicating the success or failure of the operation. If any of the
//     specified tokens cannot be found, an error of type `refresh.ErrTokenNotFound` is returned.
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

// GetRefreshTokens returns a slice of all the refresh tokens associated with the given user ID.
// If no tokens are found, this method returns an error of type `refresh.ErrTokenNotFound`.
//
// Args:
//   - ctx: a `context.Context` used to execute Redis commands.
//   - userID: a string representing the ID of the user to fetch tokens for.
//
// Returns:
//   - []refresh.Token: a slice of refresh tokens associated with the given user ID.
//   - error: an error value indicating the success or failure of the operation. If no tokens are
//     found, an error of type `refresh.ErrTokenNotFound` is returned.
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
