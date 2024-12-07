package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/malytinKonstantin/go-messenger-mono/auth-service/infrastructure/database/model"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/repository"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/auth/mocks"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthenticateUser_Success(t *testing.T) {
	// Инициализация
	ctx := context.Background()
	userRepo := new(mocks.UserCredentialsRepository)
	jwtSecret := "test-secret"
	email := "test@example.com"
	password := "password123"
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	userID := uuid.New().String()

	// Настройка ожидаемого поведения моков
	userRepo.On("GetByEmail", mock.Anything, email).Return(&model.UserCredential{
		UserID:       userID,
		Email:        email,
		PasswordHash: string(passwordHash),
	}, nil)

	usecase := auth.NewAuthenticateUserUsecase(userRepo, jwtSecret)

	// Выполнение
	token, err := usecase.Execute(ctx, email, password)

	// Проверка
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Проверяем, что токен корректный
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})
	assert.NoError(t, err)
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, userID, claims["user_id"])
}

func TestAuthenticateUser_InvalidPassword(t *testing.T) {
	// Инициализация
	ctx := context.Background()
	userRepo := new(mocks.UserCredentialsRepository)
	jwtSecret := "test-secret"
	email := "test@example.com"
	password := "password123"
	wrongPassword := "wrongpassword"
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Настройка ожидаемого поведения моков
	userRepo.On("GetByEmail", mock.Anything, email).Return(&model.UserCredential{
		UserID:       uuid.New().String(),
		Email:        email,
		PasswordHash: string(passwordHash),
	}, nil)

	usecase := auth.NewAuthenticateUserUsecase(userRepo, jwtSecret)

	// Выполнение
	token, err := usecase.Execute(ctx, email, wrongPassword)

	// Проверка
	assert.Error(t, err)
	assert.Equal(t, "invalid password", err.Error())
	assert.Empty(t, token)
}

func TestAuthenticateUser_UserNotFound(t *testing.T) {
	// Инициализация
	ctx := context.Background()
	userRepo := new(mocks.UserCredentialsRepository)
	jwtSecret := "test-secret"
	email := "nonexistent@example.com"
	password := "password123"

	// Настройка ожидаемого поведения моков
	userRepo.On("GetByEmail", mock.Anything, email).Return(nil, repository.ErrNotFound)

	usecase := auth.NewAuthenticateUserUsecase(userRepo, jwtSecret)

	// Выполнение
	token, err := usecase.Execute(ctx, email, password)

	// Проверка
	assert.Error(t, err)
	assert.Equal(t, repository.ErrNotFound, err)
	assert.Empty(t, token)
}
