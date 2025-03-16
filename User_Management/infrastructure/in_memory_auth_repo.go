package infrastructure

import (
	"errors"
	"sync"
	"time"

	"User_Management/internal/domain"
)

type inMemoryAuthRepository struct {
	mutex   sync.RWMutex
	auths   map[int]*domain.Authentication
	counter int
}

func NewInMemoryAuthRepository() domain.AuthRepository {
	return &inMemoryAuthRepository{
		auths:   make(map[int]*domain.Authentication),
		counter: 0,
	}
}

func (r *inMemoryAuthRepository) Create(auth *domain.Authentication) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.counter++
	auth.ID = r.counter
	r.auths[auth.UserID] = auth
	return nil
}

func (r *inMemoryAuthRepository) FindByUserID(userID int) (*domain.Authentication, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	auth, exists := r.auths[userID]
	if !exists {
		return nil, errors.New("authentication not found")
	}
	return auth, nil
}

func (r *inMemoryAuthRepository) UpdateToken(userID int, token string, loginAt time.Time) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	auth, exists := r.auths[userID]
	if !exists {
		return errors.New("authentication not found")
	}

	auth.Token = token
	auth.LoginAt = loginAt
	return nil
}

func (r *inMemoryAuthRepository) ClearToken(userID int, logoutAt time.Time) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	auth, exists := r.auths[userID]
	if !exists {
		return errors.New("authentication not found")
	}

	auth.Token = ""
	auth.LogoutAt = logoutAt
	return nil
}
