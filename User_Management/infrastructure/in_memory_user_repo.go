package infrastructure

import (
	"errors"
	"sync"

	"User_Management/internal/domain"
)

type inMemoryUserRepository struct {
	mutex   sync.RWMutex
	users   map[int]*domain.User
	counter int
}

func NewInMemoryUserRepository() domain.UserRepository {
	return &inMemoryUserRepository{
		users:   make(map[int]*domain.User),
		counter: 0,
	}
}

func (r *inMemoryUserRepository) Create(user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.counter++
	user.ID = r.counter
	r.users[user.ID] = user
	return nil
}

func (r *inMemoryUserRepository) FindByID(id int) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *inMemoryUserRepository) FindByEmail(email string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}
