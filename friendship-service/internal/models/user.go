package models

type User struct {
	UserID    string   `json:"user_id"`
	Nickname  string   `json:"nickname"`
	AvatarURL string   `json:"avatar_url"`
	AddedAt   int64    `json:"added_at"`
	Friends   []string `json:"friends"`
}
