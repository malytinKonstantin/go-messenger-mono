package models

import (
	"time"

	"github.com/gocql/gocql"
)

type NotificationType int32

const (
	NotificationTypeUnspecified   NotificationType = 0
	NotificationTypeNewMessage    NotificationType = 1
	NotificationTypeFriendRequest NotificationType = 2
	NotificationTypeSystem        NotificationType = 3
)

type Notification struct {
	NotificationID gocql.UUID
	UserID         gocql.UUID
	Message        string
	Type           NotificationType
	CreatedAt      time.Time
	IsRead         bool
}
