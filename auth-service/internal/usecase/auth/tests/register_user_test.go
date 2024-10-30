package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/repository"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/malytinKonstantin/go-messenger-mono/auth-service/infrastructure/database/model"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/auth/mocks"
)

func TestRegisterUser_Success(t *testing.T) {
	// Инициализация
	ctx := context.Background()
	userRepo := new(mocks.UserCredentialsRepository)
	email := "newuser@example.com"
	password := "password123"

	// Настройка ожидаемого поведения моков
	userRepo.On("GetByEmail", mock.Anything, email).Return(nil, repository.ErrNotFound)
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.UserCredential")).Return(nil)

	usecase := auth.NewRegisterUserUsecase(userRepo)

	// Выполнение
	user, err := usecase.Execute(ctx, email, password)

	// Проверка
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.False(t, user.IsVerified)
	assert.NotEmpty(t, user.PasswordHash)
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	// Инициализация
	ctx := context.Background()
	userRepo := new(mocks.UserCredentialsRepository)
	email := "existinguser@example.com"
	password := "password123"

	existingUser := &model.UserCredential{
		UserID: uuid.New().String(),
		Email:  email,
	}

	// Настройка ожидаемого поведения моков
	userRepo.On("GetByEmail", mock.Anything, email).Return(existingUser, nil)

	usecase := auth.NewRegisterUserUsecase(userRepo)

	// Выполнение
	user, err := usecase.Execute(ctx, email, password)

	// Проверка
	assert.Error(t, err)
	assert.Equal(t, "user with this email already exists", err.Error())
	assert.Nil(t, user)
}
