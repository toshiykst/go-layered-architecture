package model

// UserID is identifier for user.
type UserID string

// User is a user domain model.
type User struct {
	ID    UserID
	Name  string
	Email string
}
