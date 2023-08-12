package response

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/toshiykst/go-layerd-architecture/app/usecase/dto"
)

func TestToUsersFromDTO(t *testing.T) {
	tests := []struct {
		name  string
		dtous []dto.User
		want  []User
	}{
		{
			name: "Returns response users",
			dtous: []dto.User{
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
			got := ToUsersFromDTO(tt.dtous)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf(
					"ToUsersFromDTO(%v)=%v; want %v\ndiffers: (-got +want)\n%s",
					tt.dtous, got, tt.want, diff,
				)
			}
		})
	}
}
