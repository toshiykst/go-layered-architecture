package memory

import (
	"errors"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/domain/repository"
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

func (r *memoryGroupRepository) List(f repository.GroupListFilter) (model.Groups, error) {
	var result model.Groups
	for _, g := range r.s.groups {
		if len(f.UserIDs) > 0 {
			found := false
			for _, uID := range f.UserIDs {
				for _, guID := range g.UserIDs() {
					if uID == guID {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				continue
			}
		}

		result = append(result, g)
	}

	return result, nil
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

func (r *memoryGroupRepository) Delete(g *model.Group) error {
	gID := g.ID()
	var result model.Groups
	for _, sg := range r.s.groups {
		if sg.ID() != gID {
			result = append(result, g)
		}
	}
	r.s.groups = result
	return nil
}

func (r *memoryGroupRepository) AddUsers(gID model.GroupID, uIDs []model.UserID) error {
	for i, g := range r.s.groups {
		if g.ID() == gID {
			mg, err := model.NewGroup(gID, g.Name(), append(g.UserIDs(), uIDs...))
			if err != nil {
				return err
			}
			r.s.groups[i] = mg
			return nil
		}
	}

	return nil
}

func (r *memoryGroupRepository) RemoveUsers(gID model.GroupID, uIDs []model.UserID) error {
	for i, g := range r.s.groups {
		if g.ID() != gID {
			continue
		}

		var removed []model.UserID
		for _, guID := range r.s.groups[i].UserIDs() {
			found := false
			for _, uID := range uIDs {
				if guID == uID {
					found = true
					break
				}
			}
			if !found {
				removed = append(removed, guID)
			}
		}

		mg, err := model.NewGroup(g.ID(), g.Name(), removed)
		if err != nil {
			return err
		}

		r.s.groups[i] = mg
		return nil

	}
	return nil
}

func (r *memoryGroupRepository) RemoveUsersFromAll(uIDs []model.UserID) error {
	for i, g := range r.s.groups {
		var removed []model.UserID
		for _, guID := range r.s.groups[i].UserIDs() {
			found := false
			for _, uID := range uIDs {
				if guID == uID {
					found = true
					break
				}
			}
			if !found {
				removed = append(removed, guID)
			}
		}
		mg, err := model.NewGroup(g.ID(), g.Name(), removed)
		if err != nil {
			return err
		}

		r.s.groups[i] = mg
		return nil
	}
	return nil
}
