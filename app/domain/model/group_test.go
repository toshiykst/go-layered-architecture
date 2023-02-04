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
