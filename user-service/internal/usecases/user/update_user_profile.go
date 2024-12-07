package user

import (
	"context"
	"time"

	"github.com/malytinKonstantin/go-messenger-mono/user-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/user-service/internal/repositories"
)

type UpdateUserProfileUsecase interface {
	Execute(ctx context.Context, profile *models.UserProfile) error
}

type updateUserProfileUsecase struct {
	userRepo repositories.UserRepository
}

func NewUpdateUserProfileUsecase(userRepo repositories.UserRepository) UpdateUserProfileUsecase {
	return &updateUserProfileUsecase{
		userRepo: userRepo,
	}
}

func (uc *updateUserProfileUsecase) Execute(ctx context.Context, profile *models.UserProfile) error {
	// Проверка существования пользователя
	existingProfile, err := uc.userRepo.GetUser(ctx, profile.UserID)
	if err != nil {
		return err
	}

	// Если никнейм меняется, проверить его уникальность
	if profile.Nickname != "" && profile.Nickname != existingProfile.Nickname {
		exists, err := uc.userRepo.NicknameExists(ctx, profile.Nickname)
		if err != nil {
			return err
		}
		if exists {
			return ErrNicknameAlreadyExists
		}
	}

	// Обновление полей профиля
	if profile.Nickname != "" {
		existingProfile.Nickname = profile.Nickname
	}
	if profile.Bio != "" {
		existingProfile.Bio = profile.Bio
	}
	if profile.AvatarURL != "" {
		existingProfile.AvatarURL = profile.AvatarURL
	}
	existingProfile.UpdatedAt = time.Now()

	return uc.userRepo.UpdateUserProfile(ctx, existingProfile)
}
