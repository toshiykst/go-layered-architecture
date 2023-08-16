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

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/labstack/echo"
	"go.uber.org/mock/gomock"

	"github.com/toshiykst/go-layerd-architecture/app/handler/response"
	mockusecase "github.com/toshiykst/go-layerd-architecture/app/mock/usecase"
	"github.com/toshiykst/go-layerd-architecture/app/usecase"
	"github.com/toshiykst/go-layerd-architecture/app/usecase/dto"
)

func TestNewGroupHandler(t *testing.T) {
	t.Run("Returns a group handler instance", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		uc := mockusecase.NewMockGroupUsecase(ctrl)
		got := NewGroupHandler(uc)
		want := &GroupHandler{uc: uc}
		if diff := cmp.Diff(got, want, cmpopts.IgnoreUnexported(GroupHandler{})); diff != "" {
			t.Errorf(
				"NewGroupHandler(%v)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
				uc, got, want, diff,
			)
		}
	})
}

func TestGroupHandler_CreateGroup(t *testing.T) {
	tests := []struct {
		name            string
		req             *CreateGroupRequest
		newGroupUsecase func(ctrl *gomock.Controller) usecase.GroupUsecase
		wantStatus      int
		wantRes         *CreateGroupResponse
		wantErrRes      *response.ErrorResponse
	}{
		{
			name: "Create a group and returns the group response with no users",
			req: &CreateGroupRequest{
				Name:    "TEST_GROUP_NAME",
				UserIDs: nil,
			},
			newGroupUsecase: func(ctrl *gomock.Controller) usecase.GroupUsecase {
				uc := mockusecase.NewMockGroupUsecase(ctrl)
				uc.EXPECT().
					CreateGroup(gomock.Any()).
					DoAndReturn(func(in *dto.CreateGroupInput) (*dto.CreateGroupOutput, error) {
						return &dto.CreateGroupOutput{
							Group: dto.Group{
								GroupID: "TEST_GROUP_ID",
								Name:    in.Name,
								Users:   []dto.User{},
							},
						}, nil
					})
				return uc
			},
			wantStatus: http.StatusCreated,
			wantRes: &CreateGroupResponse{
				Group: response.Group{
					GroupID: "TEST_GROUP_ID",
					Name:    "TEST_GROUP_NAME",
					Users:   []response.User{},
				},
			},
			wantErrRes: nil,
		},
		{
			name: "Create a group and returns the group response with users",
			req: &CreateGroupRequest{
				Name: "TEST_GROUP_NAME",
				UserIDs: []string{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_3",
				},
			},
			newGroupUsecase: func(ctrl *gomock.Controller) usecase.GroupUsecase {
				uc := mockusecase.NewMockGroupUsecase(ctrl)
				uc.EXPECT().
					CreateGroup(gomock.Any()).
					DoAndReturn(func(in *dto.CreateGroupInput) (*dto.CreateGroupOutput, error) {
						return &dto.CreateGroupOutput{
							Group: dto.Group{
								GroupID: "TEST_GROUP_ID",
								Name:    in.Name,
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
						}, nil
					})
				return uc
			},
			wantStatus: http.StatusCreated,
			wantRes: &CreateGroupResponse{
				Group: response.Group{
					GroupID: "TEST_GROUP_ID",
					Name:    "TEST_GROUP_NAME",
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
			},
			wantErrRes: nil,
		},
		{
			name: "Returns invalid arguments error response when some of user ids are invalid",
			req: &CreateGroupRequest{
				Name: "TEST_GROUP_NAME",
				UserIDs: []string{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_4",
				},
			},
			newGroupUsecase: func(ctrl *gomock.Controller) usecase.GroupUsecase {
				uc := mockusecase.NewMockGroupUsecase(ctrl)
				uc.EXPECT().
					CreateGroup(gomock.Any()).
					Return(nil, usecase.ErrInvalidUserIDs)
				return uc
			},
			wantStatus: http.StatusBadRequest,
			wantRes:    nil,
			wantErrRes: &response.ErrorResponse{
				Code:    response.ErrorCodeInvalidArguments,
				Status:  http.StatusBadRequest,
				Message: usecase.ErrInvalidUserIDs.Error(),
			},
		},
		{
			name: "Returns internal server error response",
			req: &CreateGroupRequest{
				Name:    "TEST_GROUP_NAME",
				UserIDs: nil,
			},
			newGroupUsecase: func(ctrl *gomock.Controller) usecase.GroupUsecase {
				uc := mockusecase.NewMockGroupUsecase(ctrl)
				uc.EXPECT().
					CreateGroup(gomock.Any()).
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
				"https://example.com:8080/groups",
				bytes.NewBuffer(reqJson),
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc := tt.newGroupUsecase(ctrl)

			h := NewGroupHandler(uc)

			err := h.CreateGroup(c)
			if err != nil {
				t.Fatalf("want no error, but has error: %v", err)
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
				var got *CreateGroupResponse
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

func TestGroupHandler_GetGroup(t *testing.T) {
	tests := []struct {
		name            string
		id              string
		newGroupUsecase func(ctrl *gomock.Controller) usecase.GroupUsecase
		wantStatus      int
		wantRes         *GetGroupResponse
		wantErrRes      *response.ErrorResponse
	}{
		{
			name: "Returns the group response",
			id:   "TEST_GROUP_ID",
			newGroupUsecase: func(ctrl *gomock.Controller) usecase.GroupUsecase {
				uc := mockusecase.NewMockGroupUsecase(ctrl)
				uc.EXPECT().
					GetGroup(gomock.Any()).
					DoAndReturn(func(in *dto.GetGroupInput) (*dto.GetGroupOutput, error) {
						return &dto.GetGroupOutput{
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
						}, nil
					})
				return uc
			},
			wantStatus: http.StatusOK,
			wantRes: &GetGroupResponse{
				Group: response.Group{
					GroupID: "TEST_GROUP_ID",
					Name:    "TEST_GROUP_NAME",
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
			},
			wantErrRes: nil,
		},
		{
			name: "Returns group not found error response",
			id:   "TEST_GROUP_ID",
			newGroupUsecase: func(ctrl *gomock.Controller) usecase.GroupUsecase {
				uc := mockusecase.NewMockGroupUsecase(ctrl)
				uc.EXPECT().
					GetGroup(gomock.Any()).
					Return(nil, usecase.ErrGroupNotFound)
				return uc
			},
			wantStatus: http.StatusNotFound,
			wantRes:    nil,
			wantErrRes: &response.ErrorResponse{
				Code:    response.ErrorCodeGroupNotFound,
				Status:  http.StatusNotFound,
				Message: usecase.ErrGroupNotFound.Error(),
			},
		},
		{
			name: "Returns internal server error response",
			id:   "TEST_GROUP_ID",
			newGroupUsecase: func(ctrl *gomock.Controller) usecase.GroupUsecase {
				uc := mockusecase.NewMockGroupUsecase(ctrl)
				uc.EXPECT().
					GetGroup(gomock.Any()).
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
				"https://example.com:8080/groups/"+tt.id,
				nil,
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/groups/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.id)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc := tt.newGroupUsecase(ctrl)

			h := NewGroupHandler(uc)

			err := h.GetGroup(c)
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
				var got *GetGroupResponse
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

func TestGroupHandler_GetGroups(t *testing.T) {
	tests := []struct {
		name            string
		newGroupUsecase func(ctrl *gomock.Controller) usecase.GroupUsecase
		wantStatus      int
		wantRes         *GetGroupsResponse
		wantErrRes      *response.ErrorResponse
	}{
		{
			name: "Returns groups response",
			newGroupUsecase: func(ctrl *gomock.Controller) usecase.GroupUsecase {
				uc := mockusecase.NewMockGroupUsecase(ctrl)
				uc.EXPECT().
					GetGroups(gomock.Any()).
					DoAndReturn(func(in *dto.GetGroupsInput) (*dto.GetGroupsOutput, error) {
						return &dto.GetGroupsOutput{
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
						}, nil
					})
				return uc
			},
			wantStatus: http.StatusOK,
			wantRes: &GetGroupsResponse{
				Groups: []response.Group{
					{
						GroupID: "TEST_GROUP_ID_1",
						Name:    "TEST_GROUP_NAME_1",
						Users: []response.User{
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
						},
					},
					{
						GroupID: "TEST_GROUP_ID_3",
						Name:    "TEST_GROUP_NAME_3",
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
				},
			},
			wantErrRes: nil,
		},
		{
			name: "Returns internal server error response",
			newGroupUsecase: func(ctrl *gomock.Controller) usecase.GroupUsecase {
				uc := mockusecase.NewMockGroupUsecase(ctrl)
				uc.EXPECT().
					GetGroups(gomock.Any()).
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
				"https://example.com:8080/groups",
				nil,
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/groups")

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc := tt.newGroupUsecase(ctrl)

			h := NewGroupHandler(uc)

			err := h.GetGroups(c)
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
				var got *GetGroupsResponse
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

func TestGroupHandler_UpdateGroup(t *testing.T) {
	tests := []struct {
		name            string
		gID             string
		req             *UpdateGroupRequest
		newGroupUsecase func(ctrl *gomock.Controller) usecase.GroupUsecase
		wantStatus      int
		wantErrRes      *response.ErrorResponse
	}{
		{
			name: "Update a group",
			gID:  "TEST_GROUP_ID",
			req: &UpdateGroupRequest{
				Name: "TEST_GROUP_NAME",
			},
			newGroupUsecase: func(ctrl *gomock.Controller) usecase.GroupUsecase {
				uc := mockusecase.NewMockGroupUsecase(ctrl)
				uc.EXPECT().
					UpdateGroup(gomock.Any()).
					DoAndReturn(func(in *dto.UpdateGroupInput) (*dto.UpdateGroupOutput, error) {
						return &dto.UpdateGroupOutput{}, nil
					})
				return uc
			},
			wantStatus: http.StatusNoContent,
			wantErrRes: nil,
		},
		{
			name: "Returns group not found error response",
			gID:  "TEST_GROUP_ID",
			newGroupUsecase: func(ctrl *gomock.Controller) usecase.GroupUsecase {
				uc := mockusecase.NewMockGroupUsecase(ctrl)
				uc.EXPECT().
					UpdateGroup(gomock.Any()).
					Return(nil, usecase.ErrGroupNotFound)
				return uc
			},
			wantStatus: http.StatusNotFound,
			wantErrRes: &response.ErrorResponse{
				Code:    response.ErrorCodeGroupNotFound,
				Status:  http.StatusNotFound,
				Message: usecase.ErrGroupNotFound.Error(),
			},
		},
		{
			name: "Returns internal server error response",
			gID:  "TEST_GROUP_ID",
			req: &UpdateGroupRequest{
				Name: "TEST_GROUP_NAME",
			},
			newGroupUsecase: func(ctrl *gomock.Controller) usecase.GroupUsecase {
				uc := mockusecase.NewMockGroupUsecase(ctrl)
				uc.EXPECT().
					UpdateGroup(gomock.Any()).
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
				fmt.Sprintf("https://example.com:8080/groups/%s", tt.gID),
				bytes.NewBuffer(reqJson),
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc := tt.newGroupUsecase(ctrl)

			h := NewGroupHandler(uc)

			err := h.UpdateGroup(c)
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

func TestGroupHandler_DeleteGroup(t *testing.T) {
	tests := []struct {
		name            string
		gID             string
		newGroupUsecase func(ctrl *gomock.Controller) usecase.GroupUsecase
		wantStatus      int
		wantErrRes      *response.ErrorResponse
	}{
		{
			name: "Delete a group",
			gID:  "TEST_GROUP_ID",
			newGroupUsecase: func(ctrl *gomock.Controller) usecase.GroupUsecase {
				uc := mockusecase.NewMockGroupUsecase(ctrl)
				uc.EXPECT().
					DeleteGroup(gomock.Any()).
					DoAndReturn(func(in *dto.DeleteGroupInput) (*dto.DeleteGroupOutput, error) {
						return &dto.DeleteGroupOutput{}, nil
					})
				return uc
			},
			wantStatus: http.StatusNoContent,
			wantErrRes: nil,
		},
		{
			name: "Returns group not found error response",
			gID:  "TEST_GROUP_ID",
			newGroupUsecase: func(ctrl *gomock.Controller) usecase.GroupUsecase {
				uc := mockusecase.NewMockGroupUsecase(ctrl)
				uc.EXPECT().
					DeleteGroup(gomock.Any()).
					Return(nil, usecase.ErrGroupNotFound)
				return uc
			},
			wantStatus: http.StatusNotFound,
			wantErrRes: &response.ErrorResponse{
				Code:    response.ErrorCodeGroupNotFound,
				Status:  http.StatusNotFound,
				Message: usecase.ErrGroupNotFound.Error(),
			},
		},
		{
			name: "Returns internal server error response",
			gID:  "TEST_GROUP_ID",
			newGroupUsecase: func(ctrl *gomock.Controller) usecase.GroupUsecase {
				uc := mockusecase.NewMockGroupUsecase(ctrl)
				uc.EXPECT().
					DeleteGroup(gomock.Any()).
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
				fmt.Sprintf("https://example.com:8080/groups/%s", tt.gID),
				nil,
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc := tt.newGroupUsecase(ctrl)

			h := NewGroupHandler(uc)

			err := h.DeleteGroup(c)
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
