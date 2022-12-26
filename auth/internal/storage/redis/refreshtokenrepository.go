package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v9"
	"github.com/udholdenhed/unotes/auth/internal/domain/refreshtoken"
)

type refreshTokenRepository struct {
	client *redis.Client
}

const RefreshTokensPrefix = "REFRESH_TOKENS_"

var _ refreshtoken.Repository = (*refreshTokenRepository)(nil)

func NewRefreshTokenRepository(client *redis.Client) refreshtoken.Repository {
	return &refreshTokenRepository{client: client}
}

func (r *refreshTokenRepository) SetRefreshToken(ctx context.Context, userID string, token refreshtoken.Token) error {
	tokens, err := r.GetAllRefreshTokens(ctx, userID)
	if err != nil {
		return fmt.Errorf("redis: %w", err)
	}
	tokens = append(tokens, token)

	b, err := json.Marshal(tokens)
	if err != nil {
		return fmt.Errorf("redis: %w", err)
	}

	if err := r.client.Set(ctx, RefreshTokensPrefix+userID, string(b), 0).Err(); err != nil {
		return fmt.Errorf("redis: %w", err)
	}

	return nil
}

func (r *refreshTokenRepository) GetRefreshToken(ctx context.Context, userID string, token refreshtoken.Token) (*refreshtoken.Token, error) {
	tokens, err := r.GetAllRefreshTokens(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}

	for _, t := range tokens {
		if t == token {
			return &t, nil
		}
	}

	return nil, nil
}

func (r *refreshTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, token refreshtoken.Token) error {
	tokens, err := r.GetAllRefreshTokens(ctx, userID)
	if err != nil {
		return fmt.Errorf("redis: %w", err)
	}

	for i, t := range tokens {
		if t == token {
			tokens = append(tokens[:i], tokens[i+1:]...)
			break
		}
	}

	b, err := json.Marshal(tokens)
	if err != nil {
		return fmt.Errorf("redis: %w", err)
	}

	if err := r.client.Set(ctx, RefreshTokensPrefix+userID, string(b), 0).Err(); err != nil {
		return fmt.Errorf("redis: %w", err)
	}

	return nil
}

func (r *refreshTokenRepository) GetAllRefreshTokens(ctx context.Context, userID string) ([]refreshtoken.Token, error) {
	result, err := r.client.Get(ctx, RefreshTokensPrefix+userID).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("redis: %w", err)
	} else if err == redis.Nil {
		result = "[]"
	}

	tokens := make([]refreshtoken.Token, 0)
	if err := json.Unmarshal([]byte(result), &tokens); err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}

	if len(tokens) == 0 {
		return nil, nil
	}
	return tokens, nil
}
