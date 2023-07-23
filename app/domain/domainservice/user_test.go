package domainservice

import (
	"testing"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/domain/repository"
	mockrepository "github.com/toshiykst/go-layerd-architecture/app/mock/domain/repository"
)

func TestUserService_Exists(t *testing.T) {
	tests := []struct {
		name              string
		newMockRepository func() repository.Repository
		uID               model.UserID
		want              bool
		wantErr           bool
	}{
		{
			name: "Returns true if the user id is for an existing user",
			uID:  model.UserID("TEST_USER_ID"),
			want: true,
			newMockRepository: func() repository.Repository {
				s := mockrepository.NewStore()
				s.AddUsers(model.NewUser("TEST_USER_ID", "TEST_USER_NAME", "TEST_USER_EMAIL"))
				r := mockrepository.NewMockRepository(s)
				return r
			},
		},
		{
			name: "Returns false if the user id is not for an existing user",
			uID:  model.UserID("TEST_USER_ID"),
			want: false,
			newMockRepository: func() repository.Repository {
				s := mockrepository.NewStore()
				r := mockrepository.NewMockRepository(s)
				return r
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &userService{r: tt.newMockRepository()}
			got, err := us.Exists(tt.uID)
			if err != nil {
				t.Fatalf("want no err, but has error %v", err)
			}
			if got != tt.want {
				t.Errorf("us.Exists(%s)=%t, nil; want %t, nil", tt.uID, got, tt.want)
			}
		})
	}
}

func TestUserService_ExistsAll(t *testing.T) {
	tests := []struct {
		name              string
		newMockRepository func() repository.Repository
		uIDs              []model.UserID
		want              bool
		wantErr           bool
	}{
		{
			name: "Returns true if user ids are existing",
			uIDs: []model.UserID{
				"TEST_USER_ID_1",
				"TEST_USER_ID_2",
				"TEST_USER_ID_3",
			},
			want: true,
			newMockRepository: func() repository.Repository {
				s := mockrepository.NewStore()
				s.AddUsers(
					model.NewUser("TEST_USER_ID_1", "TEST_USER_NAME_1", "TEST_USER_EMAIL_1"),
					model.NewUser("TEST_USER_ID_2", "TEST_USER_NAME_2", "TEST_USER_EMAIL_2"),
					model.NewUser("TEST_USER_ID_3", "TEST_USER_NAME_3", "TEST_USER_EMAIL_3"),
				)
				r := mockrepository.NewMockRepository(s)
				return r
			},
		},
		{
			name: "Returns false if some user ids are not existing",
			uIDs: []model.UserID{
				"TEST_USER_ID_1",
				"TEST_USER_ID_2",
				"TEST_USER_ID_4",
			},
			want: false,
			newMockRepository: func() repository.Repository {
				s := mockrepository.NewStore()
				s.AddUsers(
					model.NewUser("TEST_USER_ID_1", "TEST_USER_NAME_1", "TEST_USER_EMAIL_1"),
					model.NewUser("TEST_USER_ID_2", "TEST_USER_NAME_2", "TEST_USER_EMAIL_2"),
					model.NewUser("TEST_USER_ID_3", "TEST_USER_NAME_3", "TEST_USER_EMAIL_3"),
				)
				r := mockrepository.NewMockRepository(s)
				return r
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &userService{r: tt.newMockRepository()}
			got, err := us.ExistsAll(tt.uIDs)
			if err != nil {
				t.Fatalf("want no error, but has error %v", err)
			}
			if got != tt.want {
				t.Errorf("us.ExistsAll(%s)=%t, nil; want %t, nil", tt.uIDs, got, tt.want)
			}
		})
	}
}
