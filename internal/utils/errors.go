package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"error"`
}

type ErrorCode int

var (
	ErrorCodeBadRequest    ErrorCode = 400
	ErrorCodeUnauthorized  ErrorCode = 401
	ErrorCodeForbidden     ErrorCode = 403
	ErrorCodeNotFound      ErrorCode = 404
	ErrorCodeInternalError ErrorCode = 500
)

func ErrCodeToInt(err ErrorCode) int {
	switch err {
	case ErrorCodeBadRequest:
		return http.StatusBadRequest
	case ErrorCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrorCodeForbidden:
		return http.StatusForbidden
	case ErrorCodeNotFound:
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}

func (app *AppError) Error() string {
	return app.Message
}

func NewError(code ErrorCode, msg string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg}
}

func WrapError(err error, message string, code ErrorCode) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err}
}

func ResponseError(c *gin.Context, err error) {
	// vì AppError là error nên ép kiểu sang AppError
	if err, ok := err.(*AppError); ok {
		c.JSON(ErrCodeToInt(err.Code), gin.H{
			"code":    err.Code,
			"message": err.Message,
		})
		return
	}
	// Trường hợp err là error bình thường
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    ErrorCodeInternalError,
		"message": "Internal server error",
	})
}

func ResponseSuccess(c *gin.Context, statusCode int, data interface{}, message string) {
	if statusCode == http.StatusNoContent {
		c.Status(statusCode)
		return
	}

	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}
