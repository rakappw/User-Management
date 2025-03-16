package domain

import "time"

type Authentication struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	Password string    `json:"-"`
	Token    string    `json:"token,omitempty"`
	LoginAt  time.Time `json:"login_at,omitempty"`
	LogoutAt time.Time `json:"logout_at,omitempty"`
}

type AuthRepository interface {
	Create(auth *Authentication) error
	FindByUserID(userID int) (*Authentication, error)
	UpdateToken(userID int, token string, loginAt time.Time) error
	ClearToken(userID int, logoutAt time.Time) error
}
