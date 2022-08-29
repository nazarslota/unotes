package memory

import (
	"context"
	"errors"

	"github.com/udholdenhed/unotes/auth/internal/domain/refreshtoken"
)

type RefreshTokenRepository struct {
	RefreshTokens map[string][]refreshtoken.RefreshToken
}

var _ refreshtoken.Repository = (*RefreshTokenRepository)(nil)

func NewRefreshTokenRepository() refreshtoken.Repository {
	return &RefreshTokenRepository{RefreshTokens: make(map[string][]refreshtoken.RefreshToken)}
}

func (r *RefreshTokenRepository) SetRefreshToken(_ context.Context, uid string, t refreshtoken.RefreshToken) error {
	if r.RefreshTokens[uid] == nil {
		r.RefreshTokens[uid] = make([]refreshtoken.RefreshToken, 0)
	}
	r.RefreshTokens[uid] = append(r.RefreshTokens[uid], t)
	return nil
}

func (r *RefreshTokenRepository) GetRefreshToken(_ context.Context, uid string, t refreshtoken.RefreshToken) (*refreshtoken.RefreshToken, error) {
	if len(r.RefreshTokens[uid]) == 0 {
		return nil, nil
	}

	for _, v := range r.RefreshTokens[uid] {
		if v.Token == t.Token {
			return &v, nil
		}
	}
	return nil, nil
}

func (r *RefreshTokenRepository) DeleteRefreshToken(_ context.Context, uid string, t refreshtoken.RefreshToken) error {
	if len(r.RefreshTokens[uid]) == 0 {
		return errors.New("no such t found")
	}

	for i, v := range r.RefreshTokens[uid] {
		if v.Token == t.Token {
			r.RefreshTokens[uid] = append(r.RefreshTokens[uid][:i], r.RefreshTokens[uid][i+1:]...)
		}
	}
	return nil
}

func (r *RefreshTokenRepository) GetAllRefreshTokens(_ context.Context, uid string) ([]refreshtoken.RefreshToken, error) {
	if len(r.RefreshTokens) == 0 {
		return nil, nil
	}
	return r.RefreshTokens[uid], nil
}
