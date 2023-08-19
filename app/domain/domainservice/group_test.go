package domainservice

import (
	"testing"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/domain/repository"
	"github.com/toshiykst/go-layerd-architecture/app/infrastructure/memory"
)

func TestGroupService_Exists(t *testing.T) {
	tests := []struct {
		name                string
		newMemoryRepository func() repository.Repository
		gID                 model.GroupID
		want                bool
		wantErr             bool
	}{
		{
			name: "Returns true if the group id is an existing group's",
			gID:  model.GroupID("TEST_GROUP_ID"),
			want: true,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				s.AddGroups(model.MustNewGroup("TEST_GROUP_ID", "TEST_GROUP_NAME", []model.UserID{}))
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
		{
			name: "Returns false if the user id is not for an existing user",
			gID:  model.GroupID("TEST_GROUP_ID"),
			want: false,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &groupService{r: tt.newMemoryRepository()}
			got, err := gs.Exists(tt.gID)
			if err != nil {
				t.Fatalf("want no error, but has error %v", err)
			}
			if got != tt.want {
				t.Errorf("gs.Exists(%s)=%t, nil; want %t, nil", tt.gID, got, tt.want)
			}
		})
	}
}

func TestGroupService_HasUsersAny(t *testing.T) {
	tests := []struct {
		name                string
		newMemoryRepository func() repository.Repository
		uIDs                []model.UserID
		want                bool
		wantErr             bool
	}{
		{
			name: "Returns true if any groups have users specified by user ids of arg",
			uIDs: []model.UserID{
				"TEST_USER_ID_1",
				"TEST_USER_ID_2",
				"TEST_USER_ID_3",
			},
			want: true,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				s.AddGroups(model.MustNewGroup("TEST_GROUP_ID", "TEST_GROUP_NAME", []model.UserID{
					"TEST_USER_ID_1",
				}))
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
		{
			name: "Returns false if all groups do not have users specified by user ids of arg",
			uIDs: []model.UserID{
				"TEST_USER_ID_1",
				"TEST_USER_ID_2",
				"TEST_USER_ID_3",
			},
			want: false,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				s.AddGroups(model.MustNewGroup("TEST_GROUP_ID", "TEST_GROUP_NAME", []model.UserID{
					"TEST_USER_ID_4",
				}))
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &groupService{r: tt.newMemoryRepository()}
			got, err := gs.HasUsersAny(tt.uIDs)
			if err != nil {
				t.Fatalf("want no error, but has error %v", err)
			}
			if got != tt.want {
				t.Errorf("gs.HasUsersAny(%v)=%t, nil; want %t, nil", tt.uIDs, got, tt.want)
			}
		})
	}
}
