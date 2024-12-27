package routes

import (
	"identity-service/internal/auth/jwt"
	"identity-service/internal/handlers"
	"identity-service/internal/middleware"
	"identity-service/internal/repositories"

	"github.com/gin-gonic/gin"
)

func TenantRoutes(router *gin.Engine, handler *handlers.TenantHandler, keyManager *jwt.KeyManager, userRepo repositories.UserRepository) {
	// All tenant routes require authentication
	tenantGroup := router.Group("/api/tenants")
	jwtMiddleware := middleware.NewJWTAuthMiddleware(keyManager, userRepo)
	tenantGroup.Use(jwtMiddleware.RequireAuth())
	{
		// Tenant management
		tenantGroup.GET("", handler.ListTenants)         // List tenants (with pagination and filters)
		tenantGroup.POST("", handler.CreateTenant)       // Create new tenant
		tenantGroup.GET("/:id", handler.GetTenant)       // Get tenant details
		tenantGroup.PUT("/:id", handler.UpdateTenant)    // Update tenant
		tenantGroup.DELETE("/:id", handler.DeleteTenant) // Delete tenant

		// Tenant operations
		tenantGroup.POST("/:id/switch", handler.SwitchTenant)   // Switch active tenant
		tenantGroup.POST("/:id/upgrade", handler.UpgradeTenant) // Upgrade tenant plan

		// Tenant settings
		tenantGroup.GET("/:id/settings", handler.GetTenantSettings)    // Get tenant settings
		tenantGroup.PUT("/:id/settings", handler.UpdateTenantSettings) // Update tenant settings

		// Tenant members
		tenantGroup.GET("/:id/members", handler.ListTenantMembers)               // List tenant members
		tenantGroup.POST("/:id/invites", handler.CreateTenantInvite)             // Create tenant invite
		tenantGroup.DELETE("/:id/invites/:inviteId", handler.DeleteTenantInvite) // Delete tenant invite

		// Tenant features
		tenantGroup.GET("/:id/features", handler.ListTenantFeatures)   // List tenant features
		tenantGroup.PUT("/:id/features", handler.UpdateTenantFeatures) // Update tenant features
	}
}
