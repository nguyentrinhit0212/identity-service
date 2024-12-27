package routes

import (
	"identity-service/internal/auth/jwt"
	"identity-service/internal/handlers"
	"identity-service/internal/middleware"
	"identity-service/internal/repositories"

	"github.com/gin-gonic/gin"
)

func SecurityRoutes(router *gin.Engine, handler *handlers.SecurityHandler, keyManager *jwt.KeyManager, userRepo repositories.UserRepository) {
	// All security routes require authentication
	securityGroup := router.Group("/api/security")
	jwtMiddleware := middleware.NewJWTAuthMiddleware(keyManager, userRepo)
	securityGroup.Use(jwtMiddleware.RequireAuth())
	{
		// IP Whitelist management
		whitelistGroup := securityGroup.Group("/whitelist")
		{
			whitelistGroup.GET("", handler.ListWhitelistedIPs)         // List whitelisted IPs
			whitelistGroup.POST("", handler.AddWhitelistedIP)          // Add IP to whitelist
			whitelistGroup.DELETE("/:id", handler.RemoveWhitelistedIP) // Remove IP from whitelist
		}

		// API Keys management
		apiKeyGroup := securityGroup.Group("/api-keys")
		{
			apiKeyGroup.GET("", handler.ListAPIKeys)         // List API keys
			apiKeyGroup.POST("", handler.CreateAPIKey)       // Create new API key
			apiKeyGroup.GET("/:id", handler.GetAPIKey)       // Get API key details
			apiKeyGroup.PUT("/:id", handler.UpdateAPIKey)    // Update API key
			apiKeyGroup.DELETE("/:id", handler.DeleteAPIKey) // Delete API key
		}

		// Audit logs
		securityGroup.GET("/audit-logs", handler.ListAuditLogs)        // List audit logs
		securityGroup.GET("/audit-logs/:id", handler.GetAuditLogEntry) // Get audit log entry

		// Security policies
		securityGroup.GET("/policies", handler.ListSecurityPolicies)     // List security policies
		securityGroup.PUT("/policies", handler.UpdateSecurityPolicies)   // Update security policies
		securityGroup.POST("/policies/test", handler.TestSecurityPolicy) // Test security policy
	}
}
