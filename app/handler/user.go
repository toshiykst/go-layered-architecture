package handler

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/toshiykst/go-layerd-architecture/app/usecase"
)

type userHandler struct {
	uc usecase.UserUsecase
}

// NewUserHandler returns a user handler struct.
func NewUserHandler(uc usecase.UserUsecase) *userHandler {
	return &userHandler{uc: uc}
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

func (h *userHandler) CreateUser(c echo.Context) error {
	req := &CreateUserRequest{}
	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return err
	}

	in := &usecase.CreateUserInput{
		Name:  req.Name,
		Email: req.Email,
	}
	out, err := h.uc.CreateUser(in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	return c.JSON(http.StatusCreated, &CreateUserResponse{
		UserID: out.UserID,
		Name:   out.Name,
		Email:  out.Email,
	})
}

func (h *userHandler) GetUser(c echo.Context) error {
	id := c.Param("id")

	in := &usecase.GetUserInput{
		UserID: id,
	}

	out, err := h.uc.GetUser(in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	return c.JSON(http.StatusOK, out)
}
