package database

import (
	"errors"

	"gorm.io/gorm"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/infrastructure/database/datamodel"
)

type dbGroupRepository struct {
	db *gorm.DB
}

func (r *dbGroupRepository) Find(gID model.GroupID) (*model.Group, error) {
	dmg := &datamodel.Group{ID: string(gID)}

	if err := r.db.First(dmg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var dmgus datamodel.GroupUsers
	if err := r.db.Where("group_id = ?", gID).Find(&dmgus).Error; err != nil {
		return nil, err
	}

	return dmg.ToModel(dmgus), nil
}

func (r *dbGroupRepository) List() ([]*model.Group, error) {
	return nil, nil
}

func (r *dbGroupRepository) Create(g *model.Group) (*model.Group, error) {
	dmg := datamodel.NewGroup(g.ID(), g.Name())
	dmgus := datamodel.NewGroupUsers(g.ID(), g.UserIDs())

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.db.Create(dmg).Error; err != nil {
			return err
		}

		if len(dmgus) == 0 {
			return nil
		}

		if err := r.db.Create(dmgus).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return dmg.ToModel(dmgus), nil
}

func (r *dbGroupRepository) Update(g *model.Group) error {
	return nil
}

func (r *dbGroupRepository) Delete(gID model.GroupID) error {
	return nil
}
