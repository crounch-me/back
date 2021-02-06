package adapters

import (
	"errors"
	"sync"

	"github.com/crounch-me/back/internal/account/domain/authorizations"
)

type AuthorizationsMemoryRepository struct {
	authorizations map[string]string
	lock           *sync.RWMutex
}

func NewAuthorizationsMemoryRepository() authorizations.Repository {
	return &AuthorizationsMemoryRepository{
		authorizations: make(map[string]string, 0),
		lock:           &sync.RWMutex{},
	}
}

func (r *AuthorizationsMemoryRepository) AddAuthorization(userUUID, token string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.authorizations[token] = userUUID

	return nil
}

func (r *AuthorizationsMemoryRepository) GetUserUUIDByToken(token string) (string, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	userUUID, ok := r.authorizations[token]
	if !ok {
		return "", errors.New("authorization not found")
	}

	return userUUID, nil
}

func (r *AuthorizationsMemoryRepository) RemoveAuthorization(userUUID, token string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.authorizations[token]; !ok {
		return errors.New("authorization not found")
	}

	delete(r.authorizations, token)

	return nil
}
