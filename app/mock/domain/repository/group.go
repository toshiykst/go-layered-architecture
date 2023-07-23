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
	r.s.AddGroups(g)
	return g, nil
}

func (r *mockGroupRepository) Update(g *model.Group) error {
	return nil
}

func (r *mockGroupRepository) Delete(id model.GroupID) error {
	return nil
}

func (r *mockGroupRepository) AddUsers(gID model.GroupID, uIDs []model.UserID) error {
	for i, g := range r.s.groups {
		if g.ID() == gID {
			r.s.groups[i] = model.NewGroup(gID, g.Name(), uIDs)
			return nil
		}
	}

	return nil
}

func (r *mockGroupRepository) RemoveUsers(gID model.GroupID, uIDs []model.UserID) error {
	return nil
}
