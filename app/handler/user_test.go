package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	"github.com/toshiykst/go-layerd-architecture/app/usecase/dto"
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
					DoAndReturn(func(in *dto.CreateUserInput) (*dto.CreateUserOutput, error) {
						return &dto.CreateUserOutput{
							User: dto.User{
								UserID: "TEST_USER_ID",
								Name:   in.Name,
								Email:  in.Email,
							},
						}, nil
					})
				return uc
			},
			wantStatus: http.StatusCreated,
			wantRes: &CreateUserResponse{
				User: response.User{
					UserID: "TEST_USER_ID",
					Name:   "TEST_USER_NAME",
					Email:  "TEST_USER_EMAIL",
				},
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
					DoAndReturn(func(in *dto.GetUserInput) (*dto.GetUserOutput, error) {
						return &dto.GetUserOutput{
							User: dto.User{
								UserID: in.UserID,
								Name:   "TEST_USER_NAME",
								Email:  "TEST_USER_EMAIL",
							},
						}, nil
					})
				return uc
			},
			wantStatus: http.StatusOK,
			wantRes: &GetUserResponse{
				User: response.User{
					UserID: "TEST_USER_ID",
					Name:   "TEST_USER_NAME",
					Email:  "TEST_USER_EMAIL",
				},
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

func TestUserHandler_GetUsers(t *testing.T) {
	tests := []struct {
		name           string
		newUserUsecase func(ctrl *gomock.Controller) usecase.UserUsecase
		wantStatus     int
		wantRes        *GetUsersResponse
		wantErrRes     *response.ErrorResponse
	}{
		{
			name: "Returns users response",
			newUserUsecase: func(ctrl *gomock.Controller) usecase.UserUsecase {
				uc := mockusecase.NewMockUserUsecase(ctrl)
				uc.EXPECT().
					GetUsers(gomock.Any()).
					Return(
						&dto.GetUsersOutput{
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
						}, nil)

				return uc
			},
			wantStatus: http.StatusOK,
			wantRes: &GetUsersResponse{
				Users: []response.User{
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
			wantErrRes: nil,
		},
		{
			name: "Returns internal server error response",
			newUserUsecase: func(ctrl *gomock.Controller) usecase.UserUsecase {
				uc := mockusecase.NewMockUserUsecase(ctrl)
				uc.EXPECT().
					GetUsers(gomock.Any()).
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
				"https://example.com:8080/users",
				nil,
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/users")

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc := tt.newUserUsecase(ctrl)

			h := NewUserHandler(uc)

			err := h.GetUsers(c)
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
				var got *GetUsersResponse
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

func TestUserHandler_UpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		req            *UpdateUserRequest
		newUserUsecase func(ctrl *gomock.Controller) usecase.UserUsecase
		wantStatus     int
		wantErrRes     *response.ErrorResponse
	}{
		{
			name: "Update a user",
			req: &UpdateUserRequest{
				UserID: "TEST_USER_ID",
				Name:   "TEST_USER_NAME",
				Email:  "TEST_USER_EMAIL",
			},
			newUserUsecase: func(ctrl *gomock.Controller) usecase.UserUsecase {
				uc := mockusecase.NewMockUserUsecase(ctrl)
				uc.EXPECT().
					UpdateUser(gomock.Any()).
					DoAndReturn(func(in *dto.UpdateUserInput) (*dto.UpdateUserOutput, error) {
						return &dto.UpdateUserOutput{}, nil
					})
				return uc
			},
			wantStatus: http.StatusNoContent,
			wantErrRes: nil,
		},
		{
			name: "Returns internal server error response",
			req: &UpdateUserRequest{
				UserID: "TEST_USER_ID",
				Name:   "TEST_USER_NAME",
				Email:  "TEST_USER_EMAIL",
			},
			newUserUsecase: func(ctrl *gomock.Controller) usecase.UserUsecase {
				uc := mockusecase.NewMockUserUsecase(ctrl)
				uc.EXPECT().
					UpdateUser(gomock.Any()).
					Return(nil, errors.New("an error occurred"))
				return uc
			},
			wantStatus: http.StatusInternalServerError,
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
				http.MethodPut,
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

			err := h.UpdateUser(c)
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

func TestUserHandler_DeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		uID            string
		newUserUsecase func(ctrl *gomock.Controller) usecase.UserUsecase
		wantStatus     int
		wantErrRes     *response.ErrorResponse
	}{
		{
			name: "Delete a user",
			uID:  "TEST_USER_ID",
			newUserUsecase: func(ctrl *gomock.Controller) usecase.UserUsecase {
				uc := mockusecase.NewMockUserUsecase(ctrl)
				uc.EXPECT().
					DeleteUser(gomock.Any()).
					DoAndReturn(func(in *dto.DeleteUserInput) (*dto.DeleteUserOutput, error) {
						return &dto.DeleteUserOutput{}, nil
					})
				return uc
			},
			wantStatus: http.StatusNoContent,
			wantErrRes: nil,
		},
		{
			name: "Returns internal server error response",
			uID:  "TEST_USER_ID",
			newUserUsecase: func(ctrl *gomock.Controller) usecase.UserUsecase {
				uc := mockusecase.NewMockUserUsecase(ctrl)
				uc.EXPECT().
					DeleteUser(gomock.Any()).
					Return(nil, errors.New("an error occurred"))
				return uc
			},
			wantStatus: http.StatusInternalServerError,
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
				http.MethodDelete,
				fmt.Sprintf("https://example.com:8080/users/%s", "TEST_USER_ID"),
				nil,
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc := tt.newUserUsecase(ctrl)

			h := NewUserHandler(uc)

			err := h.DeleteUser(c)
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
