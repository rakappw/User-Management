package usecase

import (
	"errors"
	"time"

	"User_Management/internal/domain"
	"User_Management/internal/presenter"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	Register(input presenter.RegisterUserInput) (*presenter.UserResponse, error)
	GetUserByID(userID int) (*presenter.UserResponse, error)
}

type userUseCase struct {
	userRepo domain.UserRepository
	authRepo domain.AuthRepository
}

func NewUserUseCase(userRepo domain.UserRepository, authRepo domain.AuthRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

func (uc *userUseCase) Register(input presenter.RegisterUserInput) (*presenter.UserResponse, error) {
	existingUser, _ := uc.userRepo.FindByEmail(input.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	user := &domain.User{
		Name:      input.Name,
		Email:     input.Email,
		CreatedAt: time.Now(),
	}

	err := uc.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	auth := &domain.Authentication{
		UserID:   user.ID,
		Password: string(hashedPassword),
	}

	err = uc.authRepo.Create(auth)
	if err != nil {
		return nil, err
	}

	return &presenter.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (uc *userUseCase) GetUserByID(userID int) (*presenter.UserResponse, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return &presenter.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}
