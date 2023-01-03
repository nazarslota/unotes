package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/nazarslota/unotes/auth/internal/domain/user"
)

type userRepository struct {
	db *sqlx.DB
}

var _ user.Repository = (*userRepository)(nil)

func NewUserRepository(db *sqlx.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) SaveOne(ctx context.Context, user *user.User) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to save the user: %w", ctx.Err())
	default:
	}

	query := fmt.Sprintf(`INSERT INTO users (id, username, password_hash) VALUES ($1, $2, $3)`)
	if _, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.PasswordHash); err != nil {
		return fmt.Errorf("failed to save the user to the database: %w", err)
	}
	return nil
}

func (r *userRepository) FindOne(ctx context.Context, username string) (*user.User, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("failed to find the user: %w", ctx.Err())
	default:
	}

	query := fmt.Sprintf(`SELECT * FROM users WHERE username = $1`)

	u := &user.User{}
	if err := r.db.GetContext(ctx, u, query, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to fetch the user from the database: %w", err)
	}
	return u, nil
}
