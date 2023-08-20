//go:generate mockgen -source=$GOFILE -package=mock$GOPACKAGE -destination=../../mock/domain/$GOPACKAGE/$GOFILE
package factory

import (
	"github.com/google/uuid"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

type GroupFactory interface {
	Create(name string, uIDs []model.UserID) (*model.Group, error)
}

type groupFactory struct{}

func NewGroupFactory() GroupFactory {
	return &groupFactory{}
}

func (uf groupFactory) Create(name string, uIDs []model.UserID) (*model.Group, error) {
	generated, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	g, err := model.NewGroup(model.GroupID(generated.String()), name, uIDs)
	if err != nil {
		return nil, err
	}

	return g, nil
}
