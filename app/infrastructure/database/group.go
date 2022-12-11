package database

import (
	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"gorm.io/gorm"
)

type groupDataModel struct {
	ID   string `gorm:"primaryKey"`
	Name string
}

type groupUserDataModel struct {
	UserID  string `gorm:"primaryKey"`
	GroupID string `gorm:"primaryKey"`
}

type groupUserDataModels []*groupUserDataModel

func (gudms groupUserDataModels) UserIDs() []model.UserID {
	uIDs := make([]model.UserID, len(gudms))
	for _, v := range gudms {
		uIDs = append(uIDs, model.UserID(v.UserID))
	}
	return uIDs
}

func toGroupModel(gdm *groupDataModel, gudms groupUserDataModels) *model.Group {
	return &model.Group{
		ID:      model.GroupID(gdm.ID),
		Name:    gdm.Name,
		UserIDs: gudms.UserIDs(),
	}
}

func toGroupDataModel(g *model.Group) *groupDataModel {
	return &groupDataModel{
		ID:   string(g.ID),
		Name: g.Name,
	}
}

func toGroupUserDataModels(g *model.Group) groupUserDataModels {
	gudms := make(groupUserDataModels, len(g.UserIDs))
	for _, uID := range g.UserIDs {
		gudms = append(gudms, &groupUserDataModel{
			GroupID: string(g.ID),
			UserID:  string(uID),
		})
	}
	return gudms
}

type dbGroupRepository struct {
	conn *gorm.DB
}

func (db *dbGroupRepository) Find(id model.GroupID) (*model.Group, error) {
	return nil, nil
}

func (db *dbGroupRepository) FindByName(name string) (*model.Group, error) {
	return nil, nil
}

func (db *dbGroupRepository) List() ([]*model.Group, error) {
	return nil, nil
}

func (db *dbGroupRepository) Create(g *model.Group) (*model.Group, error) {
	return nil, nil
}

func (db *dbGroupRepository) Update(g *model.Group) error {
	return nil
}

func (db *dbGroupRepository) Delete(id model.GroupID) error {
	return nil
}
