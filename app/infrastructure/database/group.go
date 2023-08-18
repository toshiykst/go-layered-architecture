package database

import (
	"errors"

	"gorm.io/gorm"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/domain/repository"
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

func (r *dbGroupRepository) List(f repository.GroupListFilter) (model.Groups, error) {
	var (
		dmgs  datamodel.Groups
		dmgus datamodel.GroupUsers
	)

	gdb := r.db
	gudb := r.db

	if len(f.UserIDs) > 0 {
		if err := gudb.Where("user_id IN (?)", f.UserIDs).Find(&dmgus).Error; err != nil {
			return nil, err
		}
		if len(dmgus) == 0 {
			return nil, nil

		}

		if err := gdb.Where("group_id IN (?)", dmgus.GroupIDs()).Find(&dmgs).Error; err != nil {
			return nil, err
		}
	} else {
		if err := gdb.Find(&dmgs).Error; err != nil {
			return nil, err
		}
		if len(dmgs) == 0 {
			return nil, nil
		}

		if err := gudb.Where("group_id IN (?)", dmgs.IDs()).Find(&dmgus).Error; err != nil {
			return nil, err
		}
	}

	return dmgs.ToModel(dmgus), nil
}

func (r *dbGroupRepository) ListByUserIDs(uIDs []model.UserID) (model.Groups, error) {
	if len(uIDs) == 0 {
		return nil, nil
	}

	gudb := r.db
	var dmgus datamodel.GroupUsers
	if err := gudb.Where("user_id IN (?)", uIDs).Find(&dmgus).Error; err != nil {
		return nil, err
	}
	if len(dmgus) == 0 {
		return nil, nil
	}

	gdb := r.db
	var dmgs datamodel.Groups
	if err := gdb.Where("group_id IN (?)", dmgus.GroupIDs()).Find(&dmgs).Error; err != nil {
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
	return r.db.Delete(&datamodel.Group{
		ID: string(gID),
	}).Error
}

func (r *dbGroupRepository) AddUsers(gID model.GroupID, uIDs []model.UserID) error {
	if gID == "" {
		return errors.New("group id must not be empty")
	}
	if len(uIDs) == 0 {
		return errors.New("user ids must not be empty")
	}
	dmgus := datamodel.NewGroupUsers(gID, uIDs)
	return r.db.Create(dmgus).Error
}

func (r *dbGroupRepository) RemoveUsers(gID model.GroupID, uIDs []model.UserID) error {
	if gID == "" {
		return errors.New("group id must not be empty")
	}
	if len(uIDs) == 0 {
		return errors.New("user ids must not be empty")
	}

	return r.db.
		Where("group_id = ?", gID).
		Where("user_id IN (?)", uIDs).
		Delete(&datamodel.Group{}).
		Error
}

func (r *dbGroupRepository) RemoveUsersFromAll(uIDs []model.UserID) error {
	if len(uIDs) == 0 {
		return errors.New("user ids must not be empty")
	}

	return r.db.
		Where("user_id IN (?)", uIDs).
		Delete(&datamodel.Group{}).
		Error
}
