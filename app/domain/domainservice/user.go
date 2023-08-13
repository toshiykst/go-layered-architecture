package domainservice

import (
	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/domain/repository"
)

type UserService interface {
	Exists(uID model.UserID) (bool, error)
	ExistsAll(uIDs []model.UserID) (bool, error)
}

type userService struct {
	r repository.Repository
}

func NewUserService(r repository.Repository) UserService {
	return &userService{r: r}
}

func (us *userService) Exists(uID model.UserID) (bool, error) {
	u, err := us.r.User().Find(uID)
	if err != nil {
		return false, err
	}
	if u == nil {
		return false, nil
	}
	return true, nil
}

func (us *userService) ExistsAll(uIDs []model.UserID) (bool, error) {
	users, err := us.r.User().List(repository.UserListFilter{
		UserIDs: uIDs,
	})
	if err != nil {
		return false, err
	}

	if len(users) != len(uIDs) {
		return false, nil
	}

	return true, nil
}
