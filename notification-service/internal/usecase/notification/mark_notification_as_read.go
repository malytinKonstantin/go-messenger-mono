package notification

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/repositories"
)

type MarkNotificationAsReadUsecase interface {
	Execute(ctx context.Context, notificationID, userID gocql.UUID) error
}

type markNotificationAsReadUsecase struct {
	notificationRepo repositories.NotificationRepository
}

func NewMarkNotificationAsReadUsecase(notificationRepo repositories.NotificationRepository) MarkNotificationAsReadUsecase {
	return &markNotificationAsReadUsecase{
		notificationRepo: notificationRepo,
	}
}

func (uc *markNotificationAsReadUsecase) Execute(ctx context.Context, notificationID, userID gocql.UUID) error {
	return uc.notificationRepo.MarkNotificationAsRead(ctx, notificationID, userID)
}
