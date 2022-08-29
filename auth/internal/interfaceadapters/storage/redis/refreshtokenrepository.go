package redis

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v9"
	"github.com/udholdenhed/unotes/auth/internal/domain/refreshtoken"
)

type RefreshTokenRepository struct {
	RedisClient *redis.Client
}

const UserRefreshTokensPrefix = "USER_REFRESH_TOKENS_"

var _ refreshtoken.Repository = (*RefreshTokenRepository)(nil)

func NewRefreshTokenRepository(client *redis.Client) refreshtoken.Repository {
	return &RefreshTokenRepository{RedisClient: client}
}

func (r *RefreshTokenRepository) SetRefreshToken(ctx context.Context, uid string, t refreshtoken.RefreshToken) error {
	tokens, err := r.GetAllRefreshTokens(ctx, uid)
	if err != nil {
		return err
	}
	tokens = append(tokens, t)

	b, err := json.Marshal(tokens)
	if err != nil {
		return err
	}
	return r.RedisClient.Set(ctx, UserRefreshTokensPrefix+uid, string(b), 0).Err()
}

func (r *RefreshTokenRepository) GetRefreshToken(ctx context.Context, uid string, t refreshtoken.RefreshToken) (*refreshtoken.RefreshToken, error) {
	tokens, err := r.GetAllRefreshTokens(ctx, uid)
	if err != nil {
		return nil, err
	}

	for _, v := range tokens {
		if v == t {
			return &v, nil
		}
	}
	return nil, nil
}

func (r *RefreshTokenRepository) DeleteRefreshToken(ctx context.Context, uid string, t refreshtoken.RefreshToken) error {
	tokens, err := r.GetAllRefreshTokens(ctx, uid)
	if err != nil {
		return err
	}

	for i, v := range tokens {
		if v == t {
			tokens = append(tokens[:i], tokens[i+1:]...)
			break
		}
	}

	b, err := json.Marshal(tokens)
	if err != nil {
		return err
	}
	return r.RedisClient.Set(ctx, UserRefreshTokensPrefix+uid, string(b), 0).Err()
}

func (r *RefreshTokenRepository) GetAllRefreshTokens(ctx context.Context, uid string) ([]refreshtoken.RefreshToken, error) {
	result, err := r.RedisClient.Get(ctx, UserRefreshTokensPrefix+uid).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	} else if err == redis.Nil {
		result = "[]"
	}

	tokens := make([]refreshtoken.RefreshToken, 0)
	if err := json.Unmarshal([]byte(result), &tokens); err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		return nil, nil
	}
	return tokens, nil
}
