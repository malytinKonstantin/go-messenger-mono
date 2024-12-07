package models

import (
	"github.com/gocql/gocql"
)

type NotificationPreferences struct {
	UserID        gocql.UUID
	NewMessage    bool
	FriendRequest bool
	System        bool
}
