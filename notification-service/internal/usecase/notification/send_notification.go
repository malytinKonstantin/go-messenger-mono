package notification

import (
	"context"
	"time"

	"github.com/gocql/gocql"
	"github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/repositories"
)

type SendNotificationUsecase interface {
	Execute(ctx context.Context, userID gocql.UUID, message string, notifType models.NotificationType) error
}

type sendNotificationUsecase struct {
	notificationRepo repositories.NotificationRepository
}

func NewSendNotificationUsecase(notificationRepo repositories.NotificationRepository) SendNotificationUsecase {
	return &sendNotificationUsecase{
		notificationRepo: notificationRepo,
	}
}

func (uc *sendNotificationUsecase) Execute(ctx context.Context, userID gocql.UUID, message string, notifType models.NotificationType) error {
	notification := &models.Notification{
		NotificationID: gocql.TimeUUID(),
		UserID:         userID,
		Message:        message,
		Type:           notifType,
		CreatedAt:      time.Now(),
		IsRead:         false,
	}

	return uc.notificationRepo.CreateNotification(ctx, notification)
}
