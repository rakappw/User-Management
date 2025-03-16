package usecase

import (
	"errors"
	"time"

	"User_Management/internal/domain"
	"User_Management/internal/presenter"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Login(input presenter.LoginUserInput) (*presenter.LoginResponse, error)
	Logout(userID int) (*presenter.LogoutResponse, error)
}

type authUseCase struct {
	userRepo    domain.UserRepository
	authRepo    domain.AuthRepository
	jwtSecret   string
	tokenExpiry time.Duration
}

func NewAuthUseCase(userRepo domain.UserRepository, authRepo domain.AuthRepository, jwtSecret string, tokenExpiry time.Duration) AuthUseCase {
	return &authUseCase{
		userRepo:    userRepo,
		authRepo:    authRepo,
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
	}
}

func (uc *authUseCase) Login(input presenter.LoginUserInput) (*presenter.LoginResponse, error) {
	user, err := uc.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	auth, err := uc.authRepo.FindByUserID(user.ID)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(input.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(uc.tokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return nil, err
	}

	err = uc.authRepo.UpdateToken(user.ID, signedToken, time.Now())
	if err != nil {
		return nil, err
	}

	return &presenter.LoginResponse{
		Token: signedToken,
	}, nil
}

func (uc *authUseCase) Logout(userID int) (*presenter.LogoutResponse, error) {
	err := uc.authRepo.ClearToken(userID, time.Now())
	if err != nil {
		return nil, err
	}

	return &presenter.LogoutResponse{
		Message: "Logout successful",
	}, nil
}
