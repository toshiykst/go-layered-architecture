package usecase

import (
	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/usecase/dto"
)

func convertUsersToDTO(us model.Users) []dto.User {
	result := make([]dto.User, len(us))
	for i, u := range us {
		result[i] = dto.User{
			UserID: string(u.ID()),
			Name:   u.Name(),
			Email:  u.Email(),
		}
	}
	return result
}
