package database

import (
	"gorm.io/gorm"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

type dbGroupRepository struct {
	db *gorm.DB
}

func (r *dbGroupRepository) Find(gID model.GroupID) (*model.Group, error) {
	return nil, nil
}

func (r *dbGroupRepository) List() ([]*model.Group, error) {
	return nil, nil
}

func (r *dbGroupRepository) Create(g *model.Group) (*model.Group, error) {
	return nil, nil
}

func (r *dbGroupRepository) Update(g *model.Group) error {
	return nil
}

func (r *dbGroupRepository) Delete(gID model.GroupID) error {
	return nil
}
