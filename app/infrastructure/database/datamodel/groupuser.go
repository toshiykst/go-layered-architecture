package datamodel

import "github.com/toshiykst/go-layerd-architecture/app/domain/model"

type GroupUser struct {
	GroupID string `gorm:"primaryKey"`
	UserID  string `gorm:"primaryKey"`
}

func NewGroupUser(gID model.GroupID, uID model.UserID) *GroupUser {
	return &GroupUser{
		GroupID: string(gID),
		UserID:  string(uID),
	}
}

type GroupUsers []*GroupUser

func NewGroupUsers(gID model.GroupID, uIDs []model.UserID) GroupUsers {
	gus := make(GroupUsers, len(uIDs))
	for i, uID := range uIDs {
		gus[i] = NewGroupUser(gID, uID)
	}
	return gus
}

func (gus GroupUsers) UserIDs() []model.UserID {
	uIDs := make([]model.UserID, len(gus))
	for i, v := range gus {
		uIDs[i] = model.UserID(v.UserID)
	}
	return uIDs
}
