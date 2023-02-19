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
			us := &UserService{r: tt.newMockRepository()}
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
