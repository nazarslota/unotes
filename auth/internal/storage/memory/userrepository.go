package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/udholdenhed/unotes/auth/internal/domain/user"
)

type userRepository struct {
	users map[string]user.User
	mutex sync.RWMutex
}

var _ user.Repository = (*userRepository)(nil)

func NewUserRepository() user.Repository {
	return &userRepository{
		users: make(map[string]user.User),
		mutex: sync.RWMutex{},
	}
}

func (r *userRepository) SaveOne(ctx context.Context, user *user.User) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to save the user: %w", ctx.Err())
	default:
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.users[user.ID] = *user
	return nil
}

func (r *userRepository) FindOne(ctx context.Context, username string) (*user.User, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("failed to find the user: %w", ctx.Err())
	default:
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, u := range r.users {
		if username == u.Username {
			return &u, nil
		}
	}
	return nil, user.ErrUserNotFound
}
