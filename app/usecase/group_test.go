package usecase

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"

	"github.com/toshiykst/go-layerd-architecture/app/domain/domainservice"
	"github.com/toshiykst/go-layerd-architecture/app/domain/factory"
	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/domain/repository"
	mockdomainservice "github.com/toshiykst/go-layerd-architecture/app/mock/domain/domainservice"
	mockfactory "github.com/toshiykst/go-layerd-architecture/app/mock/domain/factory"
	mockrepository "github.com/toshiykst/go-layerd-architecture/app/mock/domain/repository"
	"github.com/toshiykst/go-layerd-architecture/app/usecase/dto"
)

func TestGroupUsecase_CreateGroup(t *testing.T) {
	tests := []struct {
		name              string
		in                *dto.CreateGroupInput
		want              *dto.CreateGroupOutput
		wantErr           error
		newMockRepository func() repository.Repository
		newMockFactory    func(ctrl *gomock.Controller) factory.GroupFactory
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
			newMockFactory: func(ctrl *gomock.Controller) factory.GroupFactory {
				f := mockfactory.NewMockGroupFactory(ctrl)
				f.EXPECT().
					Create("TEST_GROUP_NAME").
					Return(model.NewGroup("TEST_GROUP_ID", "TEST_GROUP_NAME", []model.UserID{}), nil)
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
			newMockFactory: func(ctrl *gomock.Controller) factory.GroupFactory {
				f := mockfactory.NewMockGroupFactory(ctrl)
				f.EXPECT().
					Create("TEST_GROUP_NAME").
					Return(model.NewGroup("TEST_GROUP_ID", "TEST_GROUP_NAME", []model.UserID{}), nil)
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
			newMockFactory: func(ctrl *gomock.Controller) factory.GroupFactory {
				f := mockfactory.NewMockGroupFactory(ctrl)
				f.EXPECT().
					Create("TEST_GROUP_NAME").
					Return(model.NewGroup("TEST_GROUP_ID", "TEST_GROUP_NAME", []model.UserID{}), nil)
				return f
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			r := tt.newMockRepository()
			uc := NewGroupUsecase(
				tt.newMockRepository(),
				tt.newMockFactory(ctrl),
				domainservice.NewUserService(r),
			)
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
		name              string
		in                *dto.GetGroupInput
		want              *dto.GetGroupOutput
		wantErr           error
		newMockRepository func() repository.Repository
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
			newMockRepository: func() repository.Repository {
				s := mockrepository.NewStore()
				s.AddGroups(model.NewGroup(
					"TEST_GROUP_ID",
					"TEST_GROUP_NAME",
					[]model.UserID{
						"TEST_USER_ID_1",
						"TEST_USER_ID_2",
						"TEST_USER_ID_3",
					}),
				)
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
			name: "Returns error if the group does not exist",
			in: &dto.GetGroupInput{
				GroupID: "TEST_GROUP_ID",
			},
			want:    nil,
			wantErr: ErrGroupNotFound,
			newMockRepository: func() repository.Repository {
				s := mockrepository.NewStore()
				r := mockrepository.NewMockRepository(s)
				return r
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			r := tt.newMockRepository()
			uc := NewGroupUsecase(r, mockfactory.NewMockGroupFactory(ctrl), mockdomainservice.NewMockUserService(ctrl))
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
