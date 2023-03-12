package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/nazarslota/unotes/auth/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	userA = user.User{
		ID:           "867620c3-d77a-4d96-821c-ca65cbca8318",
		Username:     "user-a-username",
		PasswordHash: "user-a-password-hash",
	}
)

var repository *UserRepository

func init() {
	db, err := NewPostgreSQL(context.Background(), Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
	})
	if err != nil {
		panic(err)
	}

	repository, err = NewUserRepository(db)
	if err != nil {
		panic(err)
	}
}

func TestNewUserRepository(t *testing.T) {
	t.Run("should create new user repository", func(t *testing.T) {
		db, err := NewPostgreSQL(context.Background(), Config{
			Host:     "localhost",
			Port:     "5432",
			Username: "postgres",
			Password: "postgres",
			DBName:   "postgres",
			SSLMode:  "disable",
		})
		require.NoError(t, err)
		require.NotNil(t, db)

		repository, err := NewUserRepository(db)
		assert.NoError(t, err)
		assert.NotNil(t, repository)
	})

	t.Run("should return error when db is nil", func(t *testing.T) {
		repository, err := NewUserRepository(nil)
		assert.EqualError(t, err, "db is nil")
		assert.Nil(t, repository)
	})
}

func TestUserRepository_SaveUser(t *testing.T) {
	t.Run("should successfully save user", func(t *testing.T) {
		err := repository.SaveUser(context.Background(), userA)
		assert.NoError(t, err)

		query := fmt.Sprintf(`SELECT * FROM users WHERE users.id = $1`)

		var result user.User
		err = repository.db.Get(&result, query, userA.ID)
		require.NotErrorIs(t, err, sql.ErrNoRows)

		assert.NoError(t, err)
		assert.Equal(t, userA, result)

		t.Cleanup(func() {
			query := fmt.Sprintf(`DELETE FROM users WHERE id = $1`)
			_, _ = repository.db.Exec(query, userA.ID)
		})
	})

	t.Run("should return error if context is invalid", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := repository.SaveUser(ctx, userA)
		assert.ErrorIs(t, err, context.Canceled)

		t.Cleanup(func() {
			query := fmt.Sprintf(`DELETE FROM users WHERE id = $1`)
			_, _ = repository.db.Exec(query, userA.ID)
		})
	})

	t.Run("should return error if user already exists", func(t *testing.T) {
		query := fmt.Sprintf(`INSERT INTO users (id, username, password_hash) VALUES ($1, $2, $3) ON CONFLICT (username) DO NOTHING`)
		_, err := repository.db.Exec(query, userA.ID, userA.Username, userA.PasswordHash)
		require.NoError(t, err)

		err = repository.SaveUser(context.Background(), userA)
		assert.ErrorIs(t, err, user.ErrUserAlreadyExists)

		t.Cleanup(func() {
			query := fmt.Sprintf(`DELETE FROM users WHERE id = $1`)
			_, _ = repository.db.Exec(query, userA.ID)
		})
	})
}

func TestUserRepository_FindUserByUserID(t *testing.T) {
	t.Run("should successfully find user by id", func(t *testing.T) {
		query := fmt.Sprintf(`INSERT INTO users (id, username, password_hash) VALUES ($1, $2, $3) ON CONFLICT (username) DO NOTHING`)
		_, err := repository.db.Exec(query, userA.ID, userA.Username, userA.PasswordHash)
		require.NoError(t, err)

		result, err := repository.FindUserByUserID(context.Background(), userA.ID)
		assert.NoError(t, err)
		assert.Equal(t, userA, result)

		t.Cleanup(func() {
			query := fmt.Sprintf(`DELETE FROM users WHERE id = $1`)
			_, _ = repository.db.Exec(query, userA.ID)
		})
	})

	t.Run("should return error if context is invalid", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		result, err := repository.FindUserByUserID(ctx, userA.ID)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Empty(t, result)

		t.Cleanup(func() {
			query := fmt.Sprintf(`DELETE FROM users WHERE id = $1`)
			_, _ = repository.db.Exec(query, userA.ID)
		})
	})

	t.Run("should return error if user does not exist", func(t *testing.T) {
		result, err := repository.FindUserByUserID(context.Background(), userA.ID)
		assert.ErrorIs(t, err, user.ErrUserNotFound)
		assert.Empty(t, result)
	})
}

func TestUserRepository_FindUserByUsername(t *testing.T) {
	t.Run("should successfully find user by username", func(t *testing.T) {
		query := fmt.Sprintf(`INSERT INTO users (id, username, password_hash) VALUES ($1, $2, $3) ON CONFLICT (username) DO NOTHING`)
		_, err := repository.db.Exec(query, userA.ID, userA.Username, userA.PasswordHash)
		require.NoError(t, err)

		result, err := repository.FindUserByUsername(context.Background(), userA.Username)
		assert.NoError(t, err)
		assert.Equal(t, userA, result)

		t.Cleanup(func() {
			query := fmt.Sprintf(`DELETE FROM users WHERE username = $1`)
			_, _ = repository.db.Exec(query, userA.Username)
		})
	})

	t.Run("should return error if context is invalid", func(t *testing.T) {
		query := fmt.Sprintf(`INSERT INTO users (id, username, password_hash) VALUES ($1, $2, $3) ON CONFLICT (username) DO NOTHING`)
		_, err := repository.db.Exec(query, userA.ID, userA.Username, userA.PasswordHash)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		result, err := repository.FindUserByUsername(ctx, userA.Username)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Empty(t, result)

		t.Cleanup(func() {
			query := fmt.Sprintf(`DELETE FROM users WHERE username = $1`)
			_, _ = repository.db.Exec(query, userA.Username)
		})
	})

	t.Run("should return error if user does not exist", func(t *testing.T) {
		result, err := repository.FindUserByUsername(context.Background(), userA.Username)
		assert.ErrorIs(t, err, user.ErrUserNotFound)
		assert.Empty(t, result)
	})
}
