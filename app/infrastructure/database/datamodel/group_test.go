package datamodel

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

func TestNewGroup(t *testing.T) {
	type args struct {
		id   model.GroupID
		name string
	}
	tests := []struct {
		name string
		args args
		want *Group
	}{
		{
			name: "Creates a datamodel user",
			args: args{
				id:   model.GroupID("TEST_GROUP_ID"),
				name: "TEST_GROUP_NAME",
			},
			want: &Group{
				ID:   "TEST_GROUP_ID",
				Name: "TEST_GROUP_NAME",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGroup(tt.args.id, tt.args.name)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf(
					"NewGroup(%s,%s)=%v; want %v\ndiffers: (-got +want)\n%s",
					tt.args.id, tt.args.name, got, tt.want, diff,
				)
			}
		})
	}
}

func TestGroup_ToModel(t *testing.T) {
	type args struct {
		gus GroupUsers
	}
	tests := []struct {
		name  string
		group *Group
		args  args
		want  *model.Group
	}{
		{
			name: "Convert to model.Group",
			group: &Group{
				ID:   "TEST_GROUP_ID",
				Name: "TEST_GROUP_NAME",
			},
			args: args{
				gus: GroupUsers{
					{
						GroupID: "TEST_GROUP_ID",
						UserID:  "TEST_USER_ID_1",
					},
					{
						GroupID: "TEST_GROUP_ID",
						UserID:  "TEST_USER_ID_2",
					},
					{
						GroupID: "TEST_GROUP_ID",
						UserID:  "TEST_USER_ID_3",
					},
				},
			},
			want: model.NewGroup(
				"TEST_GROUP_ID",
				"TEST_GROUP_NAME",
				[]model.UserID{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_3",
				},
			),
		},
		{
			name:  "Returns nil when the receiver is nil",
			group: nil,
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.group
			got := g.ToModel(tt.args.gus)
			if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(model.Group{})); diff != "" {
				t.Errorf(
					"g.ToModel(%v)=%v; want=%v,receiver=%v\ndiffers: (-got +want)\n%s",
					tt.args.gus, got, tt.want, g, diff,
				)
			}
		})
	}
}

func TestGroups_IDs(t *testing.T) {
	tests := []struct {
		name   string
		groups Groups
		want   []string
	}{
		{
			name: "Returns group ids",
			groups: Groups{
				{
					ID:   "TEST_GROUP_ID_1",
					Name: "TEST_GROUP_NAME_1",
				},
				{
					ID:   "TEST_GROUP_ID_2",
					Name: "TEST_GROUP_NAME_2",
				},
				{
					ID:   "TEST_GROUP_ID_3",
					Name: "TEST_GROUP_NAME_3",
				},
			},
			want: []string{
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

func TestGroups_ToModel(t *testing.T) {
	type args struct {
		gus GroupUsers
	}
	tests := []struct {
		name   string
		groups Groups
		args   args
		want   model.Groups
	}{
		{
			name: "Convert to model.Groups",
			groups: Groups{
				{
					ID:   "TEST_GROUP_ID_1",
					Name: "TEST_GROUP_NAME_1",
				},
				{
					ID:   "TEST_GROUP_ID_2",
					Name: "TEST_GROUP_NAME_2",
				},
				{
					ID:   "TEST_GROUP_ID_3",
					Name: "TEST_GROUP_NAME_3",
				},
			},
			args: args{
				gus: GroupUsers{
					{
						GroupID: "TEST_GROUP_ID_1",
						UserID:  "TEST_USER_ID_1",
					},
					{
						GroupID: "TEST_GROUP_ID_2",
						UserID:  "TEST_USER_ID_2",
					},
					{
						GroupID: "TEST_GROUP_ID_3",
						UserID:  "TEST_USER_ID_3",
					},
				},
			},
			want: model.Groups{
				model.NewGroup(
					"TEST_GROUP_ID_1",
					"TEST_GROUP_NAME_1",
					[]model.UserID{"TEST_USER_ID_1"},
				),
				model.NewGroup(
					"TEST_GROUP_ID_2",
					"TEST_GROUP_NAME_2",
					[]model.UserID{"TEST_USER_ID_2"},
				),
				model.NewGroup(
					"TEST_GROUP_ID_3",
					"TEST_GROUP_NAME_3",
					[]model.UserID{"TEST_USER_ID_3"},
				),
			},
		},
		{
			name: "Returns nil when the receiver is nil",
			args: args{
				gus: GroupUsers{},
			},
			groups: nil,
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := tt.groups
			got := gs.ToModel(tt.args.gus)
			if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(model.Group{})); diff != "" {
				t.Errorf(
					"gs.ToModel(%v)=%v; want=%v,receiver=%v\ndiffers: (-got +want)\n%s",
					tt.args.gus, got, tt.want, gs, diff,
				)
			}
		})
	}
}
