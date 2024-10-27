package password

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type ChangePasswordUsecase interface {
	Execute(ctx context.Context, token uuid.UUID, newPassword string) error
}

type changePasswordUsecase struct {
	userRepo  repository.UserCredentialsRepository
	tokenRepo repository.ResetPasswordTokenRepository
}

func NewChangePasswordUsecase(userRepo repository.UserCredentialsRepository, tokenRepo repository.ResetPasswordTokenRepository) ChangePasswordUsecase {
	return &changePasswordUsecase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (uc *changePasswordUsecase) Execute(ctx context.Context, token uuid.UUID, newPassword string) error {
	// Получаем токен сброса пароля
	resetToken, err := uc.tokenRepo.GetByToken(ctx, token.String())
	if err != nil {
		return err
	}

	// Проверяем, не истек ли токен
	if resetToken.ExpiresAt.Before(time.Now()) {
		return errors.New("password reset token has expired")
	}

	// Хешируем новый пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userID, err := uuid.Parse(resetToken.UserID)
	if err != nil {
		return err
	}

	// Обновляем пароль пользователя
	err = uc.userRepo.UpdatePassword(ctx, userID, string(hashedPassword))
	if err != nil {
		return err
	}

	// Удаляем использованный токен
	err = uc.tokenRepo.Delete(ctx, token.String())
	if err != nil {
		return err
	}

	return nil
}
