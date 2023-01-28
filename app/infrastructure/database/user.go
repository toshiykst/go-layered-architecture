package database

import (
	"gorm.io/gorm"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

type dbUserRepository struct {
	conn *gorm.DB
}

func (db *dbUserRepository) Find(uID model.UserID) (*model.User, error) {
	return nil, nil
}

func (db *dbUserRepository) FindByName(name string) (*model.User, error) {
	return nil, nil
}

func (db *dbUserRepository) List() ([]*model.User, error) {
	return nil, nil
}

func (db *dbUserRepository) Create(u *model.User) (*model.User, error) {
	return nil, nil
}

func (db *dbUserRepository) Update(u *model.User) error {
	return nil
}

func (db *dbUserRepository) Delete(id model.UserID) error {
	return nil
}
