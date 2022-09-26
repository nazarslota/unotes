package user

import "context"

type Repository interface {
	Create(ctx context.Context, u User) error
	FindByUsername(ctx context.Context, username string) (*User, error)
}
