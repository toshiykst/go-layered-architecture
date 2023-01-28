package datamodel

import "github.com/toshiykst/go-layerd-architecture/app/domain/model"

type User struct {
	ID    string `gorm:"primaryKey"`
	Name  string
	Email string
}

func NewUser(uID model.UserID, name, email string) *User {
	return &User{
		ID:    string(uID),
		Name:  name,
		Email: email,
	}
}

func ToUserModel(u *User) *model.User {
	return &model.User{
		ID:    model.UserID(u.ID),
		Name:  u.Name,
		Email: u.Email,
	}
}
