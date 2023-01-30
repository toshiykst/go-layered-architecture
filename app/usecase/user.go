package usecase

import "github.com/toshiykst/go-layerd-architecture/app/domain/repository"

type CreateUserInput struct {
	Name  string
	Email string
}

type CreateUserOutput struct {
	Name  string
	Email string
}

type GetUserInput struct {
	UserID string
}

type GetUserOutput struct {
	UserID string
	Name   string
	Email  string
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

type DeleteUserInput struct {
	UserID string
}

type DeleteUserOutput struct {
}

type UserUsecase interface {
	CreateUser(in *CreateUserInput) (*CreateUserOutput, error)
	GetUser(in *GetUserInput) (*GetUserOutput, error)
	UpdateUser(in *UpdateUserInput) (*UpdateUserOutput, error)
	DeleteUser(in *DeleteUserInput) (*DeleteUserOutput, error)
}

type userUsecase struct {
	r repository.Repository
}

func NewUserUsecase(r repository.Repository) UserUsecase {
	return &userUsecase{r: r}
}

func (uc *userUsecase) CreateUser(in *CreateUserInput) (*CreateUserOutput, error) {
	return &CreateUserOutput{}, nil
}

func (uc *userUsecase) GetUser(in *GetUserInput) (*GetUserOutput, error) {
	return &GetUserOutput{}, nil
}

func (uc *userUsecase) UpdateUser(in *UpdateUserInput) (*UpdateUserOutput, error) {
	return &UpdateUserOutput{}, nil
}

func (uc *userUsecase) DeleteUser(in *DeleteUserInput) (*DeleteUserOutput, error) {
	return &DeleteUserOutput{}, nil
}
