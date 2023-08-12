package dto

import "github.com/toshiykst/go-layerd-architecture/app/domain/model"

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

func ToUsersFromModel(mus model.Users) []User {
	result := make([]User, len(mus))
	for i, mu := range mus {
		result[i] = User{
			UserID: string(mu.ID()),
			Name:   mu.Name(),
			Email:  mu.Email(),
		}
	}
	return result
}

func ToModelUserIDs(ids []string) []model.UserID {
	uIDs := make([]model.UserID, len(ids))
	for i, id := range ids {
		uIDs[i] = model.UserID(id)
	}
	return uIDs
}
