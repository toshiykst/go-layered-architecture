package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"

	"github.com/toshiykst/go-layerd-architecture/app/handler/response"
	"github.com/toshiykst/go-layerd-architecture/app/usecase"
	"github.com/toshiykst/go-layerd-architecture/app/usecase/dto"
)

type UserHandler struct {
	uc usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc: uc}
}

type (
	CreateUserRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	CreateUserResponse struct {
		User response.User `json:"user"`
	}
)

func (h *UserHandler) CreateUser(c echo.Context) error {
	req := &CreateUserRequest{}
	if err := c.Bind(req); err != nil {
		return response.Error(c, response.ErrorCodeInvalidArguments, http.StatusBadRequest, err)
	}

	in := &dto.CreateUserInput{
		Name:  req.Name,
		Email: req.Email,
	}
	out, err := h.uc.CreateUser(in)
	if err != nil {
		return response.ErrorInternal(c, err)
	}

	return response.Created(c, &CreateUserResponse{
		User: response.User{
			UserID: out.User.UserID,
			Name:   out.User.Name,
			Email:  out.User.Email,
		},
	})
}

type (
	GetUserResponse struct {
		User response.User `json:"user"`
	}
)

func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")
	in := &dto.GetUserInput{
		UserID: id,
	}

	out, err := h.uc.GetUser(in)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			return response.Error(c, response.ErrorCodeUserNotFound, http.StatusNotFound, err)
		} else {
			return response.ErrorInternal(c, err)
		}
	}

	return response.OK(c, &GetUserResponse{
		User: response.User{
			UserID: out.User.UserID,
			Name:   out.User.Name,
			Email:  out.User.Email,
		},
	})
}

type (
	GetUsersResponse struct {
		Users []response.User `json:"users"`
	}
)

func (h *UserHandler) GetUsers(c echo.Context) error {
	out, err := h.uc.GetUsers(nil)
	if err != nil {
		return response.ErrorInternal(c, err)
	}

	return response.OK(c, &GetUsersResponse{
		Users: response.ToUsersFromDTO(out.Users),
	})
}

type (
	UpdateUserRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)

func (h *UserHandler) UpdateUser(c echo.Context) error {
	uID := c.Param("id")

	req := &UpdateUserRequest{}
	if err := c.Bind(req); err != nil {
		return response.Error(c, response.ErrorCodeInvalidArguments, http.StatusBadRequest, err)
	}

	in := &dto.UpdateUserInput{
		UserID: uID,
		Name:   req.Name,
		Email:  req.Email,
	}

	_, err := h.uc.UpdateUser(in)
	if err != nil {
		return response.ErrorInternal(c, err)
	}

	return response.NoContent(c)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	uID := c.Param("id")

	in := &dto.DeleteUserInput{
		UserID: uID,
	}

	_, err := h.uc.DeleteUser(in)
	if err != nil {
		return response.ErrorInternal(c, err)
	}

	return response.NoContent(c)
}
