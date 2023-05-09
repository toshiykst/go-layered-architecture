package dto

type (
	CreateUserInput struct {
		Name  string
		Email string
	}

	CreateUserOutput struct {
		UserID string
		Name   string
		Email  string
	}
)
