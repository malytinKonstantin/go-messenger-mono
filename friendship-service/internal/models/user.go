package models

import (
	"github.com/mindstand/gogm/v2"
)

type User struct {
	gogm.BaseNode

	UserID    string           `gogm:"name=user_id;unique"`
	Nickname  string           `gogm:"name=nickname"`
	AvatarURL string           `gogm:"name=avatar_url"`
	Friends   []*User          `gogm:"direction=both;relationship=FRIENDS_WITH"`
	Requests  []*FriendRequest `gogm:"direction=both;relationship=HAS_REQUEST"`
}
