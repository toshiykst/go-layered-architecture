package model

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidUser = errors.New("invalid user")
)

const (
	maxUserNameLength  = 30
	maxUserEmailLength = 254
)

type UserID string

type User struct {
	id    UserID
	name  string
	email string
}

func NewUser(id UserID, name, email string) (*User, error) {
	if id == "" {
		return nil, fmt.Errorf("user id must not empty: %w", ErrInvalidUser)
	}

	if name == "" {
		return nil, fmt.Errorf("user name must not empty: %w", ErrInvalidUser)
	}
	if len(name) > maxUserNameLength {
		return nil, fmt.Errorf("exceeds the max user name length: %w", ErrInvalidUser)
	}

	if email == "" {
		return nil, fmt.Errorf("user email must not empty: %w", ErrInvalidUser)
	}
	if len(email) > maxUserEmailLength {
		return nil, fmt.Errorf("exceeds the max user name length: %w", ErrInvalidUser)
	}

	return &User{
		id:    id,
		name:  name,
		email: email,
	}, nil
}

func MustNewUser(id UserID, name, email string) *User {
	u, err := NewUser(id, name, email)
	if err != nil {
		panic(err)
	}
	return u
}

func (u *User) ID() UserID {
	if u == nil {
		return ""
	}
	return u.id
}

func (u *User) Name() string {
	if u == nil {
		return ""
	}
	return u.name
}

func (u *User) Email() string {
	if u == nil {
		return ""
	}
	return u.email
}

type Users []*User

func (us Users) ByUserID() map[UserID]*User {
	result := make(map[UserID]*User, len(us))
	for _, u := range us {
		result[u.ID()] = u
	}
	return result
}
