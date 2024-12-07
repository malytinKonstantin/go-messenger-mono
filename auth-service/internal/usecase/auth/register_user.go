package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/infrastructure/database/model"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserUsecase interface {
	Execute(ctx context.Context, email, password string) (*model.UserCredential, error)
}

type registerUserUsecase struct {
	userRepo repository.UserCredentialsRepository
}

func NewRegisterUserUsecase(userRepo repository.UserCredentialsRepository) RegisterUserUsecase {
	return &registerUserUsecase{userRepo: userRepo}
}

func (uc *registerUserUsecase) Execute(ctx context.Context, email, password string) (*model.UserCredential, error) {
	// Проверяем, существует ли пользователь с таким email
	existingUser, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаем нового пользователя
	user := &model.UserCredential{
		UserID:       uuid.New().String(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		IsVerified:   false,
	}

	// Сохраняем пользователя в базе данных
	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
