package database

import (
	"errors"

	"gorm.io/gorm"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/infrastructure/database/datamodel"
)

type dbUserRepository struct {
	db *gorm.DB
}

func (r *dbUserRepository) Find(uID model.UserID) (*model.User, error) {
	dmu := &datamodel.User{ID: string(uID)}

	if err := r.db.First(dmu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return datamodel.ToUserModel(dmu), nil
}

func (r *dbUserRepository) FindByName(name string) (*model.User, error) {
	return nil, nil
}

func (r *dbUserRepository) List() ([]*model.User, error) {
	return nil, nil
}

func (r *dbUserRepository) Create(u *model.User) (*model.User, error) {
	dmu := datamodel.NewUser(u.ID(), u.Name(), u.Email())

	if err := r.db.Create(dmu).Error; err != nil {
		return nil, err
	}

	return datamodel.ToUserModel(dmu), nil
}

func (r *dbUserRepository) Update(u *model.User) error {
	return nil
}

func (r *dbUserRepository) Delete(id model.UserID) error {
	return nil
}
