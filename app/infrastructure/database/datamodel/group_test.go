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
