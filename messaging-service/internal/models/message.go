package models

import (
	"time"

	"github.com/gocql/gocql"
)

type MessageStatus int32

const (
	StatusUnspecified MessageStatus = 0
	StatusSent        MessageStatus = 1
	StatusDelivered   MessageStatus = 2
	StatusRead        MessageStatus = 3
)

type Message struct {
	MessageID      gocql.UUID
	SenderID       gocql.UUID
	RecipientID    gocql.UUID
	ConversationID gocql.UUID
	Content        string
	Status         MessageStatus
	Timestamp      time.Time
}
