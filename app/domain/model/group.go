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

func (gs Groups) IDs() []GroupID {
	if gs == nil {
		return nil
	}
	gIDs := make([]GroupID, len(gs))
	for i, g := range gs {
		gIDs[i] = g.ID()
	}
	return gIDs
}

func (gs Groups) UserIDs() []UserID {
	if gs == nil {
		return nil
	}
	m := make(map[UserID]struct{})
	var uIDs []UserID
	for _, g := range gs {
		for _, uID := range g.UserIDs() {
			if _, ok := m[uID]; !ok {
				m[uID] = struct{}{}
				uIDs = append(uIDs, uID)
			}
		}
	}
	return uIDs
}
