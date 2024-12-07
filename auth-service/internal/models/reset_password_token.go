package models

import (
	"time"

	"github.com/google/uuid"
)

type ResetPasswordToken struct {
	Token     uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"token"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (ResetPasswordToken) TableName() string {
	return "reset_password_tokens"
}
