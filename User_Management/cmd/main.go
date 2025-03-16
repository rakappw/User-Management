package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"User_Management/delivery/http"
	"User_Management/infrastructure"
	"User_Management/internal/usecase"
)

const (
	jwtSecret   = "your-super-secret-key"
	tokenExpiry = 24 * time.Hour
	serverAddr  = ":8080"
)

func main() {
	userRepo := infrastructure.NewInMemoryUserRepository()
	authRepo := infrastructure.NewInMemoryAuthRepository()

	userUseCase := usecase.NewUserUseCase(userRepo, authRepo)
	authUseCase := usecase.NewAuthUseCase(userRepo, authRepo, jwtSecret, tokenExpiry)

	router := gin.Default()

	http.NewUserController(router, userUseCase, authUseCase, jwtSecret)

	log.Println("Server starting on", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
