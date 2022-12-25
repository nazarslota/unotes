package refreshtoken

import "context"

type Repository interface {
	SetRefreshToken(ctx context.Context, userID string, token Token) error
	GetRefreshToken(ctx context.Context, userID string, token Token) (*Token, error)
	DeleteRefreshToken(ctx context.Context, userID string, token Token) error
	GetAllRefreshTokens(ctx context.Context, userID string) ([]Token, error)
}
