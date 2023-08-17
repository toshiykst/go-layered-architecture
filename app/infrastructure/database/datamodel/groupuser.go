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

func (gus GroupUsers) GroupIDs() []string {
	if gus == nil {
		return nil
	}
	m := make(map[string]struct{})
	var gIDs []string
	for _, gu := range gus {
		gID := gu.GroupID
		if _, ok := m[gID]; !ok {
			m[gID] = struct{}{}
			gIDs = append(gIDs, gID)
		}
	}
	return gIDs
}

func (gus GroupUsers) ModelUserIDs() []model.UserID {
	uIDs := make([]model.UserID, len(gus))
	for i, v := range gus {
		uIDs[i] = model.UserID(v.UserID)
	}
	return uIDs
}

func (gus GroupUsers) ModelUserIDsByGroupID() map[string][]model.UserID {
	result := make(map[string][]model.UserID, len(gus))
	for _, gu := range gus {
		gID := gu.GroupID
		var uIDs []model.UserID
		if found, ok := result[gID]; ok {
			uIDs = append(found, model.UserID(gu.UserID))
		} else {
			uIDs = []model.UserID{model.UserID(gu.UserID)}
		}
		result[gID] = uIDs
	}
	return result
}
