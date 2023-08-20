package usecase

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	"github.com/toshiykst/go-layerd-architecture/app/domain/domainservice"
	"github.com/toshiykst/go-layerd-architecture/app/domain/factory"
	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/domain/repository"
	"github.com/toshiykst/go-layerd-architecture/app/infrastructure/memory"
	mockfactory "github.com/toshiykst/go-layerd-architecture/app/mock/domain/factory"
	"github.com/toshiykst/go-layerd-architecture/app/usecase/dto"
)

func TestGroupUsecase_CreateGroup(t *testing.T) {
	tests := []struct {
		name                string
		in                  *dto.CreateGroupInput
		want                *dto.CreateGroupOutput
		wantErr             error
		newMemoryRepository func() repository.Repository
		newMockFactory      func(ctrl *gomock.Controller) factory.GroupFactory
	}{
		{
			name: "Creates a new group with no users",
			in: &dto.CreateGroupInput{
				Name:    "TEST_GROUP_NAME",
				UserIDs: nil,
			},
			want: &dto.CreateGroupOutput{
				Group: dto.Group{
					GroupID: "TEST_GROUP_ID",
					Name:    "TEST_GROUP_NAME",
					Users:   []dto.User{},
				},
			},
			wantErr: nil,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				r := memory.NewMemoryRepository(s)
				return r
			},
			newMockFactory: func(ctrl *gomock.Controller) factory.GroupFactory {
				f := mockfactory.NewMockGroupFactory(ctrl)
				f.EXPECT().
					Create("TEST_GROUP_NAME", []model.UserID{}).
					Return(model.MustNewGroup("TEST_GROUP_ID", "TEST_GROUP_NAME", []model.UserID{}), nil)
				return f
			},
		},
		{
			name: "Creates a new group with users",
			in: &dto.CreateGroupInput{
				Name: "TEST_GROUP_NAME",
				UserIDs: []string{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_3",
				},
			},
			want: &dto.CreateGroupOutput{
				Group: dto.Group{
					GroupID: "TEST_GROUP_ID",
					Name:    "TEST_GROUP_NAME",
					Users: []dto.User{
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
			},
			wantErr: nil,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				s.AddUsers(
					model.MustNewUser("TEST_USER_ID_1", "TEST_USER_NAME_1", "TEST_USER_EMAIL_1"),
					model.MustNewUser("TEST_USER_ID_2", "TEST_USER_NAME_2", "TEST_USER_EMAIL_2"),
					model.MustNewUser("TEST_USER_ID_3", "TEST_USER_NAME_3", "TEST_USER_EMAIL_3"),
				)
				r := memory.NewMemoryRepository(s)
				return r
			},
			newMockFactory: func(ctrl *gomock.Controller) factory.GroupFactory {
				uIDs := []model.UserID{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_3",
				}
				f := mockfactory.NewMockGroupFactory(ctrl)
				f.EXPECT().
					Create("TEST_GROUP_NAME", uIDs).
					Return(model.MustNewGroup("TEST_GROUP_ID", "TEST_GROUP_NAME", uIDs), nil)
				return f
			},
		},
		{
			name: "Returns error if any of group inputs are invalid",
			in: &dto.CreateGroupInput{
				Name:    "TEST_GROUP_NAME_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				UserIDs: []string{},
			},
			want:    nil,
			wantErr: ErrInvalidGroupInput,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				s.AddUsers(
					model.MustNewUser("TEST_USER_ID_1", "TEST_USER_NAME_1", "TEST_USER_EMAIL_1"),
					model.MustNewUser("TEST_USER_ID_2", "TEST_USER_NAME_2", "TEST_USER_EMAIL_2"),
					model.MustNewUser("TEST_USER_ID_3", "TEST_USER_NAME_3", "TEST_USER_EMAIL_3"),
				)
				r := memory.NewMemoryRepository(s)
				return r
			},
			newMockFactory: func(ctrl *gomock.Controller) factory.GroupFactory {
				f := mockfactory.NewMockGroupFactory(ctrl)
				f.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, model.ErrInvalidGroup)
				return f
			},
		},
		{
			name: "Returns error if some of user ids are not existed",
			in: &dto.CreateGroupInput{
				Name: "TEST_GROUP_NAME",
				UserIDs: []string{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_4",
				},
			},
			want:    nil,
			wantErr: ErrInvalidUserIDs,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				s.AddUsers(
					model.MustNewUser("TEST_USER_ID_1", "TEST_USER_NAME_1", "TEST_USER_EMAIL_1"),
					model.MustNewUser("TEST_USER_ID_2", "TEST_USER_NAME_2", "TEST_USER_EMAIL_2"),
					model.MustNewUser("TEST_USER_ID_3", "TEST_USER_NAME_3", "TEST_USER_EMAIL_3"),
				)
				r := memory.NewMemoryRepository(s)
				return r
			},
			newMockFactory: func(ctrl *gomock.Controller) factory.GroupFactory {
				uIDs := []model.UserID{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_4",
				}
				f := mockfactory.NewMockGroupFactory(ctrl)
				f.EXPECT().
					Create("TEST_GROUP_NAME", uIDs).
					Return(model.MustNewGroup("TEST_GROUP_ID", "TEST_GROUP_NAME", uIDs), nil)
				return f
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := tt.newMemoryRepository()
			gs := domainservice.NewGroupService(r)
			us := domainservice.NewUserService(r)
			uc := NewGroupUsecase(r, tt.newMockFactory(ctrl), gs, us)

			got, err := uc.CreateGroup(tt.in)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("want an error, but has no error")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf(
						"uc.CreateGroup(%v)=_, %v; want _, %v",
						tt.in, got, tt.wantErr,
					)
				}
			} else {
				if err != nil {
					t.Fatalf("want no error, but has error %v", err)
				}
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf(
						"uc.CreateGroup(%v)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
						tt.in, got, tt.want, diff,
					)
				}
			}
		})
	}
}

func TestGroupUsecase_GetGroup(t *testing.T) {
	tests := []struct {
		name                string
		in                  *dto.GetGroupInput
		want                *dto.GetGroupOutput
		wantErr             error
		newMemoryRepository func() repository.Repository
	}{
		{
			name: "Returns a group",
			in: &dto.GetGroupInput{
				GroupID: "TEST_GROUP_ID",
			},
			want: &dto.GetGroupOutput{
				Group: dto.Group{
					GroupID: "TEST_GROUP_ID",
					Name:    "TEST_GROUP_NAME",
					Users: []dto.User{
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
			},
			wantErr: nil,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				s.AddGroups(model.MustNewGroup(
					"TEST_GROUP_ID",
					"TEST_GROUP_NAME",
					[]model.UserID{
						"TEST_USER_ID_1",
						"TEST_USER_ID_2",
						"TEST_USER_ID_3",
					}),
				)
				s.AddUsers(
					model.MustNewUser("TEST_USER_ID_1", "TEST_USER_NAME_1", "TEST_USER_EMAIL_1"),
					model.MustNewUser("TEST_USER_ID_2", "TEST_USER_NAME_2", "TEST_USER_EMAIL_2"),
					model.MustNewUser("TEST_USER_ID_3", "TEST_USER_NAME_3", "TEST_USER_EMAIL_3"),
				)
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
		{
			name: "Returns error if the group does not exist",
			in: &dto.GetGroupInput{
				GroupID: "TEST_GROUP_ID",
			},
			want:    nil,
			wantErr: ErrGroupNotFound,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := tt.newMemoryRepository()
			gs := domainservice.NewGroupService(r)
			us := domainservice.NewUserService(r)
			uc := NewGroupUsecase(r, mockfactory.NewMockGroupFactory(ctrl), gs, us)

			got, err := uc.GetGroup(tt.in)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("want an error, but has no error")
				}

				if err.Error() != tt.wantErr.Error() {
					t.Errorf(
						"uc.GetGroup(%v)=_, %v; want _, %v",
						tt.in, got, tt.wantErr,
					)
				}
			} else {
				if err != nil {
					t.Fatalf("want no err, but has error %v", err)
				}
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf(
						"uc.GetGroup(%v)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
						tt.in, got, tt.want, diff,
					)
				}
			}
		})
	}
}

func TestGroupUsecase_GetGroups(t *testing.T) {
	tests := []struct {
		name                string
		want                *dto.GetGroupsOutput
		wantErr             error
		newMemoryRepository func() repository.Repository
	}{
		{
			name: "Returns groups",
			want: &dto.GetGroupsOutput{
				Groups: []dto.Group{
					{
						GroupID: "TEST_GROUP_ID_1",
						Name:    "TEST_GROUP_NAME_1",
						Users: []dto.User{
							{
								UserID: "TEST_USER_ID_1",
								Name:   "TEST_USER_NAME_1",
								Email:  "TEST_USER_EMAIL_1",
							},
						},
					},
					{
						GroupID: "TEST_GROUP_ID_2",
						Name:    "TEST_GROUP_NAME_2",
						Users: []dto.User{
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
						},
					},
					{
						GroupID: "TEST_GROUP_ID_3",
						Name:    "TEST_GROUP_NAME_3",
						Users: []dto.User{
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
				},
			},
			wantErr: nil,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				s.AddGroups(
					model.MustNewGroup(
						"TEST_GROUP_ID_1",
						"TEST_GROUP_NAME_1",
						[]model.UserID{
							"TEST_USER_ID_1",
						},
					),
					model.MustNewGroup(
						"TEST_GROUP_ID_2",
						"TEST_GROUP_NAME_2",
						[]model.UserID{
							"TEST_USER_ID_1",
							"TEST_USER_ID_2",
						},
					),
					model.MustNewGroup(
						"TEST_GROUP_ID_3",
						"TEST_GROUP_NAME_3",
						[]model.UserID{
							"TEST_USER_ID_1",
							"TEST_USER_ID_2",
							"TEST_USER_ID_3",
						},
					),
				)
				s.AddUsers(
					model.MustNewUser("TEST_USER_ID_1", "TEST_USER_NAME_1", "TEST_USER_EMAIL_1"),
					model.MustNewUser("TEST_USER_ID_2", "TEST_USER_NAME_2", "TEST_USER_EMAIL_2"),
					model.MustNewUser("TEST_USER_ID_3", "TEST_USER_NAME_3", "TEST_USER_EMAIL_3"),
				)
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
		{
			name: "Returns empty groups if not exist",
			want: &dto.GetGroupsOutput{
				Groups: []dto.Group{},
			},
			wantErr: nil,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
		{
			name: "Returns groups if all groups have no user ids",
			want: &dto.GetGroupsOutput{
				Groups: []dto.Group{
					{
						GroupID: "TEST_GROUP_ID_1",
						Name:    "TEST_GROUP_NAME_1",
						Users:   []dto.User{},
					},
					{
						GroupID: "TEST_GROUP_ID_2",
						Name:    "TEST_GROUP_NAME_2",
						Users:   []dto.User{},
					},
					{
						GroupID: "TEST_GROUP_ID_3",
						Name:    "TEST_GROUP_NAME_3",
						Users:   []dto.User{},
					},
				},
			},
			wantErr: nil,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				s.AddGroups(
					model.MustNewGroup(
						"TEST_GROUP_ID_1",
						"TEST_GROUP_NAME_1",
						[]model.UserID{},
					),
					model.MustNewGroup(
						"TEST_GROUP_ID_2",
						"TEST_GROUP_NAME_2",
						[]model.UserID{},
					),
					model.MustNewGroup(
						"TEST_GROUP_ID_3",
						"TEST_GROUP_NAME_3",
						[]model.UserID{},
					),
				)
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := tt.newMemoryRepository()
			gs := domainservice.NewGroupService(r)
			us := domainservice.NewUserService(r)
			uc := NewGroupUsecase(r, mockfactory.NewMockGroupFactory(ctrl), gs, us)

			in := &dto.GetGroupsInput{}
			got, err := uc.GetGroups(in)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("want an error, but has no error")
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf(
						"uc.GetGroups(%v)=_, %v; want _, %v",
						in, got, tt.wantErr,
					)
				}
			} else {
				if err != nil {
					t.Fatalf("want no err, but has error %v", err)
				}
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf(
						"uc.GetGroups(%v)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
						in, got, tt.want, diff,
					)
				}
			}
		})
	}
}

func TestGroupUsecase_UpdateGroup(t *testing.T) {
	tests := []struct {
		name                string
		in                  *dto.UpdateGroupInput
		wantGroup           *model.Group
		newMemoryRepository func() repository.Repository
		wantErr             error
	}{
		{
			name: "Update a group",
			in: &dto.UpdateGroupInput{
				GroupID: "TEST_GROUP_ID",
				Name:    "TEST_GROUP_NAME_UPDATED",
			},
			wantGroup: model.MustNewGroup(
				"TEST_GROUP_ID",
				"TEST_GROUP_NAME_UPDATED",
				[]model.UserID{},
			),
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				s.AddGroups(model.MustNewGroup(
					"TEST_GROUP_ID",
					"TEST_GROUP_NAME_UPDATED",
					[]model.UserID{},
				))
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
		{
			name: "Returns error if any of group inputs are invalid",
			in: &dto.UpdateGroupInput{
				GroupID: "TEST_GROUP_ID",
				Name:    "TEST_GROUP_NAME_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			},
			wantGroup: nil,
			wantErr:   ErrInvalidGroupInput,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
		{
			name: "Returns error if the group does not exist",
			in: &dto.UpdateGroupInput{
				GroupID: "TEST_GROUP_ID",
				Name:    "TEST_GROUP_NAME_UPDATED",
			},
			wantGroup: nil,
			wantErr:   ErrGroupNotFound,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := mockfactory.NewMockGroupFactory(ctrl)
			r := tt.newMemoryRepository()
			gs := domainservice.NewGroupService(r)
			us := domainservice.NewUserService(r)
			uc := NewGroupUsecase(r, f, gs, us)

			_, err := uc.UpdateGroup(tt.in)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("want an error, but has no error")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf(
						"uc.UpdateGroup(%v)=_, %v; want _, %v",
						tt.in, err, tt.wantErr,
					)
				}
			} else {
				if err != nil {
					t.Fatalf("want no error, but has error %v", err)
				}
				gID := model.GroupID(tt.in.GroupID)
				got, _ := r.Group().Find(gID)
				if diff := cmp.Diff(got, tt.wantGroup, cmp.AllowUnexported(model.Group{})); diff != "" {
					t.Errorf(
						"r.Group().Find(%s)=%v, _; want %v, nil\ndiffers: (-got +want)\n%s",
						gID, got, tt.wantGroup, diff,
					)
				}
			}
		})
	}
}

func TestGroupUsecase_DeleteGroup(t *testing.T) {
	tests := []struct {
		name                string
		in                  *dto.DeleteGroupInput
		newMemoryRepository func() repository.Repository
		wantErr             error
	}{
		{
			name: "Delete a group",
			in: &dto.DeleteGroupInput{
				GroupID: "TEST_GROUP_ID",
			},
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				s.AddGroups(model.MustNewGroup(
					"TEST_GROUP_ID",
					"TEST_GROUP_NAME",
					[]model.UserID{},
				))
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
		{
			name: "Delete a group that has users",
			in: &dto.DeleteGroupInput{
				GroupID: "TEST_GROUP_ID",
			},
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				s.AddGroups(model.MustNewGroup(
					"TEST_GROUP_ID",
					"TEST_GROUP_NAME",
					[]model.UserID{
						"TEST_USER_ID_1",
						"TEST_USER_ID_2",
						"TEST_USER_ID_3",
					},
				))
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
		{
			name: "Returns error if the group does not exist",
			in: &dto.DeleteGroupInput{
				GroupID: "TEST_GROUP_ID",
			},
			wantErr: ErrGroupNotFound,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := mockfactory.NewMockGroupFactory(ctrl)
			r := tt.newMemoryRepository()
			gs := domainservice.NewGroupService(r)
			us := domainservice.NewUserService(r)
			uc := NewGroupUsecase(r, f, gs, us)

			_, err := uc.DeleteGroup(tt.in)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("want an error, but has no error")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf(
						"uc.DeleteGroup(%v)=_, %v; want _, %v",
						tt.in, err, tt.wantErr,
					)
				}
			} else {
				if err != nil {
					t.Fatalf("want no error, but has error %v", err)
				}
				gID := model.GroupID(tt.in.GroupID)
				got, _ := r.Group().Find(gID)
				if got != nil {
					t.Errorf("r.Group().Find(%s)=%v, _; want nil, nil", gID, got)
				}
			}
		})
	}
}
