package memory

import (
	"context"
	"sync"

	"github.com/udholdenhed/unotes/user-service/internal/domain/user"
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

func (r *userRepository) Create(ctx context.Context, u user.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		r.users[u.ID] = u
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		u, ok := r.users[id]
		if ok {
			return &u, nil
		}
	}
	return nil, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		for _, u := range r.users {
			if username == u.Username {
				return &u, nil
			}
		}
	}
	return nil, nil
}
