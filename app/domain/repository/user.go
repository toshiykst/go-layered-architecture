package repository

import (
	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

type UserListFilter struct {
	UserIDs []model.UserID
}

// UserRepositoryQuery is interface for query methods of user.
type UserRepositoryQuery interface {
	Find(uID model.UserID) (*model.User, error)
	List(f UserListFilter) (model.Users, error)
}

// UserRepositoryCommand is interface for query and command methods of user.
type UserRepositoryCommand interface {
	UserRepositoryQuery
	Create(u *model.User) (*model.User, error)
	Update(u *model.User) error
	Delete(uID model.UserID) error
}
