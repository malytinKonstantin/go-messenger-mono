package models

type User struct {
	UserID    string   `json:"user_id"`
	Nickname  string   `json:"nickname"`
	AvatarURL string   `json:"avatar_url"`
	Friends   []string `json:"friends"`
}
