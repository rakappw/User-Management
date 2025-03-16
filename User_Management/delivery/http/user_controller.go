package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"User_Management/internal/presenter"
	"User_Management/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type UserController struct {
	userUseCase usecase.UserUseCase
	authUseCase usecase.AuthUseCase
	jwtSecret   string
}

func NewUserController(router *gin.Engine, userUseCase usecase.UserUseCase, authUseCase usecase.AuthUseCase, jwtSecret string) {
	controller := &UserController{
		userUseCase: userUseCase,
		authUseCase: authUseCase,
		jwtSecret:   jwtSecret,
	}

	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)
	router.GET("/profile", controller.AuthMiddleware(), controller.GetProfile)
	router.POST("/logout", controller.AuthMiddleware(), controller.Logout)
}

func (c *UserController) Register(ctx *gin.Context) {
	var input presenter.RegisterUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.userUseCase.Register(input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *UserController) Login(ctx *gin.Context) {
	var input presenter.LoginUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.authUseCase.Login(input)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	response, err := c.userUseCase.GetUserByID(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *UserController) Logout(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	response, err := c.authUseCase.Logout(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *UserController) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		if tokenString == authHeader {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(c.jwtSecret), nil
		})

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, ok := claims["user_id"]
			if !ok {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
				return
			}

			var userIDInt int
			switch v := userID.(type) {
			case float64:
				userIDInt = int(v)
			case string:
				id, err := strconv.Atoi(v)
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID format"})
					return
				}
				userIDInt = id
			default:
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID type"})
				return
			}

			ctx.Set("userID", userIDInt)
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
	}
}
