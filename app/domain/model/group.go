package model

// GroupID is identifier for group.
type GroupID string

// Group is a user domain model.

type Group struct {
	ID      string
	Name    string
	UserIDs []UserID
}
