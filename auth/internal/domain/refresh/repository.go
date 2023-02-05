package refresh

import (
	"context"
	"errors"
)

type Repository interface {
	SaveRefreshToken(ctx context.Context, userID string, tokens Token) error
	DeleteRefreshToken(ctx context.Context, userID string, tokens Token) error
	SaveRefreshTokens(ctx context.Context, userID string, tokens []Token) error
	DeleteRefreshTokens(ctx context.Context, userID string, tokens []Token) error
	GetRefreshTokens(ctx context.Context, userID string) ([]Token, error)
}

var ErrTokenNotFound = errors.New("token not found")
