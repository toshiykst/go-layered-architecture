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

func (g *Group) ToModel(gus GroupUsers) *model.Group {
	if g == nil {
		return nil
	}
	return model.NewGroup(
		model.GroupID(g.ID),
		g.Name,
		gus.ModelUserIDs(),
	)
}

type Groups []*Group

func (gs Groups) IDs() []string {
	if gs == nil {
		return nil
	}
	gIDs := make([]string, len(gs))
	for i, g := range gs {
		gIDs[i] = g.ID
	}
	return gIDs
}

func (gs Groups) ToModel(gus GroupUsers) model.Groups {
	if gs == nil {
		return nil
	}
	uIDsByGID := gus.ModelUserIDsByGroupID()
	mgs := make(model.Groups, len(gs))
	for i, g := range gs {
		mgs[i] = model.NewGroup(
			model.GroupID(g.ID),
			g.Name,
			uIDsByGID[g.ID],
		)
	}
	return mgs
}
