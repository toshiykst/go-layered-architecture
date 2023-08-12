package response

import "github.com/toshiykst/go-layerd-architecture/app/usecase/dto"

type User struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type Group struct {
	GroupID string `json:"groupId"`
	Name    string `json:"name"`
	Users   []User `json:"users"`
}

func ToUsersFromDTO(dtous []dto.User) []User {
	us := make([]User, len(dtous))
	for i, dtou := range dtous {
		us[i] = User{
			UserID: dtou.UserID,
			Name:   dtou.Name,
			Email:  dtou.Email,
		}
	}
	return us
}
