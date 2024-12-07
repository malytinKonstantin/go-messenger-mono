package preferences

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/repositories"
)

type GetPreferencesUsecase interface {
	Execute(ctx context.Context, userID gocql.UUID) (*models.NotificationPreferences, error)
}

type getPreferencesUsecase struct {
	preferencesRepo repositories.NotificationPreferencesRepository
}

func NewGetPreferencesUsecase(preferencesRepo repositories.NotificationPreferencesRepository) GetPreferencesUsecase {
	return &getPreferencesUsecase{
		preferencesRepo: preferencesRepo,
	}
}

func (uc *getPreferencesUsecase) Execute(ctx context.Context, userID gocql.UUID) (*models.NotificationPreferences, error) {
	return uc.preferencesRepo.GetPreferences(ctx, userID)
}
