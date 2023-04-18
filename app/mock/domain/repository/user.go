package mockrepository

import (
	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

type mockUserRepository struct {
	s *store
}

func (r *mockUserRepository) Find(uID model.UserID) (*model.User, error) {
	for _, u := range r.s.users {
		if u.ID() == uID {
			return u, nil
		}
	}
	return nil, nil
}

func (r *mockUserRepository) FindByName(name string) (*model.User, error) {
	return nil, nil
}

func (r *mockUserRepository) List() (model.Users, error) {
	return r.s.users, nil
}

func (r *mockUserRepository) Create(u *model.User) (*model.User, error) {
	r.s.AddUsers(u)
	return u, nil
}

func (r *mockUserRepository) Update(u *model.User) error {
	return nil
}

func (r *mockUserRepository) Delete(id model.UserID) error {
	return nil
}
