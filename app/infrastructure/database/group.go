package database

import (
	"gorm.io/gorm"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

type dbGroupRepository struct {
	db *gorm.DB
}

func (db *dbGroupRepository) Find(gID model.GroupID) (*model.Group, error) {
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
