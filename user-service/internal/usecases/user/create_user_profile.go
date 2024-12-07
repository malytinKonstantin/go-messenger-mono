package user

import (
	"context"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/malytinKonstantin/go-messenger-mono/user-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/user-service/internal/repositories"
)

type CreateUserProfileUsecase interface {
	Execute(ctx context.Context, profile *models.UserProfile) error
}

type createUserProfileUsecase struct {
	userRepo repositories.UserRepository
}

func NewCreateUserProfileUsecase(userRepo repositories.UserRepository) CreateUserProfileUsecase {
	return &createUserProfileUsecase{
		userRepo: userRepo,
	}
}

func (uc *createUserProfileUsecase) Execute(ctx context.Context, profile *models.UserProfile) error {
	// Проверка уникальности никнейма
	exists, err := uc.userRepo.NicknameExists(ctx, profile.Nickname)
	if err != nil {
		return err
	}
	if exists {
		return ErrNicknameAlreadyExists
	}

	// Установка идентификатора и временных меток
	profile.UserID = gocql.UUID(uuid.New())
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	return uc.userRepo.CreateUserProfile(ctx, profile)
}
