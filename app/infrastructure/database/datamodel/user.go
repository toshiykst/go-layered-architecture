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

func (u *User) ToModel() *model.User {
	if u == nil {
		return nil
	}
	return model.NewUser(
		model.UserID(u.ID),
		u.Name,
		u.Email,
	)
}

type Users []*User

func (us Users) ToModel() model.Users {
	if us == nil {
		return nil
	}
	mus := make(model.Users, len(us))
	for i, u := range us {
		mus[i] = model.NewUser(
			model.UserID(u.ID),
			u.Name,
			u.Email,
		)
	}
	return mus
}
