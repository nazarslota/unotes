package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
	domainrefresh "github.com/nazarslota/unotes/auth/internal/domain/refresh"
)

type RefreshTokenRepository struct {
	client *redis.Client
}

var _ domainrefresh.Repository = (*RefreshTokenRepository)(nil)

func NewRefreshTokenRepository(client *redis.Client) *RefreshTokenRepository {
	return &RefreshTokenRepository{client: client}
}

func (r *RefreshTokenRepository) SaveRefreshToken(ctx context.Context, userID string, token domainrefresh.Token) error {
	return r.SaveRefreshTokens(ctx, userID, []domainrefresh.Token{token})
}

func (r *RefreshTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, token domainrefresh.Token) error {
	return r.DeleteRefreshTokens(ctx, userID, []domainrefresh.Token{token})
}

func (r *RefreshTokenRepository) SaveRefreshTokens(ctx context.Context, userID string, tokens []domainrefresh.Token) error {
	key := fmt.Sprintf("refresh_token:%s", userID)

	pipe := r.client.TxPipeline()
	for _, token := range tokens {
		pipe.SAdd(ctx, key, token)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to execute the pipeline: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) DeleteRefreshTokens(ctx context.Context, userID string, tokens []domainrefresh.Token) error {
	key := fmt.Sprintf("refresh_token:%s", userID)
	for _, token := range tokens {
		if err := r.client.SRem(ctx, key, token).Err(); err != nil {
			return fmt.Errorf("failed to remove: %w", err)
		}
	}
	return nil
}

func (r *RefreshTokenRepository) GetRefreshTokens(ctx context.Context, userID string) ([]domainrefresh.Token, error) {
	key := fmt.Sprintf("refresh_token:%s", userID)
	members, err := r.client.SMembersMap(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}

	l := len(members)
	if l == 0 {
		return nil, domainrefresh.ErrTokenNotFound
	}

	tokens := make([]domainrefresh.Token, 0, l)
	for token := range members {
		tokens = append(tokens, domainrefresh.Token(token))
	}
	return tokens, nil
}
