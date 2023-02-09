package mockrepository

import (
	"github.com/toshiykst/go-layerd-architecture/app/domain/repository"
)

type mockRepository struct {
	s *store
}

type mockTransaction struct {
	s *store
}

func NewMockRepository(s *store) repository.Repository {
	return &mockRepository{s: s}
}

func (r *mockRepository) RunTransaction(f func(repository.Transaction) error) error {
	if err := f(&mockTransaction{s: r.s}); err != nil {
		return err
	}
	return nil
}

func (r *mockRepository) User() repository.UserRepositoryQuery {
	return &mockUserRepository{s: r.s}
}
func (tx *mockTransaction) User() repository.UserRepositoryCommand {
	return &mockUserRepository{s: tx.s}
}
func (r *mockRepository) Group() repository.GroupRepositoryQuery {
	return &mockGroupRepository{s: r.s}
}
func (tx *mockTransaction) Group() repository.GroupRepositoryCommand {
	return &mockGroupRepository{s: tx.s}
}
