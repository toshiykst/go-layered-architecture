package model

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewUser(t *testing.T) {
	type args struct {
		id    UserID
		name  string
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr error
	}{
		{
			name: "Returns user",
			args: args{
				id:    "TEST_USER_ID",
				name:  "TEST_USER_NAME",
				email: "TEST_USER_EMAIL",
			},
			want: &User{
				id:    "TEST_USER_ID",
				name:  "TEST_USER_NAME",
				email: "TEST_USER_EMAIL",
			},
			wantErr: nil,
		},
		{
			name: "Error empty user id",
			args: args{
				id:    "",
				name:  "TEST_USER_NAME",
				email: "TEST_USER_EMAIL",
			},
			want:    nil,
			wantErr: ErrInvalidUser,
		},
		{
			name: "Error empty user name",
			args: args{
				id:    "TEST_USER_ID",
				name:  "",
				email: "TEST_USER_EMAIL",
			},
			want:    nil,
			wantErr: ErrInvalidUser,
		},
		{
			name: "Error exceeds the max user name length",
			args: args{
				id:    "TEST_USER_ID",
				name:  strings.Repeat("x", maxUserNameLength+1),
				email: "TEST_USER_EMAIL",
			},
			want:    nil,
			wantErr: ErrInvalidUser,
		},
		{
			name: "Error empty user email",
			args: args{
				id:    "TEST_USER_ID",
				name:  "TEST_USER_NAME",
				email: "",
			},
			want:    nil,
			wantErr: ErrInvalidUser,
		},
		{
			name: "Error exceeds the max user email length",
			args: args{
				id:    "TEST_USER_ID",
				name:  "TEST_USER_NAME",
				email: strings.Repeat("x", maxUserEmailLength+1),
			},
			want:    nil,
			wantErr: ErrInvalidUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.id, tt.args.name, tt.args.email)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("want an error, but has no error")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf(
						"NewUser(%s, %s, %s)=_, %v; want _, %v",
						tt.args.id, tt.args.name, tt.args.email, err, tt.wantErr,
					)
				}
			} else {
				if err != nil {
					t.Fatalf("want no error, but has error %v", err)
				}
				if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(User{})); diff != "" {
					t.Errorf(
						"NewUser(%s, %s, %s)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
						tt.args.id, tt.args.name, tt.args.email, err, tt.wantErr, diff,
					)
				}
			}
		})
	}
}

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
				MustNewUser(
					"TEST_USER_ID_1",
					"TEST_USER_NAME_1",
					"TEST_USER_EMAIL_1",
				),
				MustNewUser(
					"TEST_USER_ID_2",
					"TEST_USER_NAME_2",
					"TEST_USER_EMAIL_2",
				),
				MustNewUser(
					"TEST_USER_ID_3",
					"TEST_USER_NAME_3",
					"TEST_USER_EMAIL_3",
				),
			},
			want: map[UserID]*User{
				"TEST_USER_ID_1": MustNewUser(
					"TEST_USER_ID_1",
					"TEST_USER_NAME_1",
					"TEST_USER_EMAIL_1",
				),
				"TEST_USER_ID_2": MustNewUser(
					"TEST_USER_ID_2",
					"TEST_USER_NAME_2",
					"TEST_USER_EMAIL_2",
				),
				"TEST_USER_ID_3": MustNewUser(
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
