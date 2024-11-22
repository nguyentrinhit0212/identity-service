package routes

import (
	"identity-service/internal/auth"
	"identity-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, auth *auth.OAuthHandler) {
	authGroup := router.Group("/api/auth/:provider")
	{
		authGroup.GET("/login", auth.HandleLogin)
		authGroup.GET("/callback", auth.HandleCallback)
	}
	router.GET("/api/me", middleware.AuthMiddleware(), auth.GetMe)

}
