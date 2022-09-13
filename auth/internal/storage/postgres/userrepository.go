package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/udholdenhed/unotes/auth/internal/domain/user"
)

type userRepository struct {
	db *sqlx.DB
}

var _ user.Repository = (*userRepository)(nil)

func NewUserRepository(db *sqlx.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u user.User) error {
	query := fmt.Sprintf(`INSERT INTO users (id, username, password_hash) VALUES ($1, $2, $3)`)
	if _, err := r.db.ExecContext(ctx, query, u.ID, u.Username, u.PasswordHash); err != nil {
		return fmt.Errorf("failed to create a user: %w", err)
	}
	return nil
}

func (r *userRepository) FindOne(ctx context.Context, username string) (*user.User, error) {
	query := fmt.Sprintf(`SELECT * FROM users WHERE username = $1`)

	u := &user.User{}
	if err := r.db.GetContext(ctx, u, query, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find the user: %w", err)
	}
	return u, nil
}
