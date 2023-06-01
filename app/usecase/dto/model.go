package dto

type User struct {
	UserID string
	Name   string
	Email  string
}

type Group struct {
	GroupID string
	Name    string
	Users   []User
}
