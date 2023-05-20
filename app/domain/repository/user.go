package repository

import (
	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

// UserRepositoryQuery is interface for query methods of user.
type UserRepositoryQuery interface {
	Find(id model.UserID) (*model.User, error)
	List() (model.Users, error)
}

// UserRepositoryCommand is interface for query and command methods of user.
type UserRepositoryCommand interface {
	UserRepositoryQuery
	Create(u *model.User) (*model.User, error)
	Update(u *model.User) error
	Delete(id model.UserID) error
}
