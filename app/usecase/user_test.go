package usecase

import (
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

func TestUserUsecase_CreateUser(t *testing.T) {
	tests := []struct {
		name              string
		in                *dto.CreateUserInput
		want              *dto.CreateUserOutput
		newMockRepository func() repository.Repository
		newMockFactory    func(ctrl *gomock.Controller) factory.UserFactory
		wantErr           error
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
			newMockRepository: func() repository.Repository {
				s := mockrepository.NewStore()
				r := mockrepository.NewMockRepository(s)
				return r
			},
			newMockFactory: func(ctrl *gomock.Controller) factory.UserFactory {
				f := mockfactory.NewMockUserFactory(ctrl)
				f.EXPECT().
					Create("TEST_USER_NAME", "TEST_USER_EMAIL").
					Return(model.NewUser("TEST_USER_ID", "TEST_USER_NAME", "TEST_USER_EMAIL"), nil)
				return f
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc := NewUserUsecase(
				tt.newMockRepository(),
				tt.newMockFactory(ctrl),
				mockdomainservice.NewMockUserService(ctrl),
			)
			got, err := uc.CreateUser(tt.in)
			if err != nil {
				t.Fatalf("want no err, but has error %v", err)
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf(
					"uc.CreateUser(%v)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
					tt.in, got, tt.want, diff,
				)
			}
		})
	}
}

func TestUserUsecase_GetUser(t *testing.T) {
	tests := []struct {
		name              string
		in                *dto.GetUserInput
		want              *dto.GetUserOutput
		wantErr           error
		newMockRepository func() repository.Repository
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
			newMockRepository: func() repository.Repository {
				s := mockrepository.NewStore()
				s.AddUsers(model.NewUser("TEST_USER_ID", "TEST_USER_NAME", "TEST_USER_EMAIL"))
				r := mockrepository.NewMockRepository(s)
				return r
			},
		},
		{
			name: "Returns error if the user does not exist",
			in: &dto.GetUserInput{
				UserID: "TEST_USER_ID",
			},
			want:    nil,
			wantErr: ErrUserNotFound,
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
			us := domainservice.NewUserService(r)
			uc := NewUserUsecase(
				r,
				mockfactory.NewMockUserFactory(ctrl),
				us,
			)
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
		name              string
		in                *dto.GetUsersInput
		want              *dto.GetUsersOutput
		wantErr           error
		newMockRepository func() repository.Repository
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			r := tt.newMockRepository()
			us := domainservice.NewUserService(r)
			uc := NewUserUsecase(
				r,
				mockfactory.NewMockUserFactory(ctrl),
				us,
			)
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
