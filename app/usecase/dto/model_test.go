package dto

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

func TestToUsersFromModel(t *testing.T) {
	tests := []struct {
		name string
		us   model.Users
		want []User
	}{
		{
			name: "Convert model users to dto users",
			us: model.Users{
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
			want: []User{
				{
					UserID: "TEST_USER_ID_1",
					Name:   "TEST_USER_NAME_1",
					Email:  "TEST_USER_EMAIL_1",
				},
				{
					UserID: "TEST_USER_ID_2",
					Name:   "TEST_USER_NAME_2",
					Email:  "TEST_USER_EMAIL_2",
				},
				{
					UserID: "TEST_USER_ID_3",
					Name:   "TEST_USER_NAME_3",
					Email:  "TEST_USER_EMAIL_3",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToUsersFromModel(tt.us)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf(
					"ToUsersFromModel(%v)=%v; want %v\ndiffers: (-got +want)\n%s",
					tt.us, got, tt.want, diff,
				)
			}
		})
	}
}

func TestToModelUserIDs(t *testing.T) {
	tests := []struct {
		name string
		ids  []string
		want []model.UserID
	}{
		{
			name: "Convert string slice to model.UserID slice",
			ids: []string{
				"TEST_USER_ID_1",
				"TEST_USER_ID_2",
				"TEST_USER_ID_3",
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
			got := ToModelUserIDs(tt.ids)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf(
					"toModelUserIDs(%v)=%v; want %v\ndiffers: (-got +want)\n%s",
					tt.ids, got, tt.want, diff,
				)
			}
		})
	}
}
