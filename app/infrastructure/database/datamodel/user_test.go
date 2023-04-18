package datamodel

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

func TestNewUser(t *testing.T) {
	type args struct {
		id    model.UserID
		name  string
		email string
	}
	tests := []struct {
		name string
		args args
		want *User
	}{
		{
			name: "Creates a datamodel user",
			args: args{
				id:    model.UserID("TEST_USER_ID"),
				name:  "TEST_USER_NAME",
				email: "TEST_USER_EMAIL",
			},
			want: &User{
				ID:    "TEST_USER_ID",
				Name:  "TEST_USER_NAME",
				Email: "TEST_USER_EMAIL",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUser(tt.args.id, tt.args.name, tt.args.email)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf(
					"TestNewUser(%s, %s, %s)=%v; want %v\ndiffers: (-got +want)\n%s",
					tt.args.id, tt.args.name, tt.args.email, got, tt.want, diff,
				)
			}
		})
	}
}

func TestUser_ToModel(t *testing.T) {
	tests := []struct {
		name string
		user *User
		want *model.User
	}{
		{
			name: "Convert to model.User",
			user: &User{
				ID:    "TEST_USER_ID",
				Name:  "TEST_USER_NAME",
				Email: "TEST_USER_EMAIL",
			},
			want: model.NewUser("TEST_USER_ID", "TEST_USER_NAME", "TEST_USER_EMAIL"),
		},
		{
			name: "Returns nil when the receiver is nil",
			user: nil,
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := tt.user
			got := u.ToModel()
			if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(model.User{})); diff != "" {
				t.Errorf(
					"u.ToModel()=%v; want=%v,receiver=%v\ndiffers: (-got +want)\n%s",
					got, tt.want, u, diff,
				)
			}
		})
	}
}

func TestUsers_ToModel(t *testing.T) {
	tests := []struct {
		name  string
		users Users
		want  []*model.User
	}{
		{
			name: "Convert to model.User slice",
			users: Users{
				{
					ID:    "TEST_USER_ID_1",
					Name:  "TEST_USER_NAME_1",
					Email: "TEST_USER_EMAIL_1",
				},
				{
					ID:    "TEST_USER_ID_2",
					Name:  "TEST_USER_NAME_2",
					Email: "TEST_USER_EMAIL_2",
				},
				{
					ID:    "TEST_USER_ID_3",
					Name:  "TEST_USER_NAME_3",
					Email: "TEST_USER_EMAIL_3",
				},
			},
			want: model.Users{
				model.NewUser(
					"TEST_USER_ID_1",
					"TEST_USER_NAME_1",
					"TEST_USER_EMAIL_1",
				),
				model.NewUser(
					"TEST_USER_ID_2",
					"TEST_USER_NAME_2",
					"TEST_USER_EMAIL_2",
				),
				model.NewUser(
					"TEST_USER_ID_3",
					"TEST_USER_NAME_3",
					"TEST_USER_EMAIL_3",
				),
			},
		},
		{
			name:  "Returns nil when the receiver is nil",
			users: nil,
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := tt.users
			got := us.ToModel()
			if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(model.User{})); diff != "" {
				t.Errorf(
					"us.ToModel()=%v; want=%v,receiver=%v\ndiffers: (-got +want)\n%s",
					got, tt.want, us, diff,
				)
			}
		})
	}
}
