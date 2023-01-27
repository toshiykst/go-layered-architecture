package handler

import (
	"github.com/labstack/echo"
	"github.com/toshiykst/go-layerd-architecture/app/usecase"
	"net/http"
)

type userHandler struct {
	uc usecase.UserUsecase
}

// NewUserHandler returns a user handler struct.
func NewUserHandler(uc usecase.UserUsecase) *userHandler {
	return &userHandler{uc: uc}
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
