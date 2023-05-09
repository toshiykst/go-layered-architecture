package dto

type (
	GetUserInput struct {
		UserID string
	}

	GetUserOutput struct {
		UserID string
		Name   string
		Email  string
	}
)
