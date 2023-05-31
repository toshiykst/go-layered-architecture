package mockrepository

import (
	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

type mockGroupRepository struct {
	s *store
}

func (r *mockGroupRepository) Find(gID model.GroupID) (*model.Group, error) {
	for _, g := range r.s.groups {
		if g.ID() == gID {
			return g, nil
		}
	}
	return nil, nil
}

func (r *mockGroupRepository) List() ([]*model.Group, error) {
	return nil, nil
}

func (r *mockGroupRepository) Create(g *model.Group) (*model.Group, error) {
	return nil, nil
}

func (r *mockGroupRepository) Update(g *model.Group) error {
	return nil
}

func (r *mockGroupRepository) Delete(id model.GroupID) error {
	return nil
}
