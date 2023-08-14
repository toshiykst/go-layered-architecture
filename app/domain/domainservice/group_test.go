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
				s.AddGroups(model.NewGroup("TEST_GROUP_ID", "TEST_GROUP_NAME", []model.UserID{}))
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
