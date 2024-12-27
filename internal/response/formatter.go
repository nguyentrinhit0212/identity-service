package response

import (
	"net/http"

	"identity-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, data interface{}) {
	meta := gin.H{
		"code":    http.StatusOK,
		"message": "Success",
	}

	response := gin.H{
		"data": data,
		"meta": meta,
	}

	json, err := utils.JSON.Marshal(response)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to encode response",
		})
		return
	}

	c.Data(http.StatusOK, "application/json", json)
}

func ErrorResponse(c *gin.Context, statusCode int, code string, message string, details interface{}) {
	response := gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
			"details": details,
		},
	}

	json, err := utils.JSON.Marshal(response)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to encode response",
		})
		return
	}

	c.Data(statusCode, "application/json", json)
}
