package dto

type (
	CreateUserInput struct {
		Name  string
		Email string
	}

	CreateUserOutput struct {
		User User
	}
)
