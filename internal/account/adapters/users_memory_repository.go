package adapters

import (
	"sync"

	"github.com/crounch-me/back/internal/account/app"
	"github.com/crounch-me/back/internal/account/domain/users"
)

type UsersMemoryRepository struct {
	users map[string]*users.User
	lock  *sync.RWMutex
}

func NewUsersMemoryRepository() users.Repository {
	return &UsersMemoryRepository{
		users: make(map[string]*users.User, 0),
		lock:  &sync.RWMutex{},
	}
}

func (r *UsersMemoryRepository) AddUser(user *users.User) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.users[user.UUID()] = user
	return nil
}

func (r *UsersMemoryRepository) FindByEmail(email string) (*users.User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, u := range r.users {
		if u.Email() == email {
			return u, nil
		}
	}

	return nil, app.ErrUserNotFound
}
