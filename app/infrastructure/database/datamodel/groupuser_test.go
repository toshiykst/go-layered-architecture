package datamodel_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/infrastructure/database/datamodel"
)

func TestNewGroupUser(t *testing.T) {
	type args struct {
		gID model.GroupID
		uID model.UserID
	}
	tests := []struct {
		name string
		args args
		want *datamodel.GroupUser
	}{
		{
			name: "Creates a datamodel groupuser",
			args: args{
				gID: model.GroupID("TEST_GROUP_ID"),
				uID: model.UserID("TEST_USER_ID"),
			},
			want: &datamodel.GroupUser{
				GroupID: "TEST_GROUP_ID",
				UserID:  "TEST_USER_ID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := datamodel.NewGroupUser(tt.args.gID, tt.args.uID)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf(
					"NewGroupUser(%s,%s)=%v; want %v\ndiffers: (-got +want)\n%s",
					tt.args.gID, tt.args.uID, got, tt.want, diff,
				)
			}
		})
	}
}

func TestNewGroupUsers(t *testing.T) {
	type args struct {
		gID  model.GroupID
		uIDs []model.UserID
	}
	tests := []struct {
		name string
		args args
		want datamodel.GroupUsers
	}{
		{
			name: "Creates a datamodel groupusers",
			args: args{
				gID: model.GroupID("TEST_GROUP_ID"),
				uIDs: []model.UserID{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_3",
				},
			},
			want: datamodel.GroupUsers{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := datamodel.NewGroupUsers(tt.args.gID, tt.args.uIDs)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf(
					"NewGroupUsers(%s,%v)=%v; want %v\ndiffers: (-got +want)\n%s",
					tt.args.gID, tt.args.uIDs, got, tt.want, diff,
				)
			}
		})
	}
}

func TestGroupUsers_GroupIDs(t *testing.T) {
	tests := []struct {
		name string
		gus  datamodel.GroupUsers
		want []string
	}{
		{
			name: "Returns unique group ids",
			gus: datamodel.GroupUsers{
				{
					GroupID: "TEST_GROUP_ID_1",
					UserID:  "TEST_USER_ID_1",
				},
				{
					GroupID: "TEST_GROUP_ID_1",
					UserID:  "TEST_USER_ID_2",
				},
				{
					GroupID: "TEST_GROUP_ID_2",
					UserID:  "TEST_USER_ID_3",
				},
			},
			want: []string{
				"TEST_GROUP_ID_1",
				"TEST_GROUP_ID_2",
			},
		},
		{
			name: "Returns nil when group users is nil",
			gus:  nil,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.gus.GroupIDs()
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf(
					"gus.GroupIDs()=%v; want %v, receiver=%v\ndiffers: (-got +want)\n%s",
					got, tt.want, tt.gus, diff,
				)
			}
		})
	}
}

func TestGroupUsers_ModelUserIDs(t *testing.T) {
	tests := []struct {
		name string
		gus  datamodel.GroupUsers
		want []model.UserID
	}{
		{
			name: "Returns user ids",
			gus: datamodel.GroupUsers{
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
			want: []model.UserID{
				"TEST_USER_ID_1",
				"TEST_USER_ID_2",
				"TEST_USER_ID_3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.gus.ModelUserIDs()
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf(
					"gus.UserIDs()=%v; want %v, receiver=%v\ndiffers: (-got +want)\n%s",
					got, tt.want, tt.gus, diff,
				)
			}
		})
	}
}

func TestGroupUsers_ModelUserIDsByGroupID(t *testing.T) {
	tests := []struct {
		name string
		gus  datamodel.GroupUsers
		want map[string][]model.UserID
	}{
		{
			name: "Returns map of group ids by user id",
			gus: datamodel.GroupUsers{
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
				{
					GroupID: "TEST_GROUP_ID_1",
					UserID:  "TEST_USER_ID_4",
				},
				{
					GroupID: "TEST_GROUP_ID_2",
					UserID:  "TEST_USER_ID_5",
				},
			},
			want: map[string][]model.UserID{
				"TEST_GROUP_ID_1": {"TEST_USER_ID_1", "TEST_USER_ID_4"},
				"TEST_GROUP_ID_2": {"TEST_USER_ID_2", "TEST_USER_ID_5"},
				"TEST_GROUP_ID_3": {"TEST_USER_ID_3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.gus.ModelUserIDsByGroupID()
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf(
					"gus.ModelUserIDsByGroupID()=%v; want %v, receiver=%v\ndiffers: (-got +want)\n%s",
					got, tt.want, tt.gus, diff,
				)
			}
		})
	}
}
