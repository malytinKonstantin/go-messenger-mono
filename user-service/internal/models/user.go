// Пакет models содержит модели данных для сервиса пользователей.
package models

import (
	"time"

	"github.com/gocql/gocql"
)

// UserProfile представляет профиль пользователя.
type UserProfile struct {
	UserID    gocql.UUID `db:"user_id"`
	Nickname  string     `db:"nickname"`
	Bio       string     `db:"bio"`
	AvatarURL string     `db:"avatar_url"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
}
