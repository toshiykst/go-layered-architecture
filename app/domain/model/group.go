package model

type GroupID string

type Group struct {
	id      GroupID
	name    string
	userIDs []UserID
}

func NewGroup(id GroupID, name string, uIDs []UserID) *Group {
	return &Group{
		id:      id,
		name:    name,
		userIDs: uIDs,
	}
}

func (g *Group) ID() GroupID {
	if g == nil {
		return ""
	}
	return g.id
}

func (g *Group) Name() string {
	if g == nil {
		return ""
	}
	return g.name
}

func (g *Group) UserIDs() []UserID {
	if g == nil {
		return nil
	}
	return g.userIDs
}

type Groups []*Group
