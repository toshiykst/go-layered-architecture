package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUser_ID(t *testing.T) {
	tests := []struct {
		name string
		user *User
		want UserID
	}{
		{
			name: "Returns user’s id",
			user: &User{id: UserID("TEST_USER_ID")},
			want: UserID("TEST_USER_ID"),
		},
		{
			name: "Receiver is nil",
			user: nil,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.user.ID(); got != tt.want {
				t.Errorf("ID()=%v; want %v, receiver=%v", got, tt.want, tt.user)
			}
		})
	}
}

func TestUser_Name(t *testing.T) {
	tests := []struct {
		name string
		user *User
		want string
	}{
		{
			name: "Returns user’s name",
			user: &User{name: "TEST_NAME"},
			want: "TEST_NAME",
		},
		{
			name: "Receiver is nil",
			user: nil,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.user.Name(); got != tt.want {
				t.Errorf("Name()=%v; want %v, receiver=%v", got, tt.want, tt.user)
			}
		})
	}
}

func TestUser_Email(t *testing.T) {
	tests := []struct {
		name string
		user *User
		want string
	}{
		{
			name: "Returns user’s email",
			user: &User{email: "TEST_EMAIL"},
			want: "TEST_EMAIL",
		},
		{
			name: "Receiver is nil",
			user: nil,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.user.Email(); got != tt.want {
				t.Errorf("Email()=%v; want %v, receiver=%v", got, tt.want, tt.user)
			}
		})
	}
}

func TestUsers_ByUserID(t *testing.T) {
	tests := []struct {
		name  string
		users Users
		want  map[UserID]*User
	}{
		{
			name: "Returns users map by user id",
			users: Users{
				NewUser(
					"TEST_USER_ID_1",
					"TEST_USER_NAME_1",
					"TEST_USER_EMAIL_1",
				),
				NewUser(
					"TEST_USER_ID_2",
					"TEST_USER_NAME_2",
					"TEST_USER_EMAIL_2",
				),
				NewUser(
					"TEST_USER_ID_3",
					"TEST_USER_NAME_3",
					"TEST_USER_EMAIL_3",
				),
			},
			want: map[UserID]*User{
				"TEST_USER_ID_1": NewUser(
					"TEST_USER_ID_1",
					"TEST_USER_NAME_1",
					"TEST_USER_EMAIL_1",
				),
				"TEST_USER_ID_2": NewUser(
					"TEST_USER_ID_2",
					"TEST_USER_NAME_2",
					"TEST_USER_EMAIL_2",
				),
				"TEST_USER_ID_3": NewUser(
					"TEST_USER_ID_3",
					"TEST_USER_NAME_3",
					"TEST_USER_EMAIL_3",
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := tt.users
			got := us.ByUserID()
			if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(User{})); diff != "" {
				t.Errorf(
					"us.ByUserID()=%v; want=%v,receiver=%v\ndiffers: (-got +want)\n%s",
					got, tt.want, us, diff,
				)
			}
		})
	}
}
