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
		UserID string `json:"user_id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
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
		UserID: out.UserID,
		Name:   out.Name,
		Email:  out.Email,
	})
}

type (
	GetUserResponse struct {
		UserID string `json:"user_id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
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
		UserID: out.UserID,
		Name:   out.Name,
		Email:  out.Email,
	})
}
