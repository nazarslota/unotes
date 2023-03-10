package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	domain "github.com/nazarslota/unotes/auth/internal/domain/user"
)

// UserRepository provides an implementation of the user repository for a PostgreSQL database.
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new instance of the UserRepository with the provided handle to the PostgreSQL database.
//
// If db is nil, returns an error.
func NewUserRepository(db *sqlx.DB) (*UserRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}
	return &UserRepository{db: db}, nil
}

// SaveUser saves a user to the PostgreSQL database.
//
// If the user already exists, returns `user.ErrUserAlreadyExists`.
func (r UserRepository) SaveUser(ctx context.Context, user domain.User) error {
	query := fmt.Sprintf(`INSERT INTO users (id, username, password_hash) VALUES ($1, $2, $3) ON CONFLICT (username) DO NOTHING`)

	res, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	} else if affected == 0 {
		return domain.ErrUserAlreadyExists
	}
	return nil
}

// FindUserByUserID finds a user in the PostgreSQL database by their user ID.
//
// If the user is not found, returns `user.ErrUserNotFound`.
func (r UserRepository) FindUserByUserID(ctx context.Context, userID string) (user domain.User, err error) {
	query := fmt.Sprintf(`SELECT * FROM users WHERE id = $1`)
	if err := r.db.GetContext(ctx, &user, query, userID); errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("failed to execute query: %w", err)
		return domain.User{}, errors.Join(err, domain.ErrUserNotFound)
	} else if err != nil {
		return domain.User{}, fmt.Errorf("failed to execute query: %w", err)
	}
	return user, nil
}

// FindUserByUsername finds a user in the PostgreSQL database by their username.
//
// If the user is not found, returns `user.ErrUserNotFound`.
func (r UserRepository) FindUserByUsername(ctx context.Context, username string) (user domain.User, err error) {
	query := fmt.Sprintf(`SELECT * FROM users WHERE username = $1`)
	if err := r.db.GetContext(ctx, &user, query, username); errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("failed to execute query: %w", err)
		return domain.User{}, errors.Join(err, domain.ErrUserNotFound)
	} else if err != nil {
		return domain.User{}, fmt.Errorf("failed to execute query: %w", err)
	}
	return user, nil
}
