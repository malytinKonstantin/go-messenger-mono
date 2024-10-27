package password

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/infrastructure/database/model"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/repository"
)

type ResetPasswordUsecase interface {
	Execute(ctx context.Context, email string) (*model.ResetPasswordToken, error)
}

type resetPasswordUsecase struct {
	userRepo  repository.UserCredentialsRepository
	tokenRepo repository.ResetPasswordTokenRepository
}

func NewResetPasswordUsecase(userRepo repository.UserCredentialsRepository, tokenRepo repository.ResetPasswordTokenRepository) ResetPasswordUsecase {
	return &resetPasswordUsecase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (uc *resetPasswordUsecase) Execute(ctx context.Context, email string) (*model.ResetPasswordToken, error) {
	// Получаем пользователя по email
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Создаем токен для сброса пароля
	token := &model.ResetPasswordToken{
		Token:     uuid.New().String(),
		UserID:    user.UserID,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	err = uc.tokenRepo.Create(ctx, token)
	if err != nil {
		return nil, err
	}

	// Здесь можно отправить email с токеном для сброса пароля

	return token, nil
}
