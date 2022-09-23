package memory

import (
	"context"
	"sync"

	"github.com/udholdenhed/unotes/auth-service/internal/domain/refreshtoken"
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

func (r *refreshTokenRepository) SetRefreshToken(ctx context.Context, userID string, token refreshtoken.Token) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if r.tokens[userID] == nil {
			r.tokens[userID] = make([]refreshtoken.Token, 0)
		}
		r.tokens[userID] = append(r.tokens[userID], token)
	}
	return nil
}

func (r *refreshTokenRepository) GetRefreshToken(ctx context.Context, userID string, token refreshtoken.Token) (*refreshtoken.Token, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		for _, t := range r.tokens[userID] {
			if t == token {
				return &t, nil
			}
		}
	}
	return nil, nil
}

func (r *refreshTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, token refreshtoken.Token) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		for i, t := range r.tokens[userID] {
			if t == token {
				r.tokens[userID] = append(r.tokens[userID][:i], r.tokens[userID][i+1:]...)
			}
		}
	}
	return nil
}

func (r *refreshTokenRepository) GetAllRefreshTokens(ctx context.Context, userID string) ([]refreshtoken.Token, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if len(r.tokens[userID]) == 0 {
			return nil, nil
		}
		return r.tokens[userID], nil
	}
}
