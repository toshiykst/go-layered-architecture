package database

import (
	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"gorm.io/gorm"
)

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
