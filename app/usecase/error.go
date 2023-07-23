package usecase

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrGroupNotFound  = errors.New("group not found")
	ErrInvalidUserIDs = errors.New("invalid user ids")
)
