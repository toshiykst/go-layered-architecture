package response

import (
	"net/http"

	"github.com/labstack/echo"
)

func OK(c echo.Context, body any) error {
	return c.JSON(http.StatusOK, body)
}

func Created(c echo.Context, body any) error {
	return c.JSON(http.StatusCreated, body)
}
