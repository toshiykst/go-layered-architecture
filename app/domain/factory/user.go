//go:generate mockgen -source=$GOFILE -package=mock$GOPACKAGE -destination=../../mock/domain/$GOPACKAGE/$GOFILE
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
	generated, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	u, err := model.NewUser(model.UserID(generated.String()), name, email)
	if err != nil {
		return nil, err
	}

	return u, nil
}
