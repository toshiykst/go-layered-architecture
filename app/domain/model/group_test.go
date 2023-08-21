package model_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

func TestNewGroup(t *testing.T) {
	type args struct {
		id   model.GroupID
		name string
		uIDs []model.UserID
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Group
		wantErr error
	}{
		{
			name: "Returns group",
			args: args{
				id:   "TEST_GROUP_ID",
				name: "TEST_GROUP_NAME",
				uIDs: []model.UserID{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_3",
					"TEST_USER_ID_4",
					"TEST_USER_ID_5",
				},
			},
			want: model.MustNewGroup(
				"TEST_GROUP_ID",
				"TEST_GROUP_NAME",
				[]model.UserID{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_3",
					"TEST_USER_ID_4",
					"TEST_USER_ID_5",
				},
			),
			wantErr: nil,
		},
		{
			name: "Error empty group id",
			args: args{
				id:   "",
				name: "TEST_GROUP_NAME",
				uIDs: []model.UserID{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_3",
					"TEST_USER_ID_4",
					"TEST_USER_ID_5",
				},
			},
			want:    nil,
			wantErr: model.ErrInvalidGroup,
		},
		{
			name: "Error empty group name",
			args: args{
				id:   "TEST_GROUP_ID",
				name: "",
				uIDs: []model.UserID{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_3",
					"TEST_USER_ID_4",
					"TEST_USER_ID_5",
				},
			},
			want:    nil,
			wantErr: model.ErrInvalidGroup,
		},
		{
			name: "Error exceeds the max group name length",
			args: args{
				id:   "TEST_GROUP_ID",
				name: strings.Repeat("x", model.MaxGroupNameLength+1),
				uIDs: []model.UserID{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_3",
					"TEST_USER_ID_4",
					"TEST_USER_ID_5",
				},
			},
			want:    nil,
			wantErr: model.ErrInvalidGroup,
		},
		{
			name: "Error exceeds the max group users",
			args: args{
				id:   "TEST_GROUP_ID",
				name: "TEST_GROUP_NAME",
				uIDs: []model.UserID{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_3",
					"TEST_USER_ID_4",
					"TEST_USER_ID_5",
					"TEST_USER_ID_6",
				},
			},
			want:    nil,
			wantErr: model.ErrInvalidGroup,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := model.NewGroup(tt.args.id, tt.args.name, tt.args.uIDs)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("want an error, but has no error")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf(
						"model.NewGroup(%s, %s, %v)=_, %v; want _, %v",
						tt.args.id, tt.args.name, tt.args.uIDs, err, tt.wantErr,
					)
				}
			} else {
				if err != nil {
					t.Fatalf("want no error, but has error %v", err)
				}
				if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(model.Group{})); diff != "" {
					t.Errorf(
						"model.NewGroup(%s, %s, %v)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
						tt.args.id, tt.args.name, tt.args.uIDs, err, tt.wantErr, diff,
					)
				}
			}
		})
	}
}

func TestGroup_ID(t *testing.T) {
	tests := []struct {
		name  string
		group *model.Group
		want  model.GroupID
	}{
		{
			name: "Returns group’s id",
			group: model.MustNewGroup(
				"TEST_GROUP_ID",
				"TEST_GROUP_NAME",
				[]model.UserID{},
			),
			want: model.GroupID("TEST_GROUP_ID"),
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
				t.Errorf("group.ID()=%v; want %v, receiver=%v", got, tt.want, tt.group)
			}
		})
	}
}

func TestGroup_Name(t *testing.T) {
	tests := []struct {
		name  string
		group *model.Group
		want  string
	}{
		{
			name: "Returns group’s name",
			group: model.MustNewGroup(
				"TEST_GROUP_ID",
				"TEST_GROUP_NAME",
				[]model.UserID{},
			),
			want: "TEST_GROUP_NAME",
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
				t.Errorf("group.Name()=%v; want %v, receiver=%v", got, tt.want, tt.group)
			}
		})
	}
}

func TestGroup_UserIDs(t *testing.T) {
	tests := []struct {
		name  string
		group *model.Group
		want  []model.UserID
	}{
		{
			name: "Returns group’s user ids",
			group: model.MustNewGroup(
				"TEST_GROUP_ID",
				"TEST_GROUP_NAME",
				[]model.UserID{"TEST_USER_ID_1", "TEST_USER_ID_2", "TEST_USER_ID_3"},
			),
			want: []model.UserID{"TEST_USER_ID_1", "TEST_USER_ID_2", "TEST_USER_ID_3"},
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
				t.Errorf("group.UserIDs()=%v; want %v, receiver=%v\ndiffers: (-got +want)\n%s", got, tt.want, tt.group, diff)
			}
		})
	}
}

func TestGroup_IsMaxUser(t *testing.T) {
	tests := []struct {
		name  string
		group *model.Group
		want  bool
	}{
		{
			name: "Returns true if the number of group users is max",
			group: model.MustNewGroup(
				"TEST_GROUP_ID",
				"TEST_GROUP_NAME",
				[]model.UserID{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_3",
					"TEST_USER_ID_4",
					"TEST_USER_ID_5",
				},
			),
			want: true,
		},
		{
			name: "Returns false if the number of group users is less than max",
			group: model.MustNewGroup(
				"TEST_GROUP_ID",
				"TEST_GROUP_NAME",
				[]model.UserID{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_3",
				},
			),
			want: false,
		},
		{
			name:  "Receiver is nil",
			group: nil,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.group.IsMaxUsers(); got != tt.want {
				t.Errorf("group.IsMaxUsers()=%t; want %v, receiver=%v", got, tt.want, tt.group)
			}
		})
	}
}

func TestGroups_IDs(t *testing.T) {
	tests := []struct {
		name   string
		groups model.Groups
		want   []model.GroupID
	}{
		{
			name: "Returns group ids",
			groups: model.Groups{
				model.MustNewGroup(
					"TEST_GROUP_ID_1",
					"TEST_GROUP_NAME_1",
					[]model.UserID{"TEST_USER_ID_1"},
				),
				model.MustNewGroup(
					"TEST_GROUP_ID_2",
					"TEST_GROUP_NAME_2",
					[]model.UserID{"TEST_USER_ID_2"},
				),
				model.MustNewGroup(
					"TEST_GROUP_ID_3",
					"TEST_GROUP_NAME_3",
					[]model.UserID{"TEST_USER_ID_3"},
				),
			},
			want: []model.GroupID{
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
		groups model.Groups
		want   []model.UserID
	}{
		{
			name: "Returns user ids",
			groups: model.Groups{
				model.MustNewGroup(
					"TEST_GROUP_ID_1",
					"TEST_GROUP_NAME_1",
					[]model.UserID{"TEST_USER_ID_1"},
				),
				model.MustNewGroup(
					"TEST_GROUP_ID_2",
					"TEST_GROUP_NAME_2",
					[]model.UserID{"TEST_USER_ID_1", "TEST_USER_ID_2"},
				),
				model.MustNewGroup(
					"TEST_GROUP_ID_3",
					"TEST_GROUP_NAME_3",
					[]model.UserID{"TEST_USER_ID_1", "TEST_USER_ID_2", "TEST_USER_ID_3"},
				),
			},
			want: []model.UserID{
				"TEST_USER_ID_1",
				"TEST_USER_ID_2",
				"TEST_USER_ID_3",
			},
		},
		{
			name: "Returns nil when all groups have no user ids",
			groups: model.Groups{
				model.MustNewGroup(
					"TEST_GROUP_ID_1",
					"TEST_GROUP_NAME_1",
					[]model.UserID{},
				),
				model.MustNewGroup(
					"TEST_GROUP_ID_2",
					"TEST_GROUP_NAME_2",
					[]model.UserID{},
				),
				model.MustNewGroup(
					"TEST_GROUP_ID_3",
					"TEST_GROUP_NAME_3",
					[]model.UserID{},
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
