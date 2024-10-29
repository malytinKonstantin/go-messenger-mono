package user

import "errors"

var (
	ErrNicknameAlreadyExists = errors.New("nickname already exists")
)
