package repositories

import (
	"context"

	"github.com/malytinKonstantin/go-messenger-mono/notification-service/internal/models"

	"github.com/gocql/gocql"
)

type NotificationPreferencesRepository interface {
	CreatePreferences(ctx context.Context, preferences *models.NotificationPreferences) error
	GetPreferences(ctx context.Context, userID gocql.UUID) (*models.NotificationPreferences, error)
	UpdatePreferences(ctx context.Context, preferences *models.NotificationPreferences) error
	DeletePreferences(ctx context.Context, userID gocql.UUID) error
}

type notificationPreferencesRepository struct {
	session *gocql.Session
}

func NewNotificationPreferencesRepository(session *gocql.Session) NotificationPreferencesRepository {
	return &notificationPreferencesRepository{
		session: session,
	}
}

func (r *notificationPreferencesRepository) CreatePreferences(ctx context.Context, preferences *models.NotificationPreferences) error {
	query := `INSERT INTO notification_preferences (
        user_id, new_message, friend_request, system
    ) VALUES (?, ?, ?, ?)`
	return r.session.Query(query,
		preferences.UserID,
		preferences.NewMessage,
		preferences.FriendRequest,
		preferences.System,
	).WithContext(ctx).Exec()
}

func (r *notificationPreferencesRepository) GetPreferences(ctx context.Context, userID gocql.UUID) (*models.NotificationPreferences, error) {
	query := `SELECT new_message, friend_request, system FROM notification_preferences WHERE user_id = ?`
	var preferences models.NotificationPreferences
	err := r.session.Query(query, userID).WithContext(ctx).Scan(
		&preferences.NewMessage,
		&preferences.FriendRequest,
		&preferences.System,
	)
	if err != nil {
		return nil, err
	}
	preferences.UserID = userID
	return &preferences, nil
}

func (r *notificationPreferencesRepository) UpdatePreferences(ctx context.Context, preferences *models.NotificationPreferences) error {
	query := `UPDATE notification_preferences SET new_message = ?, friend_request = ?, system = ? WHERE user_id = ?`
	return r.session.Query(query,
		preferences.NewMessage,
		preferences.FriendRequest,
		preferences.System,
		preferences.UserID,
	).WithContext(ctx).Exec()
}

func (r *notificationPreferencesRepository) DeletePreferences(ctx context.Context, userID gocql.UUID) error {
	query := `DELETE FROM notification_preferences WHERE user_id = ?`
	return r.session.Query(query, userID).WithContext(ctx).Exec()
}
