package database

import (
	"gorm.io/gorm"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/infrastructure/database/datamodel"
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
	dmu := datamodel.NewUser(u.ID, u.Name, u.Name)

	if err := db.conn.Create(dmu).Error; err != nil {
		return nil, err
	}

	return datamodel.ToUserModel(dmu), nil
}

func (db *dbUserRepository) Update(u *model.User) error {
	return nil
}

func (db *dbUserRepository) Delete(id model.UserID) error {
	return nil
}
