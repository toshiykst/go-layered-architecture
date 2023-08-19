package model

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidGroup = errors.New("invalid group id")
)

const (
	maxGroupNameLength = 30
	maxGroupUserCount  = 5
)

type GroupID string

type Group struct {
	id      GroupID
	name    string
	userIDs []UserID
}

func NewGroup(id GroupID, name string, uIDs []UserID) (*Group, error) {
	if id == "" {
		return nil, fmt.Errorf("group id must not empty: %w", ErrInvalidGroup)
	}

	if name == "" {
		return nil, fmt.Errorf("group name must not empty: %w", ErrInvalidGroup)
	}
	if len(name) > maxGroupNameLength {
		return nil, fmt.Errorf("exceeds the max group name length: %w", ErrInvalidGroup)
	}

	if len(uIDs) > maxGroupUserCount {
		return nil, fmt.Errorf("exceeds the max group users: %w", ErrInvalidGroup)
	}

	return &Group{
		id:      id,
		name:    name,
		userIDs: uIDs,
	}, nil
}

func MustNewGroup(gID GroupID, name string, uIDs []UserID) *Group {
	g, err := NewGroup(gID, name, uIDs)
	if err != nil {
		panic(err)
	}
	return g
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

func (g *Group) IsMaxUsers() bool {
	if g == nil {
		return false
	}
	if len(g.userIDs) < maxGroupUserCount {
		return false
	}
	return true
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
