package middleware

import (
	"identity-service/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware xác thực JWT token trong header Authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Lấy token từ header (Bearer token)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Kiểm tra token hợp lệ
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Lưu thông tin người dùng vào context
		c.Set("user", claims)

		// Tiếp tục với request
		c.Next()
	}
}
