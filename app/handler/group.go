package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"

	"github.com/toshiykst/go-layerd-architecture/app/handler/response"
	"github.com/toshiykst/go-layerd-architecture/app/usecase"
	"github.com/toshiykst/go-layerd-architecture/app/usecase/dto"
)

type GroupHandler struct {
	uc usecase.GroupUsecase
}

func NewGroupHandler(uc usecase.GroupUsecase) *GroupHandler {
	return &GroupHandler{uc: uc}
}

type (
	CreateGroupRequest struct {
		Name    string   `json:"name"`
		UserIDs []string `json:"userIds"`
	}

	CreateGroupResponse struct {
		Group response.Group `json:"group"`
	}
)

func (h *GroupHandler) CreateGroup(c echo.Context) error {
	req := &CreateGroupRequest{}
	if err := c.Bind(req); err != nil {
		return response.Error(c, response.ErrorCodeInvalidArguments, http.StatusBadRequest, err)
	}

	in := &dto.CreateGroupInput{
		Name:    req.Name,
		UserIDs: req.UserIDs,
	}
	out, err := h.uc.CreateGroup(in)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidUserIDs) {
			return response.Error(c, response.ErrorCodeInvalidArguments, http.StatusBadRequest, err)
		} else {
			return response.ErrorInternal(c, err)
		}
	}

	us := make([]response.User, len(out.Group.Users))
	for i, u := range out.Group.Users {
		us[i] = response.User{
			UserID: u.UserID,
			Name:   u.Name,
			Email:  u.Email,
		}
	}

	return response.Created(c, &CreateGroupResponse{
		Group: response.Group{
			GroupID: out.Group.GroupID,
			Name:    out.Group.Name,
			Users:   us,
		},
	})
}

type GetGroupResponse struct {
	Group response.Group `json:"group"`
}

func (h *GroupHandler) GetGroup(c echo.Context) error {
	id := c.Param("id")
	in := &dto.GetGroupInput{
		GroupID: id,
	}

	out, err := h.uc.GetGroup(in)
	if err != nil {
		if errors.Is(err, usecase.ErrGroupNotFound) {
			return response.Error(c, response.ErrorCodeGroupNotFound, http.StatusNotFound, err)
		} else {
			return response.ErrorInternal(c, err)
		}
	}

	us := make([]response.User, len(out.Group.Users))
	for i, u := range out.Group.Users {
		us[i] = response.User{
			UserID: u.UserID,
			Name:   u.Name,
			Email:  u.Email,
		}
	}

	return response.OK(c, &GetGroupResponse{
		Group: response.Group{
			GroupID: out.Group.GroupID,
			Name:    out.Group.Name,
			Users:   us,
		},
	})
}
