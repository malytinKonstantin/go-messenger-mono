package models

import "time"

type FriendRequest struct {
	RequestID  string    `json:"request_id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Status     string    `json:"status"` // Возможные значения: "pending", "accepted", "rejected"
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
