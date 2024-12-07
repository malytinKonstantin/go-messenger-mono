package preferences

import (
	"context"

	"github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/repositories"
)

type UpdatePreferencesUsecase interface {
	Execute(ctx context.Context, preferences *models.NotificationPreferences) error
}

type updatePreferencesUsecase struct {
	preferencesRepo repositories.NotificationPreferencesRepository
}

func NewUpdatePreferencesUsecase(preferencesRepo repositories.NotificationPreferencesRepository) UpdatePreferencesUsecase {
	return &updatePreferencesUsecase{
		preferencesRepo: preferencesRepo,
	}
}

func (uc *updatePreferencesUsecase) Execute(ctx context.Context, preferences *models.NotificationPreferences) error {
	return uc.preferencesRepo.UpdatePreferences(ctx, preferences)
}
