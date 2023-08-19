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
		if errors.Is(err, usecase.ErrInvalidGroupInput) {
			return response.Error(c, response.ErrorCodeInvalidArguments, http.StatusBadRequest, err)
		}
		if errors.Is(err, usecase.ErrInvalidUserIDs) {
			return response.Error(c, response.ErrorCodeInvalidArguments, http.StatusBadRequest, err)
		}
		return response.ErrorInternal(c, err)
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

	return response.OK(c, &GetGroupResponse{
		Group: response.Group{
			GroupID: out.Group.GroupID,
			Name:    out.Group.Name,
			Users:   response.ToUsersFromDTO(out.Group.Users),
		},
	})
}

type GetGroupsResponse struct {
	Groups []response.Group `json:"groups"`
}

func (h *GroupHandler) GetGroups(c echo.Context) error {
	out, err := h.uc.GetGroups(nil)
	if err != nil {
		return response.ErrorInternal(c, err)
	}

	gs := make([]response.Group, len(out.Groups))
	for i, g := range out.Groups {
		gs[i] = response.Group{
			GroupID: g.GroupID,
			Name:    g.Name,
			Users:   response.ToUsersFromDTO(g.Users),
		}
	}

	return response.OK(c, &GetGroupsResponse{
		Groups: gs,
	})
}

type (
	UpdateGroupRequest struct {
		Name string `json:"name"`
	}
)

func (h *GroupHandler) UpdateGroup(c echo.Context) error {
	gID := c.Param("id")

	req := &UpdateGroupRequest{}
	if err := c.Bind(req); err != nil {
		return response.Error(c, response.ErrorCodeInvalidArguments, http.StatusBadRequest, err)
	}

	in := &dto.UpdateGroupInput{
		GroupID: gID,
		Name:    req.Name,
	}

	_, err := h.uc.UpdateGroup(in)
	if err != nil {
		if errors.Is(err, usecase.ErrGroupNotFound) {
			return response.Error(c, response.ErrorCodeGroupNotFound, http.StatusNotFound, err)
		}
		if errors.Is(err, usecase.ErrInvalidGroupInput) {
			return response.Error(c, response.ErrorCodeInvalidArguments, http.StatusBadRequest, err)
		}
		return response.ErrorInternal(c, err)

	}

	return response.NoContent(c)
}

func (h *GroupHandler) DeleteGroup(c echo.Context) error {
	gID := c.Param("id")

	in := &dto.DeleteGroupInput{
		GroupID: gID,
	}

	_, err := h.uc.DeleteGroup(in)
	if err != nil {
		if errors.Is(err, usecase.ErrGroupNotFound) {
			return response.Error(c, response.ErrorCodeGroupNotFound, http.StatusNotFound, err)
		} else {
			return response.ErrorInternal(c, err)
		}
	}

	return response.NoContent(c)
}
