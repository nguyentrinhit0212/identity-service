package routes

import (
	"identity-service/internal/auth/jwt"
	"identity-service/internal/handlers"
	"identity-service/internal/middleware"
	"identity-service/internal/repositories"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, oauthHandler *handlers.OAuthHandler, keyManager *jwt.KeyManager, userRepo repositories.UserRepository) {
	// Public auth routes (no authentication required)
	authGroup := router.Group("/api/auth")
	{
		// OAuth routes
		providerGroup := authGroup.Group("/:provider")
		{
			providerGroup.GET("/login", oauthHandler.HandleLogin)       // Initiate OAuth login
			providerGroup.GET("/callback", oauthHandler.HandleCallback) // OAuth callback
		}

		// Token exchange endpoint for PKCE flow
		authGroup.POST("/token", oauthHandler.HandleTokenExchange)

		// Login
		authGroup.POST("/login", authHandler.Login) // User login with credentials

		// Session management
		authGroup.POST("/logout", authHandler.Logout)        // Logout
		authGroup.POST("/refresh", authHandler.RefreshToken) // Refresh access token
		authGroup.GET("/session", authHandler.GetSession)    // Get current session info

		// Password management
		authGroup.POST("/forgot-password", authHandler.ForgotPassword) // Request password reset
		authGroup.POST("/reset-password", authHandler.ResetPassword)   // Reset password with token
		authGroup.POST("/verify-email", authHandler.VerifyEmail)       // Verify email address
	}

	// Protected auth routes
	protectedGroup := authGroup.Group("")
	jwtMiddleware := middleware.NewJWTAuthMiddleware(keyManager, userRepo)
	protectedGroup.Use(jwtMiddleware.RequireAuth())
	{
		// Session management
		protectedGroup.GET("/sessions", authHandler.ListSessions)         // List all active sessions
		protectedGroup.DELETE("/sessions/:id", authHandler.RevokeSession) // Revoke specific session
		protectedGroup.DELETE("/sessions", authHandler.RevokeAllSessions) // Revoke all sessions except current

		// Security settings
		protectedGroup.GET("/security", authHandler.GetSecuritySettings)    // Get security settings
		protectedGroup.PUT("/security", authHandler.UpdateSecuritySettings) // Update security settings
		protectedGroup.POST("/mfa/enable", authHandler.EnableMFA)           // Enable MFA
		protectedGroup.POST("/mfa/disable", authHandler.DisableMFA)         // Disable MFA
		protectedGroup.POST("/mfa/verify", authHandler.VerifyMFA)           // Verify MFA token
	}
}
