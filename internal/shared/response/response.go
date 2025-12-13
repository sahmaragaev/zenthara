package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status       int    `json:"status"`
	IsSuccess    bool   `json:"success"`
	Message      string `json:"message,omitempty"`
	Data         any    `json:"data,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
	Details      any    `json:"details,omitempty"`
}

func Success(c *gin.Context, data any, message string) {
	resp := Response{
		Status:    http.StatusOK,
		IsSuccess: true,
		Data:      data,
		Message:   message,
	}

	c.JSON(resp.Status, resp)
}

func Error(c *gin.Context, status int, err error, details any) {
	resp := Response{
		Status:       status,
		IsSuccess:    false,
		ErrorMessage: err.Error(),
		Details:      details,
	}

	c.JSON(resp.Status, resp)
}

func BadRequest(c *gin.Context, err error, details any) {
	Error(c, http.StatusBadRequest, err, details)
}

func InternalServerError(c *gin.Context, err error, details any) {
	Error(c, http.StatusInternalServerError, err, details)
}

func Unauthorized(c *gin.Context, err error, details any) {
	Error(c, http.StatusUnauthorized, err, details)
}
