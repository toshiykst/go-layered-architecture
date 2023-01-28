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
	return model.NewGroup(
		model.GroupID(g.ID),
		g.Name,
		gus.UserIDs(),
	)
}
