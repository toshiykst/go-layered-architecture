package factory

import (
	"github.com/google/uuid"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

type UserFactory interface {
	Create(name, email string) (*model.User, error)
}

type userFactory struct{}

func NewUserFactory() UserFactory {
	return &userFactory{}
}

func (uf userFactory) Create(name, email string) (*model.User, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return model.NewUser(model.UserID(uuid.String()), name, email), nil
}
