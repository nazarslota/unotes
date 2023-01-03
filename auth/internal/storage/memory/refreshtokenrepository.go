package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/nazarslota/unotes/auth/internal/domain/refreshtoken"
)

type refreshTokenRepository struct {
	tokens map[string][]refreshtoken.Token
	mutex  sync.RWMutex
}

var _ refreshtoken.Repository = (*refreshTokenRepository)(nil)

func NewRefreshTokenRepository() refreshtoken.Repository {
	return &refreshTokenRepository{
		tokens: make(map[string][]refreshtoken.Token),
		mutex:  sync.RWMutex{},
	}
}

func (r *refreshTokenRepository) SaveOne(ctx context.Context, userID string, token *refreshtoken.Token) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to save the refresh token: %w", ctx.Err())
	default:
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.tokens[userID] == nil {
		r.tokens[userID] = make([]refreshtoken.Token, 0)
	}
	r.tokens[userID] = append(r.tokens[userID], *token)

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

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, t := range r.tokens[userID] {
		if t == *token {
			return &t, nil
		}
	}
	return nil, refreshtoken.ErrTokenNotFound
}

func (r *refreshTokenRepository) DeleteOne(ctx context.Context, userID string, token *refreshtoken.Token) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to delete the refresh token: %w", ctx.Err())
	default:
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, t := range r.tokens[userID] {
		if t == *token {
			r.tokens[userID] = append(r.tokens[userID][:i], r.tokens[userID][i+1:]...)
		}
	}
	return nil
}

func (r *refreshTokenRepository) FindMany(ctx context.Context, userID string) ([]refreshtoken.Token, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("failed to find the refresh tokens: %w", ctx.Err())
	default:
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if len(r.tokens[userID]) == 0 {
		return nil, refreshtoken.ErrTokensNotFound
	}
	return r.tokens[userID], nil
}
