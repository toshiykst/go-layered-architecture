package domainservice

import (
	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/domain/repository"
)

type UserService struct {
	r repository.Repository
}

func NewUserService(r repository.Repository) *UserService {
	return &UserService{r: r}
}

func (us *UserService) Exists(uID model.UserID) (bool, error) {
	u, err := us.r.User().Find(uID)
	if err != nil {
		return false, err
	}
	if u == nil {
		return false, nil
	}
	return true, nil
}
