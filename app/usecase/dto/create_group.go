package dto

type (
	CreateGroupInput struct {
		Name    string
		UserIDs []string
	}

	CreateGroupOutput struct {
		Group Group
	}
)
