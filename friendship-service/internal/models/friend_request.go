package models

type FriendRequest struct {
	RequestID  string `json:"request_id"`
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Status     string `json:"status"` // Возможные значения: "pending", "accepted", "rejected"
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}
