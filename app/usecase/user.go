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

type UserUsecase interface {
	CreateUser(in *dto.CreateUserInput) (*dto.CreateUserOutput, error)
	GetUser(in *dto.GetUserInput) (*dto.GetUserOutput, error)
	GetUsers(in *dto.GetUsersInput) (*dto.GetUsersOutput, error)
	UpdateUser(in *dto.UpdateUserInput) (*dto.UpdateUserOutput, error)
	DeleteUser(in *dto.DeleteUserInput) (*dto.DeleteUserOutput, error)
}

type userUsecase struct {
	r  repository.Repository
	f  factory.UserFactory
	us domainservice.UserService
	gs domainservice.GroupService
}

func NewUserUsecase(
	r repository.Repository,
	f factory.UserFactory,
	us domainservice.UserService,
	gs domainservice.GroupService,
) UserUsecase {
	return &userUsecase{r: r, f: f, us: us, gs: gs}
}

func (uc *userUsecase) CreateUser(in *dto.CreateUserInput) (*dto.CreateUserOutput, error) {
	u, err := uc.f.Create(in.Name, in.Email)
	if err != nil {
		if errors.Is(err, model.ErrInvalidUser) {
			return nil, errors.Join(ErrInvalidUserInput, err)
		}
		return nil, err
	}

	if err = uc.r.RunTransaction(func(tx repository.Transaction) error {
		if _, err := tx.User().Create(u); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &dto.CreateUserOutput{
		User: dto.User{
			UserID: string(u.ID()),
			Name:   u.Name(),
			Email:  u.Email(),
		},
	}, nil
}

func (uc *userUsecase) GetUser(in *dto.GetUserInput) (*dto.GetUserOutput, error) {
	uID := model.UserID(in.UserID)

	u, err := uc.r.User().Find(uID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		// TODO: Use custom logger(zap)
		log.Warnf("the user is not found; userID=%s", uID)
		return nil, ErrUserNotFound
	}

	return &dto.GetUserOutput{
		User: dto.User{
			UserID: string(u.ID()),
			Name:   u.Name(),
			Email:  u.Email(),
		},
	}, nil
}

func (uc *userUsecase) GetUsers(_ *dto.GetUsersInput) (*dto.GetUsersOutput, error) {
	us, err := uc.r.User().List(repository.UserListFilter{})
	if err != nil {
		return nil, err
	}

	return &dto.GetUsersOutput{
		Users: dto.ToUsersFromModel(us),
	}, nil
}

func (uc *userUsecase) UpdateUser(in *dto.UpdateUserInput) (*dto.UpdateUserOutput, error) {
	u, err := model.NewUser(model.UserID(in.UserID), in.Name, in.Email)
	if err != nil {
		return nil, errors.Join(ErrInvalidUserInput, err)
	}

	isExisted, err := uc.us.Exists(u.ID())
	if err != nil {
		return nil, err
	}
	if !isExisted {
		return nil, ErrUserNotFound
	}

	if err := uc.r.RunTransaction(func(tx repository.Transaction) error {
		if err := tx.User().Update(u); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &dto.UpdateUserOutput{}, nil
}

func (uc *userUsecase) DeleteUser(in *dto.DeleteUserInput) (*dto.DeleteUserOutput, error) {
	uID := model.UserID(in.UserID)

	isExisted, err := uc.us.Exists(uID)
	if err != nil {
		return nil, err
	}
	if !isExisted {
		return nil, ErrUserNotFound
	}

	hasGroupUser, err := uc.gs.HasUsersAny([]model.UserID{uID})
	if err != nil {
		return nil, err
	}

	if err := uc.r.RunTransaction(func(tx repository.Transaction) error {
		if hasGroupUser {
			if err := tx.Group().RemoveUsersFromAll([]model.UserID{uID}); err != nil {
				return err
			}
		}

		if err := tx.User().Delete(uID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &dto.DeleteUserOutput{}, nil
}
