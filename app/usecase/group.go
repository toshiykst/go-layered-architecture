//go:generate mockgen -source=$GOFILE -package=mock$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"errors"

	"github.com/labstack/gommon/log"

	"github.com/toshiykst/go-layerd-architecture/app/domain/domainservice"
	"github.com/toshiykst/go-layerd-architecture/app/domain/factory"
	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
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
	r  repository.Repository
	f  factory.GroupFactory
	gs domainservice.GroupService
	us domainservice.UserService
}

func NewGroupUsecase(
	r repository.Repository,
	f factory.GroupFactory,
	gs domainservice.GroupService,
	us domainservice.UserService,
) GroupUsecase {
	return &groupUsecase{r: r, f: f, gs: gs, us: us}
}

func (uc *groupUsecase) CreateGroup(in *dto.CreateGroupInput) (*dto.CreateGroupOutput, error) {
	g, err := uc.f.Create(in.Name, dto.ToModelUserIDs(in.UserIDs))
	if err != nil {
		if errors.Is(err, model.ErrInvalidGroup) {
			return nil, errors.Join(ErrInvalidGroupInput, err)
		}
		return nil, err
	}

	uIDs := g.UserIDs()
	if len(uIDs) > 0 {
		ok, err := uc.us.ExistsAll(uIDs)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrInvalidUserIDs
		}
	}

	var created *model.Group
	if err = uc.r.RunTransaction(func(tx repository.Transaction) error {
		if created, err = tx.Group().Create(g); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	var us model.Users
	if len(uIDs) > 0 {
		us, err = uc.r.User().List(repository.UserListFilter{
			UserIDs: uIDs,
		})
		if err != nil {
			return nil, err
		}
	}

	return &dto.CreateGroupOutput{
		Group: dto.Group{
			GroupID: string(created.ID()),
			Name:    created.Name(),
			Users:   dto.ToUsersFromModel(us),
		},
	}, nil
}

func (uc *groupUsecase) GetGroup(in *dto.GetGroupInput) (*dto.GetGroupOutput, error) {
	gID := model.GroupID(in.GroupID)

	g, err := uc.r.Group().Find(gID)
	if err != nil {
		return nil, err
	}
	if g == nil {
		// TODO: Use custom logger(zap)
		log.Warnf("the group is not found; groupID=%s", gID)
		return nil, ErrGroupNotFound
	}

	us, err := uc.r.User().List(repository.UserListFilter{
		UserIDs: g.UserIDs(),
	})
	if err != nil {
		return nil, err
	}

	return &dto.GetGroupOutput{
		Group: dto.Group{
			GroupID: string(g.ID()),
			Name:    g.Name(),
			Users:   dto.ToUsersFromModel(us),
		},
	}, nil
}

func (uc *groupUsecase) GetGroups(_ *dto.GetGroupsInput) (*dto.GetGroupsOutput, error) {
	gs, err := uc.r.Group().List(repository.GroupListFilter{})
	if err != nil {
		return nil, err
	}
	if len(gs) == 0 {
		return &dto.GetGroupsOutput{
			Groups: []dto.Group{},
		}, nil
	}

	uIDs := gs.UserIDs()
	if len(uIDs) == 0 {
		dtogs := make([]dto.Group, len(gs))
		for i, g := range gs {
			dtogs[i] = dto.Group{
				GroupID: string(g.ID()),
				Name:    g.Name(),
				Users:   []dto.User{},
			}
		}

		return &dto.GetGroupsOutput{
			Groups: dtogs,
		}, nil
	}

	us, err := uc.r.User().List(repository.UserListFilter{
		UserIDs: uIDs,
	})
	if err != nil {
		return nil, err
	}

	usByUID := us.ByUserID()
	dtogs := make([]dto.Group, len(gs))
	for i, g := range gs {
		guIDs := g.UserIDs()
		gus := make(model.Users, len(guIDs))
		for j, guID := range guIDs {
			if gu, ok := usByUID[guID]; ok {
				gus[j] = gu
			}
		}
		dtogs[i] = dto.Group{
			GroupID: string(g.ID()),
			Name:    g.Name(),
			Users:   dto.ToUsersFromModel(gus),
		}
	}

	return &dto.GetGroupsOutput{
		Groups: dtogs,
	}, nil
}

func (uc *groupUsecase) UpdateGroup(in *dto.UpdateGroupInput) (*dto.UpdateGroupOutput, error) {
	g, err := model.NewGroup(model.GroupID(in.GroupID), in.Name, []model.UserID{})
	if err != nil {
		return nil, errors.Join(ErrInvalidGroupInput, err)
	}

	isExisted, err := uc.gs.Exists(g.ID())
	if err != nil {
		return nil, err
	}
	if !isExisted {
		return nil, ErrGroupNotFound
	}

	if err := uc.r.RunTransaction(func(tx repository.Transaction) error {
		if err := tx.Group().Update(g); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

func (uc *groupUsecase) DeleteGroup(in *dto.DeleteGroupInput) (*dto.DeleteGroupOutput, error) {
	gID := model.GroupID(in.GroupID)

	g, err := uc.r.Group().Find(gID)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, ErrGroupNotFound
	}

	if err := uc.r.RunTransaction(func(tx repository.Transaction) error {
		if len(g.UserIDs()) > 0 {
			if err = tx.Group().RemoveUsers(g.ID(), g.UserIDs()); err != nil {
				return err
			}
		}
		if err := tx.Group().Delete(gID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &dto.DeleteGroupOutput{}, nil
}
