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

func (r *dbGroupRepository) List() (model.Groups, error) {
	db := r.db

	var dmgs datamodel.Groups
	if err := db.Find(&dmgs).Error; err != nil {
		return nil, err
	}
	if len(dmgs) == 0 {
		return nil, nil
	}

	var dmgus datamodel.GroupUsers
	if err := r.db.Where("group_id IN (?)", dmgs.IDs()).Find(&dmgus).Error; err != nil {
		return nil, err
	}

	return dmgs.ToModel(dmgus), nil
}

func (r *dbGroupRepository) Create(g *model.Group) (*model.Group, error) {
	dmg := datamodel.NewGroup(g.ID(), g.Name())

	if err := r.db.Create(dmg).Error; err != nil {
		return nil, err
	}

	return dmg.ToModel(datamodel.GroupUsers{}), nil
}

func (r *dbGroupRepository) Update(g *model.Group) error {
	if g.ID() == "" {
		return errors.New("group id must not be empty")
	}
	if err := r.db.Model(&datamodel.Group{ID: string(g.ID())}).
		Updates(map[string]any{
			"name": g.Name(),
		}).Error; err != nil {
		return err
	}

	return nil
}

func (r *dbGroupRepository) Delete(gID model.GroupID) error {
	return nil
}

func (r *dbGroupRepository) AddUsers(gID model.GroupID, uIDs []model.UserID) error {
	dmgus := datamodel.NewGroupUsers(gID, uIDs)

	if err := r.db.Create(dmgus).Error; err != nil {
		return err
	}

	return nil
}

func (r *dbGroupRepository) RemoveUsers(gID model.GroupID, uIDs []model.UserID) error {
	return nil
}
