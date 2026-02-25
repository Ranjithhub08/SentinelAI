package auth

import (
	"context"
	"errors"
	"sync"
)

// Repository defines data access for users
type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

var ErrUserNotFound = errors.New("user not found")
var ErrUserExists = errors.New("user already exists")

type inMemoryRepository struct {
	mu    sync.RWMutex
	users map[string]*User
}

// NewRepository creates a new in-memory auth repository
func NewRepository() Repository {
	return &inMemoryRepository{
		users: make(map[string]*User),
	}
}

func (r *inMemoryRepository) CreateUser(ctx context.Context, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, u := range r.users {
		if u.Email == user.Email {
			return ErrUserExists
		}
	}
	r.users[user.ID] = user
	return nil
}

func (r *inMemoryRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, ErrUserNotFound
}
