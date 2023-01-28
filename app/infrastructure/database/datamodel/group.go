package datamodel

import "github.com/toshiykst/go-layerd-architecture/app/domain/model"

type Group struct {
	ID   string `gorm:"primaryKey"`
	Name string
}

func NewGroup(gID model.GroupID, name string) *Group {
	return &Group{
		ID:   string(gID),
		Name: name,
	}
}

func ToGroupModel(g *Group, gus GroupUsers) *model.Group {
	return &model.Group{
		ID:      model.GroupID(g.ID),
		Name:    g.Name,
		UserIDs: gus.UserIDs(),
	}
}
