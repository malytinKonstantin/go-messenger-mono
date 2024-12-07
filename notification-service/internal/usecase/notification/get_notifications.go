package notification

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/repositories"
)

type GetNotificationsUsecase interface {
	Execute(ctx context.Context, userID gocql.UUID, limit int, offset int) ([]*models.Notification, int, error)
}

type getNotificationsUsecase struct {
	notificationRepo repositories.NotificationRepository
}

func NewGetNotificationsUsecase(notificationRepo repositories.NotificationRepository) GetNotificationsUsecase {
	return &getNotificationsUsecase{
		notificationRepo: notificationRepo,
	}
}

func (uc *getNotificationsUsecase) Execute(ctx context.Context, userID gocql.UUID, limit int, offset int) ([]*models.Notification, int, error) {
	return uc.notificationRepo.GetNotifications(ctx, userID, limit, offset)
}
