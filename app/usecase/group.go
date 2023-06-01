//go:generate mockgen -source=$GOFILE -package=mock$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"errors"

	"github.com/toshiykst/go-layerd-architecture/app/domain/repository"
	"github.com/toshiykst/go-layerd-architecture/app/usecase/dto"
)

type GroupUsecase interface {
	CreateGroup(in *dto.CreateGroupInput) (*dto.CreateGroupOutput, error)
	GetGroup(in *dto.GetGroupInput) (*dto.GetGroupOutput, error)
	GetGroups(in *dto.GetGroupsInput) (*dto.GetGroupsOutput, error)
	UpdateGroup(in *dto.UpdateGroupInput) (*dto.UpdateGroupOutput, error)
	DeleteGroup(in *dto.DeleteGroupInput) (*dto.DeleteGroupOutput, error)
}

type groupUsecase struct {
	r repository.Repository
}

func NewGroupUsecase(
	r repository.Repository,
) GroupUsecase {
	return &groupUsecase{r: r}
}

var (
	ErrGroupNotFound = errors.New("group not found")
)

func (uc *groupUsecase) CreateGroup(in *dto.CreateGroupInput) (*dto.CreateGroupOutput, error) {
	return nil, nil
}

func (uc *groupUsecase) GetGroup(in *dto.GetGroupInput) (*dto.GetGroupOutput, error) {
	return nil, nil
}

func (uc *groupUsecase) GetGroups(_ *dto.GetGroupsInput) (*dto.GetGroupsOutput, error) {
	return nil, nil
}

func (uc *groupUsecase) UpdateGroup(in *dto.UpdateGroupInput) (*dto.UpdateGroupOutput, error) {
	return nil, nil
}

func (uc *groupUsecase) DeleteGroup(in *dto.DeleteGroupInput) (*dto.DeleteGroupOutput, error) {
	return nil, nil
}
