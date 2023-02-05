package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	domainuser "github.com/nazarslota/unotes/auth/internal/domain/user"
)

type UserRepository struct {
	db *sqlx.DB
}

var _ domainuser.Repository = (*UserRepository)(nil)

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) SaveUser(ctx context.Context, user domainuser.User) error {
	query := fmt.Sprintf(`INSERT INTO users (id, username, password_hash) VALUES ($1, $2, $3) ON CONFLICT (username) DO NOTHING`)

	res, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	} else if affected == 0 {
		return domainuser.ErrUserAlreadyExists
	}
	return nil
}

func (r *UserRepository) FindUserByUserID(ctx context.Context, userID string) (user domainuser.User, err error) {
	query := fmt.Sprintf(`SELECT * FROM users WHERE id = $1`)
	if err := r.db.GetContext(ctx, &user, query, userID); errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("failed to execute query: %w", err)
		return domainuser.User{}, errors.Join(err, domainuser.ErrUserNotFound)
	} else if err != nil {
		return domainuser.User{}, fmt.Errorf("failed to execute query: %w", err)
	}
	return user, nil
}

func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (user domainuser.User, err error) {
	query := fmt.Sprintf(`SELECT * FROM users WHERE username = $1`)
	if err := r.db.GetContext(ctx, &user, query, username); errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("failed to execute query: %w", err)
		return domainuser.User{}, errors.Join(err, domainuser.ErrUserNotFound)
	} else if err != nil {
		return domainuser.User{}, fmt.Errorf("failed to execute query: %w", err)
	}
	return user, nil
}
