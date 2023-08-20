package usecase

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrGroupNotFound     = errors.New("group not found")
	ErrInvalidUserInput  = errors.New("invalid user input")
	ErrInvalidGroupInput = errors.New("invalid group input")
	ErrInvalidUserIDs    = errors.New("invalid user ids")
)
