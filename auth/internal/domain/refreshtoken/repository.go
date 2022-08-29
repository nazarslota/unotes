package refreshtoken

import "context"

type Repository interface {
	SetRefreshToken(ctx context.Context, uid string, t RefreshToken) error
	GetRefreshToken(ctx context.Context, uid string, t RefreshToken) (*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, uid string, t RefreshToken) error
	GetAllRefreshTokens(ctx context.Context, uid string) ([]RefreshToken, error)
}
