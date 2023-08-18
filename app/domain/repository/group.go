package repository

import (
	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

type GroupListFilter struct {
	UserIDs []model.UserID
}

// GroupRepositoryQuery is interface for query methods of group.
type GroupRepositoryQuery interface {
	Find(gID model.GroupID) (*model.Group, error)
	List(f GroupListFilter) (model.Groups, error)
}

// GroupRepositoryCommand is interface for query and command methods of group.
type GroupRepositoryCommand interface {
	GroupRepositoryQuery
	Create(g *model.Group) (*model.Group, error)
	Update(g *model.Group) error
	Delete(gID model.GroupID) error
	AddUsers(gID model.GroupID, uIDs []model.UserID) error
	RemoveUsers(gID model.GroupID, uIDs []model.UserID) error
	RemoveUsersFromAll(uIDs []model.UserID) error
}
