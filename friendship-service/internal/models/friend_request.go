package models

import (
	"time"

	"github.com/mindstand/gogm/v2"
)

type FriendRequest struct {
	gogm.BaseNode

	RequestID string    `gogm:"name=request_id;unique"`
	From      *User     `gogm:"direction=outgoing;relationship=FROM"`
	To        *User     `gogm:"direction=outgoing;relationship=TO"`
	SentAt    time.Time `gogm:"name=sent_at"`
	Status    string    `gogm:"name=status"`
}
