package repositories

import (
	"context"

	"github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/models"

	"github.com/gocql/gocql"
)

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification *models.Notification) error
	GetNotificationByID(ctx context.Context, notificationID, userID gocql.UUID) (*models.Notification, error)
	GetNotifications(ctx context.Context, userID gocql.UUID, limit int, offset int) ([]*models.Notification, int, error)
	UpdateNotification(ctx context.Context, notification *models.Notification) error
	DeleteNotification(ctx context.Context, notificationID, userID gocql.UUID) error
	MarkNotificationAsRead(ctx context.Context, notificationID, userID gocql.UUID) error
}

type notificationRepository struct {
	session *gocql.Session
}

func NewNotificationRepository(session *gocql.Session) NotificationRepository {
	return &notificationRepository{
		session: session,
	}
}

func (r *notificationRepository) CreateNotification(ctx context.Context, notification *models.Notification) error {
	query := `INSERT INTO notifications (
        notification_id, user_id, message, type, created_at, is_read
    ) VALUES (?, ?, ?, ?, ?, ?)`
	return r.session.Query(query,
		notification.NotificationID,
		notification.UserID,
		notification.Message,
		int32(notification.Type),
		notification.CreatedAt,
		notification.IsRead,
	).WithContext(ctx).Exec()
}

func (r *notificationRepository) GetNotificationByID(ctx context.Context, notificationID, userID gocql.UUID) (*models.Notification, error) {
	query := `SELECT notification_id, user_id, message, type, created_at, is_read
              FROM notifications WHERE notification_id = ? AND user_id = ?`
	var notification models.Notification
	err := r.session.Query(query, notificationID, userID).WithContext(ctx).Scan(
		&notification.NotificationID,
		&notification.UserID,
		&notification.Message,
		&notification.Type,
		&notification.CreatedAt,
		&notification.IsRead,
	)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) GetNotifications(ctx context.Context, userID gocql.UUID, limit int, offset int) ([]*models.Notification, int, error) {
	// Получаем общее количество уведомлений
	var totalCount int
	err := r.session.Query("SELECT COUNT(*) FROM notifications WHERE user_id = ?", userID).WithContext(ctx).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Формируем запрос с использованием LIMIT
	query := "SELECT notification_id, user_id, message, type, created_at, is_read FROM notifications WHERE user_id = ? ORDER BY created_at DESC LIMIT ?"
	iter := r.session.Query(query, userID, limit+offset).WithContext(ctx).Iter()

	var notifications []*models.Notification
	n := models.Notification{}
	count := 0
	for iter.Scan(&n.NotificationID, &n.UserID, &n.Message, &n.Type, &n.CreatedAt, &n.IsRead) {
		if count < offset {
			count++
			continue
		}
		notification := n // Создаем копию структуры
		notifications = append(notifications, &notification)
		if len(notifications) >= limit {
			break
		}
	}

	if err := iter.Close(); err != nil {
		return nil, 0, err
	}

	return notifications, totalCount, nil
}

func (r *notificationRepository) UpdateNotification(ctx context.Context, notification *models.Notification) error {
	query := `UPDATE notifications SET message = ?, type = ?, is_read = ? WHERE notification_id = ? AND user_id = ?`
	return r.session.Query(query,
		notification.Message,
		int32(notification.Type),
		notification.IsRead,
		notification.NotificationID,
		notification.UserID,
	).WithContext(ctx).Exec()
}

func (r *notificationRepository) DeleteNotification(ctx context.Context, notificationID, userID gocql.UUID) error {
	query := `DELETE FROM notifications WHERE notification_id = ? AND user_id = ?`
	return r.session.Query(query, notificationID, userID).WithContext(ctx).Exec()
}

func (r *notificationRepository) MarkNotificationAsRead(ctx context.Context, notificationID, userID gocql.UUID) error {
	query := `UPDATE notifications SET is_read = true WHERE notification_id = ? AND user_id = ?`
	return r.session.Query(query, notificationID, userID).WithContext(ctx).Exec()
}
