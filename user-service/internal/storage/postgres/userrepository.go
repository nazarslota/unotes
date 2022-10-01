package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/udholdenhed/unotes/user-service/internal/domain/user"
)

type userRepository struct {
	db *sqlx.DB
}

var _ user.Repository = (*userRepository)(nil)

func NewUserRepository(db *sqlx.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u user.User) error {
	query := fmt.Sprintf(`INSERT INTO users (id, username, email, password_hash) VALUES ($1, $2, $3, $4)`)
	if _, err := r.db.ExecContext(ctx, query, u.ID, u.Username, u.Email, u.PasswordHash); err != nil {
		return fmt.Errorf("postgres.userRepository.Create: %w", err)
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	query := fmt.Sprintf(`SELECT * FROM users WHERE id = $1`)

	u := &user.User{}
	if err := r.db.GetContext(ctx, u, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("postgres.userRepository.FindByID: %w", err)
	}
	return u, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	query := fmt.Sprintf(`SELECT * FROM users WHERE username = $1`)

	u := &user.User{}
	if err := r.db.GetContext(ctx, u, query, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("postgres.userRepository.FindByUsername: %w", err)
	}
	return u, nil
}
