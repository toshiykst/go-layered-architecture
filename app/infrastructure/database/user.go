package database

import (
	"errors"

	"gorm.io/gorm"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/domain/repository"
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

	return dmu.ToModel(), nil
}

func (r *dbUserRepository) List(f repository.UserListFilter) (model.Users, error) {
	db := r.db
	if len(f.UserIDs) > 0 {
		db = db.Where("id IN (?)", f.UserIDs)
	}

	var dmus datamodel.Users
	if err := db.Find(&dmus).Error; err != nil {
		return nil, err
	}

	return dmus.ToModel(), nil
}

func (r *dbUserRepository) Create(u *model.User) (*model.User, error) {
	dmu := datamodel.NewUser(u.ID(), u.Name(), u.Email())

	if err := r.db.Create(dmu).Error; err != nil {
		return nil, err
	}

	return dmu.ToModel(), nil
}

func (r *dbUserRepository) Update(u *model.User) error {
	if u.ID() == "" {
		return errors.New("user id must not be empty")
	}

	if err := r.db.Model(&datamodel.User{ID: string(u.ID())}).
		Updates(map[string]any{
			"name":  u.Name(),
			"email": u.Email(),
		}).Error; err != nil {
		return err
	}

	return nil
}

func (r *dbUserRepository) Delete(uID model.UserID) error {
	return r.db.Delete(&datamodel.User{
		ID: string(uID),
	}).Error
}
