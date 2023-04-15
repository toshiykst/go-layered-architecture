package response

import (
	"net/http"

	"github.com/labstack/echo"
)

type ErrorCode string

const (
	ErrorCodeInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrorCodeInvalidArguments    ErrorCode = "INVALID_INVALID_ARGUMENTS"
	ErrorCodeUserNotFound        ErrorCode = "USER_NOT_FOUND"
)

type ErrorResponse struct {
	Code    ErrorCode `json:"code"`
	Status  int       `json:"status"`
	Message string    `json:"message"`
}

func Error(c echo.Context, code ErrorCode, status int, err error) error {
	return c.JSON(status, &ErrorResponse{
		Code:    code,
		Status:  status,
		Message: err.Error(),
	})
}

func ErrorInternal(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, &ErrorResponse{
		Code:    ErrorCodeInternalServerError,
		Status:  http.StatusInternalServerError,
		Message: err.Error(),
	})
}
