package models

import (
	"time"

	"github.com/google/uuid"
)

type UserID uuid.UUID

type UserCredentials struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       UserID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"user_id"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"password_hash"`
	IsVerified   bool      `gorm:"not null;default:false" json:"is_verified"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (UserCredentials) TableName() string {
	return "user_credentials"
}
