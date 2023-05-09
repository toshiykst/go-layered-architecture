package dto

type (
	UpdateUserInput struct {
		UserID string
		Name   string
		Email  string
	}

	UpdateUserOutput struct {
		User User
	}
)
