package domainservice

import (
	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/domain/repository"
)

type GroupService interface {
	Exists(gID model.GroupID) (bool, error)
}

type groupService struct {
	r repository.Repository
}

func NewGroupService(r repository.Repository) GroupService {
	return &groupService{r: r}
}

func (gs *groupService) Exists(gID model.GroupID) (bool, error) {
	u, err := gs.r.Group().Find(gID)
	if err != nil {
		return false, err
	}
	if u == nil {
		return false, nil
	}
	return true, nil
}
