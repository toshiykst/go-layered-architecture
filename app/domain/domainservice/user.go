//go:generate mockgen -source=$GOFILE -package=mock$GOPACKAGE -destination=../../mock/domain/$GOPACKAGE/$GOFILE
package domainservice

import (
	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/domain/repository"
)

type UserService interface {
	Exists(uID model.UserID) (bool, error)
}

type userService struct {
	r repository.Repository
}

func NewUserService(r repository.Repository) *userService {
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
