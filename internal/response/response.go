package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Success: true,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, message string, err error) {
	errorInfo := &ErrorInfo{
		Message: message,
	}
	if err != nil {
		errorInfo.Details = err.Error()
	}

	c.JSON(status, Response{
		Success: false,
		Error:   errorInfo,
	})
}
