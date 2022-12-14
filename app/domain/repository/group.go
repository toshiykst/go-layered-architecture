package repository

import (
	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

// GroupRepositoryQuery is interface for query methods of group.
type GroupRepositoryQuery interface {
	Find(id model.GroupID) (*model.Group, error)
	FindByName(name string) (*model.Group, error)
	List() ([]*model.Group, error)
}

// GroupRepositoryCommand is interface for query and command methods of group.
type GroupRepositoryCommand interface {
	GroupRepositoryQuery
	Create(g *model.Group) (*model.Group, error)
	Update(g *model.Group) error
	Delete(id model.GroupID) error
}
