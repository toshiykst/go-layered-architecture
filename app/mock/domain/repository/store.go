package mockrepository

import "github.com/toshiykst/go-layerd-architecture/app/domain/model"

type store struct {
	users  model.Users
	groups model.Groups
}

func NewStore() *store {
	return &store{}
}

func (s *store) AddUsers(us ...*model.User) {
	s.users = append(s.users, us...)
}

func (s *store) AddGroups(gs ...*model.Group) {
	s.groups = append(s.groups, gs...)
}
