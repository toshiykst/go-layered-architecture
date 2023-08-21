package usecase_test

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
	"github.com/toshiykst/go-layerd-architecture/app/usecase"
	"github.com/toshiykst/go-layerd-architecture/app/usecase/dto"
)

func TestUserUsecase_CreateUser(t *testing.T) {
	tests := []struct {
		name                string
		in                  *dto.CreateUserInput
		want                *dto.CreateUserOutput
		wantErr             error
		newMemoryRepository func() repository.Repository
		newMockFactory      func(ctrl *gomock.Controller) factory.UserFactory
	}{
		{
			name: "Creates a new user",
			in: &dto.CreateUserInput{
				Name:  "TEST_USER_NAME",
				Email: "TEST_USER_EMAIL",
			},
			want: &dto.CreateUserOutput{
				User: dto.User{
					UserID: "TEST_USER_ID",
					Name:   "TEST_USER_NAME",
					Email:  "TEST_USER_EMAIL",
				},
			},
			wantErr: nil,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				r := memory.NewMemoryRepository(s)
				return r
			},
			newMockFactory: func(ctrl *gomock.Controller) factory.UserFactory {
				f := mockfactory.NewMockUserFactory(ctrl)
				f.EXPECT().
					Create("TEST_USER_NAME", "TEST_USER_EMAIL").
					Return(model.MustNewUser("TEST_USER_ID", "TEST_USER_NAME", "TEST_USER_EMAIL"), nil)
				return f
			},
		},
		{
			name: "Returns error if any of user inputs are invalid",
			in: &dto.CreateUserInput{
				Name:  "TEST_USER_NAME_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				Email: "TEST_USER_EMAIL",
			},
			want: &dto.CreateUserOutput{
				User: dto.User{
					UserID: "TEST_USER_ID",
					Name:   "TEST_USER_NAME",
					Email:  "TEST_USER_EMAIL",
				},
			},
			wantErr: usecase.ErrInvalidUserInput,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				r := memory.NewMemoryRepository(s)
				return r
			},
			newMockFactory: func(ctrl *gomock.Controller) factory.UserFactory {
				f := mockfactory.NewMockUserFactory(ctrl)
				f.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, model.ErrInvalidUser)
				return f
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := tt.newMemoryRepository()
			us := domainservice.NewUserService(r)
			gs := domainservice.NewGroupService(r)
			uc := usecase.NewUserUsecase(r, tt.newMockFactory(ctrl), us, gs)

			got, err := uc.CreateUser(tt.in)

			if tt.wantErr != nil {
				if err == nil {
					t.Error("want an error, but has no error")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf(
						"uc.CreateUser(%v)=_, %v; want _, %v",
						tt.in, got, tt.wantErr,
					)
				}
			} else {
				if err != nil {
					t.Fatalf("want no err, but has error %v", err)
				}
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf(
						"uc.CreateUser(%v)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
						tt.in, got, tt.want, diff,
					)
				}
			}
		})
	}
}

func TestUserUsecase_GetUser(t *testing.T) {
	tests := []struct {
		name                string
		in                  *dto.GetUserInput
		want                *dto.GetUserOutput
		wantErr             error
		newMemoryRepository func() repository.Repository
	}{
		{
			name: "Returns a user",
			in: &dto.GetUserInput{
				UserID: "TEST_USER_ID",
			},
			want: &dto.GetUserOutput{
				User: dto.User{
					UserID: "TEST_USER_ID",
					Name:   "TEST_USER_NAME",
					Email:  "TEST_USER_EMAIL",
				},
			},
			wantErr: nil,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				s.AddUsers(model.MustNewUser("TEST_USER_ID", "TEST_USER_NAME", "TEST_USER_EMAIL"))
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
		{
			name: "Returns error if the user does not exist",
			in: &dto.GetUserInput{
				UserID: "TEST_USER_ID",
			},
			want:    nil,
			wantErr: usecase.ErrUserNotFound,
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
			us := domainservice.NewUserService(r)
			gs := domainservice.NewGroupService(r)
			uc := usecase.NewUserUsecase(r, mockfactory.NewMockUserFactory(ctrl), us, gs)

			got, err := uc.GetUser(tt.in)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("want an error, but has no error")
				}

				if err.Error() != tt.wantErr.Error() {
					t.Errorf(
						"uc.GetUser(%v)=_, %v; want _, %v",
						tt.in, got, tt.wantErr,
					)
				}
			} else {
				if err != nil {
					t.Fatalf("want no err, but has error %v", err)
				}
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf(
						"uc.GetUser(%v)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
						tt.in, got, tt.want, diff,
					)
				}
			}
		})
	}
}

func TestUserUsecase_GetUsers(t *testing.T) {
	tests := []struct {
		name                string
		in                  *dto.GetUsersInput
		want                *dto.GetUsersOutput
		wantErr             error
		newMemoryRepository func() repository.Repository
	}{
		{
			name: "Returns users",
			in:   nil,
			want: &dto.GetUsersOutput{
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := tt.newMemoryRepository()
			us := domainservice.NewUserService(r)
			gs := domainservice.NewGroupService(r)
			uc := usecase.NewUserUsecase(r, mockfactory.NewMockUserFactory(ctrl), us, gs)

			got, err := uc.GetUsers(tt.in)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("want an error, but has no error")
				}

				if err.Error() != tt.wantErr.Error() {
					t.Errorf(
						"uc.GetUsers(%v)=_, %v; want _, %v",
						tt.in, got, tt.wantErr,
					)
				}
			} else {
				if err != nil {
					t.Fatalf("want no err, but has error %v", err)
				}
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf(
						"uc.GetUsers(%v)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
						tt.in, got, tt.want, diff,
					)
				}
			}
		})
	}
}

func TestUserUsecase_UpdateUser(t *testing.T) {
	tests := []struct {
		name                string
		in                  *dto.UpdateUserInput
		wantUser            *model.User
		newMemoryRepository func() repository.Repository
		wantErr             error
	}{
		{
			name: "Update a user",
			in: &dto.UpdateUserInput{
				UserID: "TEST_USER_ID",
				Name:   "TEST_USER_NAME_UPDATED",
				Email:  "TEST_USER_EMAIL_UPDATED",
			},
			wantUser: model.MustNewUser(
				"TEST_USER_ID",
				"TEST_USER_NAME_UPDATED",
				"TEST_USER_EMAIL_UPDATED",
			),
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				s.AddUsers(model.MustNewUser("TEST_USER_ID", "TEST_USER_NAME", "TEST_USER_EMAIL"))
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
		{
			name: "Returns error if any of user inputs are invalid",
			in: &dto.UpdateUserInput{
				UserID: "TEST_USER_ID",
				Name:   "TEST_USER_NAME_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				Email:  "TEST_USER_EMAIL_UPDATED",
			},
			wantUser: nil,
			wantErr:  usecase.ErrInvalidUserInput,
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
		{
			name: "Returns error if the user does not exist",
			in: &dto.UpdateUserInput{
				UserID: "TEST_USER_ID",
				Name:   "TEST_USER_NAME_UPDATED",
				Email:  "TEST_USER_EMAIL_UPDATED",
			},
			wantUser: nil,
			wantErr:  usecase.ErrUserNotFound,
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

			f := mockfactory.NewMockUserFactory(ctrl)
			r := tt.newMemoryRepository()
			us := domainservice.NewUserService(r)
			gs := domainservice.NewGroupService(r)
			uc := usecase.NewUserUsecase(r, f, us, gs)

			_, err := uc.UpdateUser(tt.in)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("want an error, but has no error")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf(
						"uc.UpdateUser(%v)=_, %v; want _, %v",
						tt.in, err, tt.wantErr,
					)
				}
			} else {
				if err != nil {
					t.Fatalf("want no err, but has error %v", err)
				}
				uID := model.UserID(tt.in.UserID)
				got, _ := r.User().Find(uID)
				if diff := cmp.Diff(got, tt.wantUser, cmp.AllowUnexported(model.User{})); diff != "" {
					t.Errorf(
						"r.User().Find(%s)=%v, _; want %v, nil\ndiffers: (-got +want)\n%s",
						uID, got, tt.wantUser, diff,
					)
				}
			}
		})
	}
}

func TestUserUsecase_DeleteUser(t *testing.T) {
	tests := []struct {
		name                string
		in                  *dto.DeleteUserInput
		newMemoryRepository func() repository.Repository
		wantErr             error
	}{
		{
			name: "Delete a user",
			in: &dto.DeleteUserInput{
				UserID: "TEST_USER_ID",
			},
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				s.AddUsers(model.MustNewUser("TEST_USER_ID", "TEST_USER_NAME", "TEST_USER_EMAIL"))
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
		{
			name: "Delete a user and remove from all groups",
			in: &dto.DeleteUserInput{
				UserID: "TEST_USER_ID",
			},
			newMemoryRepository: func() repository.Repository {
				s := memory.NewStore()
				s.AddUsers(model.MustNewUser("TEST_USER_ID", "TEST_USER_NAME", "TEST_USER_EMAIL"))
				s.AddGroups(model.MustNewGroup(
					"TEST_GROUP_ID_1",
					"TEST_GROUP_NAME_1",
					[]model.UserID{
						"TEST_USER_ID_1",
						"TEST_USER_ID_2",
					}),
				)
				s.AddGroups(model.MustNewGroup(
					"TEST_GROUP_ID_2",
					"TEST_GROUP_NAME_2",
					[]model.UserID{
						"TEST_USER_ID_1",
						"TEST_USER_ID_3",
					}),
				)
				r := memory.NewMemoryRepository(s)
				return r
			},
		},
		{
			name: "Returns error if the user does not exist",
			in: &dto.DeleteUserInput{
				UserID: "TEST_USER_ID",
			},
			wantErr: usecase.ErrUserNotFound,
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

			f := mockfactory.NewMockUserFactory(ctrl)
			r := tt.newMemoryRepository()
			us := domainservice.NewUserService(r)
			gs := domainservice.NewGroupService(r)
			uc := usecase.NewUserUsecase(r, f, us, gs)

			_, err := uc.DeleteUser(tt.in)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("want an error, but has no error")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf(
						"uc.DeleteUser(%v)=_, %v; want _, %v",
						tt.in, err, tt.wantErr,
					)
				}
			} else {
				if err != nil {
					t.Fatalf("want no err, but has error %v", err)
				}

				uID := model.UserID(tt.in.UserID)
				gotUser, _ := r.User().Find(uID)
				if gotUser != nil {
					t.Errorf("r.User().Find(%s)=%v, _; want nil", uID, gotUser)
				}

				hasGroupTargetUser, _ := gs.HasUsersAny([]model.UserID{uID})
				if hasGroupTargetUser {
					t.Errorf("any of groups have the target user")
				}
			}
		})
	}
}
