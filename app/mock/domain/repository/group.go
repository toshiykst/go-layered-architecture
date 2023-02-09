package mockrepository

import (
	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

type mockGroupRepository struct {
	s *store
}

func (db *mockGroupRepository) Find(gID model.GroupID) (*model.Group, error) {
	return nil, nil
}

func (db *mockGroupRepository) FindByName(name string) (*model.Group, error) {
	return nil, nil
}

func (db *mockGroupRepository) List() ([]*model.Group, error) {
	return nil, nil
}

func (db *mockGroupRepository) Create(g *model.Group) (*model.Group, error) {
	return nil, nil
}

func (db *mockGroupRepository) Update(g *model.Group) error {
	return nil
}

func (db *mockGroupRepository) Delete(id model.GroupID) error {
	return nil
}
