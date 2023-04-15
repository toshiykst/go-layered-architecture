package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/labstack/echo"

	"github.com/toshiykst/go-layerd-architecture/app/handler/response"
	mockusecase "github.com/toshiykst/go-layerd-architecture/app/mock/usecase"
	"github.com/toshiykst/go-layerd-architecture/app/usecase"
)

func TestNewUserHandler(t *testing.T) {
	t.Run("Returns an user handler instance", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		uc := mockusecase.NewMockUserUsecase(ctrl)
		got := NewUserHandler(uc)
		want := &UserHandler{uc: uc}
		if diff := cmp.Diff(got, want, cmpopts.IgnoreUnexported(UserHandler{})); diff != "" {
			t.Errorf(
				"NewUserHandler(%v)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
				uc, got, want, diff,
			)
		}
	})
}

func TestUserHandler_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		req            *CreateUserRequest
		newUserUsecase func(ctrl *gomock.Controller) usecase.UserUsecase
		wantStatus     int
		wantRes        *CreateUserResponse
		wantErrRes     *response.ErrorResponse
	}{
		{
			name: "Create a user and returns the user response",
			req: &CreateUserRequest{
				Name:  "TEST_USER_NAME",
				Email: "TEST_USER_EMAIL",
			},
			newUserUsecase: func(ctrl *gomock.Controller) usecase.UserUsecase {
				uc := mockusecase.NewMockUserUsecase(ctrl)
				uc.EXPECT().
					CreateUser(gomock.Any()).
					DoAndReturn(func(in *usecase.CreateUserInput) (*usecase.CreateUserOutput, error) {
						return &usecase.CreateUserOutput{
							UserID: "TEST_USER_ID",
							Name:   in.Name,
							Email:  in.Email,
						}, nil
					})
				return uc
			},
			wantStatus: http.StatusCreated,
			wantRes: &CreateUserResponse{
				UserID: "TEST_USER_ID",
				Name:   "TEST_USER_NAME",
				Email:  "TEST_USER_EMAIL",
			},
			wantErrRes: nil,
		},
		{
			name: "Returns internal server error response",
			req: &CreateUserRequest{
				Name:  "TEST_USER_NAME",
				Email: "TEST_USER_EMAIL",
			},
			newUserUsecase: func(ctrl *gomock.Controller) usecase.UserUsecase {
				uc := mockusecase.NewMockUserUsecase(ctrl)
				uc.EXPECT().
					CreateUser(gomock.Any()).
					Return(nil, errors.New("an error occurred"))
				return uc
			},
			wantStatus: http.StatusInternalServerError,
			wantRes:    nil,
			wantErrRes: &response.ErrorResponse{
				Code:    response.ErrorCodeInternalServerError,
				Status:  http.StatusInternalServerError,
				Message: "an error occurred",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqJson, _ := json.Marshal(tt.req)

			req := httptest.NewRequest(
				http.MethodPost,
				"https://example.com:8080/users",
				bytes.NewBuffer(reqJson),
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc := tt.newUserUsecase(ctrl)

			h := NewUserHandler(uc)

			err := h.CreateUser(c)
			if err != nil {
				t.Fatalf("want no err, but has error: %v", err)
			}

			res := rec.Result()

			wantStatusCode := tt.wantStatus
			gotStatusCode := res.StatusCode
			if gotStatusCode != wantStatusCode {
				t.Errorf("statusCode got = %d, want = %d", gotStatusCode, wantStatusCode)
			}

			resBody, err := io.ReadAll(rec.Body)
			if err != nil {
				t.Fatalf("Failed to read body: %s", err.Error())
			}
			defer func(Body io.ReadCloser) {
				if err := Body.Close(); err != nil {
					t.Fatalf("Failed to close body: %s", err.Error())
				}
			}(res.Body)

			if tt.wantRes != nil {
				var got *CreateUserResponse
				_ = json.Unmarshal(resBody, &got)
				if diff := cmp.Diff(got, tt.wantRes); diff != "" {
					t.Errorf(
						"response body: got = %v, want = %v\ndiffers: (-got +want)\n%s",
						got, tt.wantRes, diff,
					)
				}
			}

			if tt.wantErrRes != nil {
				var got *response.ErrorResponse
				_ = json.Unmarshal(resBody, &got)
				if diff := cmp.Diff(got, tt.wantErrRes); diff != "" {
					t.Errorf(
						"error response body: got = %v, want = %v\ndiffers: (-got +want)\n%s",
						got, tt.wantErrRes, diff,
					)
				}
			}
		})
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		newUserUsecase func(ctrl *gomock.Controller) usecase.UserUsecase
		wantStatus     int
		wantRes        *GetUserResponse
		wantErrRes     *response.ErrorResponse
	}{
		{
			name: "Returns the user response",
			id:   "TEST_USER_ID",
			newUserUsecase: func(ctrl *gomock.Controller) usecase.UserUsecase {
				uc := mockusecase.NewMockUserUsecase(ctrl)
				uc.EXPECT().
					GetUser(gomock.Any()).
					DoAndReturn(func(in *usecase.GetUserInput) (*usecase.GetUserOutput, error) {
						return &usecase.GetUserOutput{
							UserID: in.UserID,
							Name:   "TEST_USER_NAME",
							Email:  "TEST_USER_EMAIL",
						}, nil
					})
				return uc
			},
			wantStatus: http.StatusOK,
			wantRes: &GetUserResponse{
				UserID: "TEST_USER_ID",
				Name:   "TEST_USER_NAME",
				Email:  "TEST_USER_EMAIL",
			},
			wantErrRes: nil,
		},
		{
			name: "Returns user not found error response",
			id:   "TEST_USER_ID",
			newUserUsecase: func(ctrl *gomock.Controller) usecase.UserUsecase {
				uc := mockusecase.NewMockUserUsecase(ctrl)
				uc.EXPECT().
					GetUser(gomock.Any()).
					Return(nil, usecase.ErrUserNotFound)
				return uc
			},
			wantStatus: http.StatusNotFound,
			wantRes:    nil,
			wantErrRes: &response.ErrorResponse{
				Code:    response.ErrorCodeUserNotFound,
				Status:  http.StatusNotFound,
				Message: usecase.ErrUserNotFound.Error(),
			},
		},
		{
			name: "Returns internal server error response",
			id:   "TEST_USER_ID",
			newUserUsecase: func(ctrl *gomock.Controller) usecase.UserUsecase {
				uc := mockusecase.NewMockUserUsecase(ctrl)
				uc.EXPECT().
					GetUser(gomock.Any()).
					Return(nil, errors.New("an error occurred"))
				return uc
			},
			wantStatus: http.StatusInternalServerError,
			wantRes:    nil,
			wantErrRes: &response.ErrorResponse{
				Code:    response.ErrorCodeInternalServerError,
				Status:  http.StatusInternalServerError,
				Message: "an error occurred",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				"https://example.com:8080/users/"+tt.id,
				nil,
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/users/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.id)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc := tt.newUserUsecase(ctrl)

			h := NewUserHandler(uc)

			err := h.GetUser(c)
			if err != nil {
				t.Fatalf("want no err, but has error: %v", err)
			}

			res := rec.Result()

			wantStatusCode := tt.wantStatus
			gotStatusCode := res.StatusCode
			if gotStatusCode != wantStatusCode {
				t.Errorf("statusCode got = %d, want = %d", gotStatusCode, wantStatusCode)
			}

			resBody, err := io.ReadAll(rec.Body)
			if err != nil {
				t.Fatalf("Failed to read body: %s", err.Error())
			}
			defer func(Body io.ReadCloser) {
				if err := Body.Close(); err != nil {
					t.Fatalf("Failed to close body: %s", err.Error())
				}
			}(res.Body)

			if tt.wantRes != nil {
				var got *GetUserResponse
				_ = json.Unmarshal(resBody, &got)
				if diff := cmp.Diff(got, tt.wantRes); diff != "" {
					t.Errorf(
						"response body: got = %v, want = %v\ndiffers: (-got +want)\n%s",
						got, tt.wantRes, diff,
					)
				}
			}

			if tt.wantErrRes != nil {
				var got *response.ErrorResponse
				_ = json.Unmarshal(resBody, &got)
				if diff := cmp.Diff(got, tt.wantErrRes); diff != "" {
					t.Errorf(
						"error response body: got = %v, want = %v\ndiffers: (-got +want)\n%s",
						got, tt.wantErrRes, diff,
					)
				}
			}
		})
	}
}
