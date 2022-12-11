package usecase

import "github.com/toshiykst/go-layerd-architecture/app/domain/repository"

// CreateUserInput is input struct of CreateUser.
type CreateUserInput struct {
	Name  string
	Email string
}

// CreateUserOutput is output struct of CreateUser.
type CreateUserOutput struct {
	Name  string
	Email string
}

// GetUserInput is input struct of GetUser.
type GetUserInput struct {
	UserID string
}

// GetUserOutput is input struct of GetUser.
type GetUserOutput struct {
	UserID string
	Name   string
	Email  string
}

// UpdateUserInput is input struct of UpdateUser.
type UpdateUserInput struct {
	UserID string
	Name   string
	Email  string
}

// UpdateUserOutput is input struct of UpdateUser.
type UpdateUserOutput struct {
	UserID string
	Name   string
	Email  string
}

// DeleteUserInput is input struct of DeleteUser.
type DeleteUserInput struct {
	UserID string
}

// DeleteUserOutput is input struct of DeleteUser.
type DeleteUserOutput struct {
}

// UserUsecase is an usecase of user.
type UserUsecase interface {
	CreateUser(in *CreateUserInput) (*CreateUserOutput, error)
	GetUser(in *GetUserInput) (*GetUserOutput, error)
	UpdateUser(in *UpdateUserInput) (*UpdateUserOutput, error)
	DeleteUser(in *DeleteUserInput) (*DeleteUserOutput, error)
}

type userUsecaseImpl struct {
	r repository.Repository
}

// NewUserUsecase returns a user usecase instance.
func NewUserUsecase(r repository.Repository) UserUsecase {
	return &userUsecaseImpl{r: r}
}

// CreateUser creates a user.
func (uc *userUsecaseImpl) CreateUser(in *CreateUserInput) (*CreateUserOutput, error) {
	return &CreateUserOutput{}, nil
}

// GetUser returns a user.
func (uc *userUsecaseImpl) GetUser(in *GetUserInput) (*GetUserOutput, error) {
	return &GetUserOutput{}, nil
}

// UpdateUser updates a user.
func (uc *userUsecaseImpl) UpdateUser(in *UpdateUserInput) (*UpdateUserOutput, error) {
	return &UpdateUserOutput{}, nil
}

// DeleteUser deletes a user.
func (uc *userUsecaseImpl) DeleteUser(in *DeleteUserInput) (*DeleteUserOutput, error) {
	return &DeleteUserOutput{}, nil
}
