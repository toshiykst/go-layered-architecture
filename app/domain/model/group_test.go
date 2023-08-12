package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGroup_ID(t *testing.T) {
	tests := []struct {
		name  string
		group *Group
		want  GroupID
	}{
		{
			name:  "Returns group’s id",
			group: &Group{id: GroupID("TEST_GROUP_ID")},
			want:  GroupID("TEST_GROUP_ID"),
		},
		{
			name:  "Receiver is nil",
			group: nil,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.group.ID(); got != tt.want {
				t.Errorf("ID()=%v; want %v, receiver=%v", got, tt.want, tt.group)
			}
		})
	}
}

func TestGroup_Name(t *testing.T) {
	tests := []struct {
		name  string
		group *Group
		want  string
	}{
		{
			name:  "Returns group’s name",
			group: &Group{name: "TEST_NAME"},
			want:  "TEST_NAME",
		},
		{
			name:  "Receiver is nil",
			group: nil,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.group.Name(); got != tt.want {
				t.Errorf("Name()=%v; want %v, receiver=%v", got, tt.want, tt.group)
			}
		})
	}
}

func TestGroup_UserIDs(t *testing.T) {
	tests := []struct {
		name  string
		group *Group
		want  []UserID
	}{
		{
			name:  "Returns group’s user ids",
			group: &Group{userIDs: []UserID{"TEST_USER_ID_1", "TEST_USER_ID_2", "TEST_USER_ID_3"}},
			want:  []UserID{"TEST_USER_ID_1", "TEST_USER_ID_2", "TEST_USER_ID_3"},
		},
		{
			name:  "Receiver is nil",
			group: nil,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.group.UserIDs()
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("UserIDs()=%v; want %v, receiver=%v\ndiffers: (-got +want)\n%s", got, tt.want, tt.group, diff)
			}
		})
	}
}

func TestGroups_IDs(t *testing.T) {
	tests := []struct {
		name   string
		groups Groups
		want   []GroupID
	}{
		{
			name: "Returns group ids",
			groups: Groups{
				NewGroup(
					"TEST_GROUP_ID_1",
					"TEST_GROUP_NAME_1",
					[]UserID{"TEST_USER_ID_1"},
				),
				NewGroup(
					"TEST_GROUP_ID_2",
					"TEST_GROUP_NAME_2",
					[]UserID{"TEST_USER_ID_2"},
				),
				NewGroup(
					"TEST_GROUP_ID_3",
					"TEST_GROUP_NAME_3",
					[]UserID{"TEST_USER_ID_3"},
				),
			},
			want: []GroupID{
				"TEST_GROUP_ID_1",
				"TEST_GROUP_ID_2",
				"TEST_GROUP_ID_3",
			},
		},
		{
			name:   "Returns nil when the receiver is nil",
			groups: nil,
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := tt.groups
			got := gs.IDs()
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf(
					"gs.IDs()=%v; want=%v,receiver=%v\ndiffers: (-got +want)\n%s",
					got, tt.want, gs, diff,
				)
			}
		})
	}
}

func TestGroups_UserIDs(t *testing.T) {
	tests := []struct {
		name   string
		groups Groups
		want   []UserID
	}{
		{
			name: "Returns user ids",
			groups: Groups{
				NewGroup(
					"TEST_GROUP_ID_1",
					"TEST_GROUP_NAME_1",
					[]UserID{"TEST_USER_ID_1"},
				),
				NewGroup(
					"TEST_GROUP_ID_2",
					"TEST_GROUP_NAME_2",
					[]UserID{"TEST_USER_ID_1", "TEST_USER_ID_2"},
				),
				NewGroup(
					"TEST_GROUP_ID_3",
					"TEST_GROUP_NAME_3",
					[]UserID{"TEST_USER_ID_1", "TEST_USER_ID_2", "TEST_USER_ID_3"},
				),
			},
			want: []UserID{
				"TEST_USER_ID_1",
				"TEST_USER_ID_2",
				"TEST_USER_ID_3",
			},
		},
		{
			name: "Returns nil when all groups have no user ids",
			groups: Groups{
				NewGroup(
					"TEST_GROUP_ID_1",
					"TEST_GROUP_NAME_1",
					[]UserID{},
				),
				NewGroup(
					"TEST_GROUP_ID_2",
					"TEST_GROUP_NAME_2",
					[]UserID{},
				),
				NewGroup(
					"TEST_GROUP_ID_3",
					"TEST_GROUP_NAME_3",
					[]UserID{},
				),
			},
			want: nil,
		},
		{
			name:   "Returns nil when the receiver is nil",
			groups: nil,
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := tt.groups
			got := gs.UserIDs()
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf(
					"gs.UserIDs()=%v; want=%v,receiver=%v\ndiffers: (-got +want)\n%s",
					got, tt.want, gs, diff,
				)
			}
		})
	}
}
