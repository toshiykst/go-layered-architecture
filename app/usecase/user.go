package usecase

import (
	"github.com/toshiykst/go-layerd-architecture/app/domain/factory"
	"github.com/toshiykst/go-layerd-architecture/app/domain/repository"
)

type UserUsecase interface {
	CreateUser(in *CreateUserInput) (*CreateUserOutput, error)
	GetUser(in *GetUserInput) (*GetUserOutput, error)
	UpdateUser(in *UpdateUserInput) (*UpdateUserOutput, error)
	DeleteUser(in *DeleteUserInput) (*DeleteUserOutput, error)
}

type userUsecase struct {
	r repository.Repository
	f factory.UserFactory
}

func NewUserUsecase(r repository.Repository, f factory.UserFactory) UserUsecase {
	return &userUsecase{r: r, f: f}
}

type CreateUserInput struct {
	Name  string
	Email string
}

type CreateUserOutput struct {
	UserID string
	Name   string
	Email  string
}

func (uc *userUsecase) CreateUser(in *CreateUserInput) (*CreateUserOutput, error) {
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

	return &CreateUserOutput{
		UserID: string(u.ID()),
		Name:   u.Name(),
		Email:  u.Email(),
	}, nil
}

type GetUserInput struct {
	UserID string
}

type GetUserOutput struct {
	UserID string
	Name   string
	Email  string
}

func (uc *userUsecase) GetUser(in *GetUserInput) (*GetUserOutput, error) {
	return &GetUserOutput{}, nil
}

type UpdateUserInput struct {
	UserID string
	Name   string
	Email  string
}

type UpdateUserOutput struct {
	UserID string
	Name   string
	Email  string
}

func (uc *userUsecase) UpdateUser(in *UpdateUserInput) (*UpdateUserOutput, error) {
	return &UpdateUserOutput{}, nil
}

type DeleteUserInput struct {
	UserID string
}

type DeleteUserOutput struct {
}

func (uc *userUsecase) DeleteUser(in *DeleteUserInput) (*DeleteUserOutput, error) {
	return &DeleteUserOutput{}, nil
}
