package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v9"
	"github.com/nazarslota/unotes/auth/internal/domain/refreshtoken"
)

type refreshTokenRepository struct {
	client *redis.Client
}

const RefreshTokensPrefix = "REFRESH_TOKENS_"

var _ refreshtoken.Repository = (*refreshTokenRepository)(nil)

func NewRefreshTokenRepository(client *redis.Client) refreshtoken.Repository {
	return &refreshTokenRepository{client: client}
}

func (r *refreshTokenRepository) SaveOne(ctx context.Context, userID string, token *refreshtoken.Token) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to save the refresh token: %w", ctx.Err())
	default:
	}

	tokens, err := r.FindMany(ctx, userID)
	if errors.Is(err, refreshtoken.ErrTokensNotFound) {
		tokens = make([]refreshtoken.Token, 0)
	} else if err != nil {
		return fmt.Errorf("failed to receive all tokens: %w", err)
	}
	tokens = append(tokens, *token)

	b, err := json.Marshal(tokens)
	if err != nil {
		return fmt.Errorf("failed to marshal the received tokens: %w", err)
	}

	if err := r.client.Set(ctx, RefreshTokensPrefix+userID, string(b), 0).Err(); err != nil {
		return fmt.Errorf("failed to save the token: %w", err)
	}
	return nil
}

func (r *refreshTokenRepository) FindOne(
	ctx context.Context, userID string, token *refreshtoken.Token,
) (*refreshtoken.Token, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("failed to find the refresh token: %w", ctx.Err())
	default:
	}

	tokens, err := r.FindMany(ctx, userID)
	if errors.Is(err, refreshtoken.ErrTokensNotFound) {
		return nil, refreshtoken.ErrTokenNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to receive all tokens: %w", err)
	}

	for _, t := range tokens {
		if t == *token {
			return &t, nil
		}
	}
	return nil, refreshtoken.ErrTokenNotFound
}

func (r *refreshTokenRepository) DeleteOne(
	ctx context.Context, userID string, token *refreshtoken.Token,
) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to delete the refresh token: %w", ctx.Err())
	default:
	}

	tokens, err := r.FindMany(ctx, userID)
	if errors.Is(err, refreshtoken.ErrTokensNotFound) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to receive all tokens: %w", err)
	}

	for i, t := range tokens {
		if t == *token {
			tokens = append(tokens[:i], tokens[i+1:]...)
			break
		}
	}

	b, err := json.Marshal(tokens)
	if err != nil {
		return fmt.Errorf("failed to marshal the received tokens: %w", err)
	}

	if err := r.client.Set(ctx, RefreshTokensPrefix+userID, string(b), 0).Err(); err != nil {
		return fmt.Errorf("failed to save the token: %w", err)
	}

	return nil
}

func (r *refreshTokenRepository) FindMany(ctx context.Context, userID string) ([]refreshtoken.Token, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("failed to find the refresh tokens: %w", ctx.Err())
	default:
	}

	result, err := r.client.Get(ctx, RefreshTokensPrefix+userID).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to take the tokens: %w", err)
	} else if err == redis.Nil {
		result = "[]"
	}

	tokens := make([]refreshtoken.Token, 0)
	if err := json.Unmarshal([]byte(result), &tokens); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the received tokens: %w", err)
	}

	if len(tokens) == 0 {
		return nil, refreshtoken.ErrTokensNotFound
	}
	return tokens, nil
}
