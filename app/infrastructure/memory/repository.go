package memory

import (
	"github.com/toshiykst/go-layerd-architecture/app/domain/repository"
)

type memoryRepository struct {
	s *store
}

type memoryTransaction struct {
	s *store
}

func NewMemoryRepository(s *store) repository.Repository {
	return &memoryRepository{s: s}
}

func (r *memoryRepository) RunTransaction(f func(repository.Transaction) error) error {
	if err := f(&memoryTransaction{s: r.s}); err != nil {
		return err
	}
	return nil
}

func (r *memoryRepository) User() repository.UserRepositoryQuery {
	return &memoryUserRepository{s: r.s}
}
func (tx *memoryTransaction) User() repository.UserRepositoryCommand {
	return &memoryUserRepository{s: tx.s}
}
func (r *memoryRepository) Group() repository.GroupRepositoryQuery {
	return &memoryGroupRepository{s: r.s}
}
func (tx *memoryTransaction) Group() repository.GroupRepositoryCommand {
	return &memoryGroupRepository{s: tx.s}
}
