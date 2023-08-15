package memory

import (
	"errors"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

type memoryGroupRepository struct {
	s *store
}

func (r *memoryGroupRepository) Find(gID model.GroupID) (*model.Group, error) {
	for _, g := range r.s.groups {
		if g.ID() == gID {
			return g, nil
		}
	}
	return nil, nil
}

func (r *memoryGroupRepository) List() (model.Groups, error) {
	return r.s.groups, nil
}

func (r *memoryGroupRepository) Create(g *model.Group) (*model.Group, error) {
	r.s.AddGroups(g)
	return g, nil
}

func (r *memoryGroupRepository) Update(g *model.Group) error {
	if g.ID() == "" {
		return errors.New("group id must not be empty")
	}

	for i, group := range r.s.groups {
		if group.ID() == g.ID() {
			r.s.groups[i] = g
			break
		}
	}

	return nil
}

func (r *memoryGroupRepository) Delete(gID model.GroupID) error {
	var result model.Groups
	for _, g := range r.s.groups {
		if g.ID() != gID {
			result = append(result, g)
		}
	}
	r.s.groups = result
	return nil
}

func (r *memoryGroupRepository) AddUsers(gID model.GroupID, uIDs []model.UserID) error {
	for i, g := range r.s.groups {
		if g.ID() == gID {
			r.s.groups[i] = model.NewGroup(gID, g.Name(), uIDs)
			return nil
		}
	}

	return nil
}

func (r *memoryGroupRepository) RemoveUsers(gID model.GroupID, uIDs []model.UserID) error {
	for i, g := range r.s.groups {
		if g.ID() == gID {
			var removed []model.UserID
			for _, guID := range r.s.groups[i].UserIDs() {
				for _, uID := range uIDs {
					if guID != uID {
						removed = append(removed, guID)
					}
				}
			}
			r.s.groups[i] = model.NewGroup(gID, g.Name(), removed)
			return nil
		}
	}
	return nil
}
