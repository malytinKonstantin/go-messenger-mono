package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/repository"
)

type VerifyEmailUsecase interface {
	Execute(ctx context.Context, userID uuid.UUID) error
}

type verifyEmailUsecase struct {
	userRepo repository.UserCredentialsRepository
}

func NewVerifyEmailUsecase(userRepo repository.UserCredentialsRepository) VerifyEmailUsecase {
	return &verifyEmailUsecase{userRepo: userRepo}
}

func (uc *verifyEmailUsecase) Execute(ctx context.Context, userID uuid.UUID) error {
	// Обновляем статус верификации пользователя
	err := uc.userRepo.UpdateVerificationStatus(ctx, userID, true)
	if err != nil {
		return err
	}
	return nil
}
