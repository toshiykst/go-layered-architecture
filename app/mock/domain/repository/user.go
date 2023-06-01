package mockrepository

import (
	"errors"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/domain/repository"
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

func (r *mockUserRepository) List(f repository.UserListFilter) (model.Users, error) {
	var result model.Users
	for _, u := range r.s.users {
		if len(f.UserIDs) > 0 {
			found := false
			for _, fUID := range f.UserIDs {
				if u.ID() == fUID {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		result = append(result, u)
	}

	return result, nil
}

func (r *mockUserRepository) Create(u *model.User) (*model.User, error) {
	r.s.AddUsers(u)
	return u, nil
}

func (r *mockUserRepository) Update(u *model.User) error {
	if u.ID() == "" {
		return errors.New("user id must not be empty")
	}

	for i, user := range r.s.users {
		if user.ID() == u.ID() {
			r.s.users[i] = u
			break
		}
	}

	return nil
}

func (r *mockUserRepository) Delete(uID model.UserID) error {
	var result model.Users
	for _, user := range r.s.users {
		if user.ID() != uID {
			result = append(result, user)
		}
	}
	r.s.users = result
	return nil
}
