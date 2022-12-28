package refreshtoken

import "context"

type Repository interface {
	SaveOne(ctx context.Context, userID string, token *Token) error
	FindOne(ctx context.Context, userID string, token *Token) (*Token, error)
	DeleteOne(ctx context.Context, userID string, token *Token) error
	FindMany(ctx context.Context, userID string) ([]Token, error)
}
