package redis

import (
	"context"
	"testing"

	"github.com/nazarslota/unotes/auth/internal/domain/refresh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	userAID                   = "user-a-id"
	userATokenA refresh.Token = "user-a-token-a"
	userATokenB refresh.Token = "user-a-token-b"

	userBID                   = "user-b-id"
	userBTokenA refresh.Token = "user-b-token-a"
	userBTokenB refresh.Token = "user-b-token-b"
)

var repository *RefreshTokenRepository

func init() {
	cfg := Config{}
	db, err := NewRedis(context.Background(), cfg)
	if err != nil {
		panic(err)
	}

	repository, err = NewRefreshTokenRepository(db)
	if err != nil {
		panic(err)
	}
}

func TestNewRefreshTokenRepository(t *testing.T) {
	redis, err := NewRedis(context.Background(), Config{})
	require.NoError(t, err)
	require.NotNil(t, redis)

	t.Cleanup(func() {
		_ = redis.Close()
	})

	t.Run("should create a new refresh token repository", func(t *testing.T) {
		repository, err := NewRefreshTokenRepository(redis)
		require.NoError(t, err)
		require.NotNil(t, repository)
		require.NotNil(t, repository.db)
	})

	t.Run("should return an error when redis is nil", func(t *testing.T) {
		repository, err := NewRefreshTokenRepository(nil)
		require.EqualError(t, err, "redis db is nil")
		require.Nil(t, repository)
	})
}

func TestRefreshTokenRepository_SaveRefreshToken(t *testing.T) {
	t.Run("should save a refresh token", func(t *testing.T) {
		err := repository.SaveRefreshToken(context.Background(), userAID, userATokenA)
		assert.NoError(t, err)

		err = repository.SaveRefreshToken(context.Background(), userAID, userATokenB)
		assert.NoError(t, err)

		result, err := repository.db.SMembers(context.Background(), refreshTokenKeyFromUserID("user-a-id")).Result()
		require.NoError(t, err)
		assert.Contains(t, result, string(userATokenA))
		assert.Contains(t, result, string(userATokenB))

		t.Cleanup(func() {
			_ = repository.db.FlushDB(context.Background())
		})
	})

	t.Run("should return an error if context is invalid", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := repository.SaveRefreshToken(ctx, userAID, userATokenA)
		assert.ErrorIs(t, err, context.Canceled)

		err = repository.SaveRefreshToken(ctx, userAID, userATokenB)
		assert.ErrorIs(t, err, context.Canceled)

		result, err := repository.db.SMembers(context.Background(), userAID).Result()
		require.NoError(t, err)
		assert.NotContains(t, result, string(userATokenA))
		assert.NotContains(t, result, string(userATokenB))

		t.Cleanup(func() {
			_ = repository.db.FlushDB(context.Background())
		})
	})
}

func TestRefreshTokenRepository_DeleteRefreshToken(t *testing.T) {
	t.Run("should delete a refresh token", func(t *testing.T) {
		err := repository.db.SAdd(context.Background(), refreshTokenKeyFromUserID(userAID), userATokenA, userATokenB).Err()
		require.NoError(t, err)

		err = repository.DeleteRefreshToken(context.Background(), userAID, userATokenA)
		assert.NoError(t, err)

		members, err := repository.db.SMembers(context.Background(), refreshTokenKeyFromUserID(userAID)).Result()
		require.NoError(t, err)
		assert.NotContains(t, members, string(userATokenA))
		assert.Contains(t, members, string(userATokenB))

		err = repository.DeleteRefreshToken(context.Background(), userAID, userATokenB)
		assert.NoError(t, err)

		members, err = repository.db.SMembers(context.Background(), refreshTokenKeyFromUserID(userAID)).Result()
		require.NoError(t, err)
		assert.NotContains(t, members, string(userATokenA))
		assert.NotContains(t, members, string(userATokenB))

		t.Cleanup(func() {
			_ = repository.db.FlushDB(context.Background())
		})
	})

	t.Run("should return an error if context is invalid", func(t *testing.T) {
		err := repository.db.SAdd(context.Background(), refreshTokenKeyFromUserID(userAID), userATokenA, userATokenB).Err()
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err = repository.DeleteRefreshToken(ctx, userAID, userATokenA)
		assert.ErrorIs(t, err, context.Canceled)

		err = repository.DeleteRefreshToken(ctx, userAID, userATokenB)
		assert.ErrorIs(t, err, context.Canceled)

		members, err := repository.db.SMembers(context.Background(), refreshTokenKeyFromUserID(userAID)).Result()
		require.NoError(t, err)
		assert.Contains(t, members, string(userATokenA))
		assert.Contains(t, members, string(userATokenB))

		t.Cleanup(func() {
			_ = repository.db.FlushDB(context.Background())
		})
	})

	t.Run("should return an error if refresh token is not found", func(t *testing.T) {
		err := repository.DeleteRefreshToken(context.Background(), userAID, userATokenA)
		assert.ErrorIs(t, err, refresh.ErrTokenNotFound)

		err = repository.DeleteRefreshToken(context.Background(), userAID, userATokenA)
		assert.ErrorIs(t, err, refresh.ErrTokenNotFound)

		t.Cleanup(func() {
			_ = repository.db.FlushDB(context.Background())
		})
	})
}

func TestRefreshTokenRepository_GetRefreshToken(t *testing.T) {
	t.Run("should get a refresh token", func(t *testing.T) {
		err := repository.db.SAdd(context.Background(), refreshTokenKeyFromUserID(userAID), userATokenA, userATokenB).Err()
		require.NoError(t, err)

		token, err := repository.GetRefreshToken(context.Background(), userAID, userATokenA)
		assert.Equal(t, userATokenA, token)

		token, err = repository.GetRefreshToken(context.Background(), userAID, userATokenB)
		assert.Equal(t, userATokenB, token)

		t.Cleanup(func() {
			_ = repository.db.FlushDB(context.Background())
		})
	})

	t.Run("should return an error if context is invalid", func(t *testing.T) {
		err := repository.db.SAdd(context.Background(), refreshTokenKeyFromUserID(userAID), userATokenA, userATokenB).Err()
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		token, err := repository.GetRefreshToken(ctx, userAID, userATokenA)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Empty(t, token)

		token, err = repository.GetRefreshToken(ctx, userAID, userATokenB)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Empty(t, token)

		t.Cleanup(func() {
			_ = repository.db.FlushDB(context.Background())
		})
	})

	t.Run("should return an error if refresh token is not found", func(t *testing.T) {
		token, err := repository.GetRefreshToken(context.Background(), userAID, userATokenA)
		assert.ErrorIs(t, err, refresh.ErrTokenNotFound)
		assert.Empty(t, token)

		token, err = repository.GetRefreshToken(context.Background(), userAID, userATokenA)
		assert.ErrorIs(t, err, refresh.ErrTokenNotFound)
		assert.Empty(t, token)

		t.Cleanup(func() {
			_ = repository.db.FlushDB(context.Background())
		})
	})
}

func TestRefreshTokenRepository_SaveRefreshTokens(t *testing.T) {
	t.Run("should save refresh tokens", func(t *testing.T) {
		err := repository.SaveRefreshTokens(context.Background(), userAID, []refresh.Token{userATokenA, userATokenB})
		assert.NoError(t, err)

		err = repository.SaveRefreshTokens(context.Background(), userBID, []refresh.Token{userBTokenA, userBTokenB})
		assert.NoError(t, err)

		result, err := repository.db.SMembers(context.Background(), refreshTokenKeyFromUserID(userAID)).Result()
		require.NoError(t, err)

		assert.Contains(t, result, string(userATokenA))
		assert.Contains(t, result, string(userATokenB))
		assert.NotContains(t, result, string(userBTokenA))
		assert.NotContains(t, result, string(userBTokenB))

		result, err = repository.db.SMembers(context.Background(), refreshTokenKeyFromUserID(userBID)).Result()
		require.NoError(t, err)

		assert.NotContains(t, result, string(userATokenA))
		assert.NotContains(t, result, string(userATokenB))
		assert.Contains(t, result, string(userBTokenA))
		assert.Contains(t, result, string(userBTokenB))

		t.Cleanup(func() {
			_ = repository.db.FlushDB(context.Background())
		})
	})

	t.Run("should return an error if context is invalid", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := repository.SaveRefreshTokens(ctx, userAID, []refresh.Token{userATokenA, userATokenB})
		assert.ErrorIs(t, err, context.Canceled)

		result, err := repository.db.SMembers(context.Background(), refreshTokenKeyFromUserID(userAID)).Result()
		require.NoError(t, err)
		require.Empty(t, result)

		t.Cleanup(func() {
			_ = repository.db.FlushDB(context.Background())
		})
	})
}

func TestRefreshTokenRepository_DeleteRefreshTokens(t *testing.T) {
	t.Run("should delete refresh tokens", func(t *testing.T) {
		err := repository.db.SAdd(context.Background(), refreshTokenKeyFromUserID(userAID), userATokenA, userATokenB).Err()
		require.NoError(t, err)

		err = repository.db.SAdd(context.Background(), refreshTokenKeyFromUserID(userBID), userBTokenA, userBTokenB).Err()
		require.NoError(t, err)

		err = repository.DeleteRefreshTokens(context.Background(), userAID, []refresh.Token{userATokenA})
		assert.NoError(t, err)

		result, err := repository.db.SMembers(context.Background(), refreshTokenKeyFromUserID(userAID)).Result()
		require.NoError(t, err)
		assert.NotContains(t, result, string(userATokenA))
		assert.Contains(t, result, string(userATokenB))

		err = repository.DeleteRefreshTokens(context.Background(), userBID, []refresh.Token{userBTokenA, userBTokenB})
		assert.NoError(t, err)

		result, err = repository.db.SMembers(context.Background(), refreshTokenKeyFromUserID(userBID)).Result()
		require.NoError(t, err)
		assert.Empty(t, result)

		t.Cleanup(func() {
			_ = repository.db.FlushDB(context.Background())
		})
	})

	t.Run("should return an error if context is invalid", func(t *testing.T) {
		err := repository.db.SAdd(context.Background(), refreshTokenKeyFromUserID(userAID), userATokenA, userATokenB).Err()
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err = repository.DeleteRefreshTokens(ctx, userAID, []refresh.Token{userATokenA, userATokenB})
		assert.ErrorIs(t, err, context.Canceled)

		result, err := repository.db.SMembers(context.Background(), refreshTokenKeyFromUserID(userAID)).Result()
		require.NoError(t, err)
		assert.Contains(t, result, string(userATokenA))
		assert.Contains(t, result, string(userATokenB))

		t.Cleanup(func() {
			_ = repository.db.FlushDB(context.Background())
		})
	})

	t.Run("should return an error if refresh tokens is not found", func(t *testing.T) {
		err := repository.DeleteRefreshTokens(context.Background(), userAID, []refresh.Token{userATokenA, userATokenB})
		assert.ErrorIs(t, err, refresh.ErrTokenNotFound)

		t.Cleanup(func() {
			_ = repository.db.FlushDB(context.Background())
		})
	})
}

func TestRefreshTokenRepository_GetRefreshTokens(t *testing.T) {
	t.Run("should get refresh tokens", func(t *testing.T) {
		err := repository.db.SAdd(context.Background(), refreshTokenKeyFromUserID(userAID), userATokenA, userATokenB).Err()
		require.NoError(t, err)

		err = repository.db.SAdd(context.Background(), refreshTokenKeyFromUserID(userBID), userBTokenA, userBTokenB).Err()
		require.NoError(t, err)

		tokens, err := repository.GetRefreshTokens(context.Background(), userAID)
		assert.NoError(t, err)
		assert.Contains(t, tokens, userATokenA)
		assert.Contains(t, tokens, userATokenB)
		assert.NotContains(t, tokens, userBTokenA)
		assert.NotContains(t, tokens, userBTokenB)

		tokens, err = repository.GetRefreshTokens(context.Background(), userBID)
		assert.NoError(t, err)
		assert.NotContains(t, tokens, userATokenA)
		assert.NotContains(t, tokens, userATokenB)
		assert.Contains(t, tokens, userBTokenA)
		assert.Contains(t, tokens, userBTokenB)

		t.Cleanup(func() {
			_ = repository.db.FlushDB(context.Background())
		})
	})

	t.Run("should return an error if context is invalid", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		tokens, err := repository.GetRefreshTokens(ctx, userAID)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Empty(t, tokens)

		t.Cleanup(func() {
			_ = repository.db.FlushDB(context.Background())
		})
	})

	t.Run("should return an error if refresh tokens is not found", func(t *testing.T) {
		tokens, err := repository.GetRefreshTokens(context.Background(), userAID)
		assert.ErrorIs(t, err, refresh.ErrTokenNotFound)
		assert.Empty(t, tokens)

		t.Cleanup(func() {
			_ = repository.db.FlushDB(context.Background())
		})
	})
}
