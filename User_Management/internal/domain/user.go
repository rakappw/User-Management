package domain

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRepository interface {
	Create(user *User) error
	FindByID(id int) (*User, error)
	FindByEmail(email string) (*User, error)
}
