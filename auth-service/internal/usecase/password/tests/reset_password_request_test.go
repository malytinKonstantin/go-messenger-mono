package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/infrastructure/database/model"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/repository"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/password"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/password/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResetPasswordRequest_Success(t *testing.T) {
	ctx := context.Background()
	userRepo := new(mocks.UserCredentialsRepository)
	tokenRepo := new(mocks.ResetPasswordTokenRepository)

	email := "user@example.com"
	userID := uuid.New().String()
	user := &model.UserCredential{
		UserID: userID,
		Email:  email,
	}

	// Настройка поведения моков
	userRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)
	tokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.ResetPasswordToken")).Return(nil)

	usecase := password.NewResetPasswordRequestUsecase(userRepo, tokenRepo)

	// Выполнение
	token, err := usecase.Execute(ctx, email)

	// Проверка
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, userID, token.UserID)

	userRepo.AssertCalled(t, "GetByEmail", mock.Anything, email)
	tokenRepo.AssertCalled(t, "Create", mock.Anything, mock.AnythingOfType("*model.ResetPasswordToken"))
}

func TestResetPasswordRequest_UserNotFound(t *testing.T) {
	ctx := context.Background()
	userRepo := new(mocks.UserCredentialsRepository)
	tokenRepo := new(mocks.ResetPasswordTokenRepository)

	email := "nonexistent@example.com"
	expectedError := repository.ErrNotFound

	// Настройка поведения моков
	userRepo.On("GetByEmail", mock.Anything, email).Return(nil, expectedError)

	usecase := password.NewResetPasswordRequestUsecase(userRepo, tokenRepo)

	// Выполнение
	token, err := usecase.Execute(ctx, email)

	// Проверка
	assert.Error(t, err)
	assert.Nil(t, token)
	assert.Equal(t, expectedError, err)

	userRepo.AssertCalled(t, "GetByEmail", mock.Anything, email)
	tokenRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestResetPasswordRequest_TokenCreationError(t *testing.T) {
	ctx := context.Background()
	userRepo := new(mocks.UserCredentialsRepository)
	tokenRepo := new(mocks.ResetPasswordTokenRepository)

	email := "user@example.com"
	userID := uuid.New().String()
	user := &model.UserCredential{
		UserID: userID,
		Email:  email,
	}
	expectedError := errors.New("database error")

	// Настройка поведения моков
	userRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)
	tokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.ResetPasswordToken")).Return(expectedError)

	usecase := password.NewResetPasswordRequestUsecase(userRepo, tokenRepo)

	// Выполнение
	token, err := usecase.Execute(ctx, email)

	// Проверка
	assert.Error(t, err)
	assert.Nil(t, token)
	assert.Equal(t, expectedError, err)

	userRepo.AssertCalled(t, "GetByEmail", mock.Anything, email)
	tokenRepo.AssertCalled(t, "Create", mock.Anything, mock.AnythingOfType("*model.ResetPasswordToken"))
}
