package handlers

import (
	"context"
	"fmt"
	"identity-service/internal/models"
	"identity-service/internal/services"
	"identity-service/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// OAuthHandler handles OAuth-related HTTP requests
type OAuthHandler struct {
	providers   map[string]services.OAuthProvider
	userService services.UserService
	authService services.AuthService
}

// NewOAuthHandler creates a new OAuth handler instance
func NewOAuthHandler(providers map[string]services.OAuthProvider, userService services.UserService, authService services.AuthService) *OAuthHandler {
	return &OAuthHandler{
		providers:   providers,
		userService: userService,
		authService: authService,
	}
}

// HandleLogin initiates OAuth login
func (h *OAuthHandler) HandleLogin(c *gin.Context) {
	providerName := c.Param("provider")
	provider, exists := h.providers[providerName]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider not supported"})
		return
	}

	state := utils.GenerateRandomState()
	url := provider.GetAuthURL(state)
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// HandleCallback handles OAuth callback
func (h *OAuthHandler) HandleCallback(c *gin.Context) {
	providerName := c.Param("provider")
	provider, exists := h.providers[providerName]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider not supported"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not provided"})
		return
	}

	token, err := provider.ExchangeToken(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	oauthUser, err := provider.FetchUserInfo(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}

	user, err := h.userService.CreateOrUpdateUser(oauthUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or update user"})
		return
	}

	// Get user's tenants
	tenants, err := h.userService.GetUserTenants(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user tenants"})
		return
	}

	// Find personal tenant
	var personalTenant *models.Tenant
	for _, tenant := range tenants {
		if tenant.Type == models.PersonalTenant {
			personalTenant = tenant
			break
		}
	}

	if personalTenant == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Personal tenant not found"})
		return
	}

	// Create a session for the user with their personal tenant
	session, err := h.authService.CreateSession(c, user, personalTenant.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	frontendURL := fmt.Sprintf("http://localhost:3000?accessToken=%s&refreshToken=%s", session.AccessToken, session.RefreshToken)
	c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}
