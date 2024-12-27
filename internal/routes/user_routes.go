package routes

import (
	"identity-service/internal/auth/jwt"
	"identity-service/internal/handlers"
	"identity-service/internal/middleware"
	"identity-service/internal/repositories"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, handler *handlers.UserHandler, keyManager *jwt.KeyManager, userRepo repositories.UserRepository) {
	// All user routes require authentication
	userGroup := router.Group("/api/users")
	jwtMiddleware := middleware.NewJWTAuthMiddleware(keyManager, userRepo)
	userGroup.Use(jwtMiddleware.RequireAuth())
	{
		// User management
		userGroup.GET("", handler.ListUsers)         // List users (with pagination and filters)
		userGroup.POST("", handler.CreateUser)       // Create new user
		userGroup.GET("/:id", handler.GetUser)       // Get user details
		userGroup.PUT("/:id", handler.UpdateUser)    // Update user
		userGroup.DELETE("/:id", handler.DeleteUser) // Delete user

		// User-tenant relationships
		userGroup.GET("/:id/tenants", handler.ListUserTenants)                   // List user's tenants
		userGroup.POST("/:id/tenants", handler.AddUserToTenant)                  // Add user to tenant
		userGroup.DELETE("/:id/tenants/:tenantId", handler.RemoveUserFromTenant) // Remove user from tenant

		// User profile
		userGroup.GET("/me", handler.GetProfile)              // Get own profile
		userGroup.PUT("/me", handler.UpdateProfile)           // Update own profile
		userGroup.PUT("/me/password", handler.UpdatePassword) // Update password
	}
}
