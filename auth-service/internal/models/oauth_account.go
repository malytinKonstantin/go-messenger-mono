package models

import (
	"time"
)

type OauthAccount struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         UserID    `gorm:"type:uuid;not null;index" json:"user_id"`
	Provider       string    `gorm:"type:varchar(50);not null" json:"provider"`
	ProviderUserID string    `gorm:"type:varchar(255);not null" json:"provider_user_id"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (OauthAccount) TableName() string {
	return "oauth_accounts"
}
