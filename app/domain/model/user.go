package model

type UserID string

type User struct {
	id    UserID
	name  string
	email string
}

func NewUser(id UserID, name, email string) *User {
	return &User{
		id:    id,
		name:  name,
		email: email,
	}
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
