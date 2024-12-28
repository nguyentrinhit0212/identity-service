package handlers

import (
	"context"
	"fmt"
	"identity-service/internal/models"
	"identity-service/internal/services"
	"identity-service/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// OAuthHandler handles OAuth-related HTTP requests
type OAuthHandler struct {
	providers   map[string]services.OAuthProvider
	userService services.UserService
	authService services.AuthService
	pkceService services.PKCEService
}

// NewOAuthHandler creates a new OAuth handler instance
func NewOAuthHandler(providers map[string]services.OAuthProvider, userService services.UserService, authService services.AuthService, pkceService services.PKCEService) *OAuthHandler {
	return &OAuthHandler{
		providers:   providers,
		userService: userService,
		authService: authService,
		pkceService: pkceService,
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
	tenants, err := h.userService.GetUserTenants(user.ID)
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

	// Create a PKCE challenge for secure token exchange
	challengeID := uuid.New()
	challenge := &models.PKCEChallenge{
		ID:            challengeID,
		CodeChallenge: utils.GenerateRandomString(32), // This will be verified against the code_verifier from frontend
		CodeVerifier:  utils.GenerateRandomString(32), // This will be sent to frontend
		UserID:        user.ID,
		TenantID:      personalTenant.ID,
		ExpiresAt:     time.Now().Add(5 * time.Minute),
		CreatedAt:     time.Now(),
	}

	if err := h.pkceService.CreateChallenge(challenge); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create challenge"})
		return
	}

	// Redirect to frontend with the code and code_verifier
	frontendURL := fmt.Sprintf("http://localhost:3000/%s/callback?code=%s&codeVerifier=%s",
		personalTenant.Slug,
		challengeID.String(),
		challenge.CodeVerifier,
	)
	c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}

// HandleTokenExchange handles the exchange of PKCE code for tokens
func (h *OAuthHandler) HandleTokenExchange(c *gin.Context) {
	var req models.TokenExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Parse the code (which is a UUID)
	challengeID, err := uuid.Parse(req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
		return
	}

	// Get and verify the challenge
	challenge, err := h.pkceService.GetChallenge(challengeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired code"})
		return
	}

	// Verify the code verifier
	if challenge.CodeVerifier != req.CodeVerifier {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code verifier"})
		return
	}

	// Get user
	user, err := h.userService.GetUser(challenge.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// Create session
	session, err := h.authService.CreateSession(c, user, challenge.TenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Mark challenge as used
	if err := h.pkceService.MarkChallengeAsUsed(challengeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark challenge as used"})
		return
	}

	// Return tokens
	c.JSON(http.StatusOK, models.TokenResponse{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
	})
}
