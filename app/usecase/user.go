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
	r repository.Repository
	f factory.UserFactory
	s domainservice.UserService
}

func NewUserUsecase(
	r repository.Repository,
	f factory.UserFactory,
	s domainservice.UserService,
) UserUsecase {
	return &userUsecase{r: r, f: f, s: s}
}

var (
	ErrUserNotFound = errors.New("user not found")
)

func (uc *userUsecase) CreateUser(in *dto.CreateUserInput) (*dto.CreateUserOutput, error) {
	u, err := uc.f.Create(in.Name, in.Email)
	if err != nil {
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
	us, err := uc.r.User().List()
	if err != nil {
		return nil, err
	}

	return &dto.GetUsersOutput{
		Users: convertUsersToDTO(us),
	}, nil
}

func convertUsersToDTO(us model.Users) []dto.User {
	result := make([]dto.User, len(us))
	for i, u := range us {
		result[i] = dto.User{
			UserID: string(u.ID()),
			Name:   u.Name(),
			Email:  u.Email(),
		}
	}
	return result
}

func (uc *userUsecase) UpdateUser(in *dto.UpdateUserInput) (*dto.UpdateUserOutput, error) {
	u := model.NewUser(model.UserID(in.UserID), in.Name, in.Email)

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
	return &dto.DeleteUserOutput{}, nil
}
