package repositories

import (
	"context"
	"errors"

	"github.com/gocql/gocql"
	"github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/shared/platform/cassandra"
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
	session cassandra.Session
}

func NewNotificationRepository(session cassandra.Session) NotificationRepository {
	return &notificationRepository{
		session: session,
	}
}

func (r *notificationRepository) CreateNotification(ctx context.Context, notification *models.Notification) error {
	query := `INSERT INTO notifications (notification_id, user_id, message, type, is_read, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	return r.session.Query(query,
		notification.NotificationID,
		notification.UserID,
		notification.Message,
		int32(notification.Type),
		notification.IsRead,
		notification.CreatedAt,
	).WithContext(ctx).Exec()
}

func (r *notificationRepository) GetNotificationByID(ctx context.Context, notificationID, userID gocql.UUID) (*models.Notification, error) {
	query := `SELECT notification_id, user_id, message, type, is_read, created_at FROM notifications WHERE notification_id = ? AND user_id = ?`
	var n models.Notification
	if err := r.session.Query(query, notificationID, userID).WithContext(ctx).Scan(
		&n.NotificationID,
		&n.UserID,
		&n.Message,
		&n.Type,
		&n.IsRead,
		&n.CreatedAt,
	); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &n, nil
}

func (r *notificationRepository) GetNotifications(ctx context.Context, userID gocql.UUID, limit int, offset int) ([]*models.Notification, int, error) {
	query := `SELECT notification_id, user_id, message, type, is_read, created_at FROM notifications WHERE user_id = ?`
	iter := r.session.Query(query, userID).WithContext(ctx).PageSize(limit + offset).Iter()

	var notifications []*models.Notification
	var n models.Notification
	for iter.Scan(
		&n.NotificationID,
		&n.UserID,
		&n.Message,
		&n.Type,
		&n.IsRead,
		&n.CreatedAt,
	) {
		nCopy := n
		notifications = append(notifications, &nCopy)
	}

	if err := iter.Close(); err != nil {
		return nil, 0, err
	}

	total := len(notifications)
	if offset > total {
		return []*models.Notification{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return notifications[offset:end], total, nil
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
