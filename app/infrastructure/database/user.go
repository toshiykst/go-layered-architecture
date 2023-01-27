package database

import (
	"gorm.io/gorm"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

type userDataModel struct {
	ID    string `gorm:"primaryKey"`
	Name  string
	Email string
}

func toUserModel(dm *userDataModel) *model.User {
	return &model.User{
		ID:    model.UserID(dm.ID),
		Name:  dm.Name,
		Email: dm.Email,
	}
}

func toUserDataModel(u *model.User) *userDataModel {
	return &userDataModel{
		ID:    string(u.ID),
		Name:  u.Name,
		Email: u.Email,
	}
}

type dbUserRepository struct {
	conn *gorm.DB
}

func (db *dbUserRepository) Find(id model.UserID) (*model.User, error) {
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
