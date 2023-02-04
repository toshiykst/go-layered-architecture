package model

import "testing"

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
