package auth

import (
	"context"
	"fmt"
	"identity-service/internal/models"
	"identity-service/internal/response"
	"identity-service/internal/services"
	"identity-service/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OAuthHandler struct {
    providers   map[string]OAuthProviderInterface
    userService services.UserService
}

func NewOAuthHandler(providers map[string]OAuthProviderInterface, userService services.UserService) *OAuthHandler {
    return &OAuthHandler{
        providers:   providers,
        userService: userService,
    }
}

// HandleLogin xử lý redirect đến provider
func (h *OAuthHandler)HandleLogin(c *gin.Context) {
	providerName := c.Param("provider")
	provider, exists := h.providers[providerName]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider not supported"})
		return
	}

	state := "random_state" // Tạo state ngẫu nhiên
	url := provider.GetAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleCallback xử lý callback từ provider
func (h *OAuthHandler)HandleCallback(c *gin.Context) {
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

	// Đổi code lấy access token
	token, err := provider.ExchangeToken(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Lấy thông tin user từ provider
	oauthUser, err := provider.FetchUserInfo(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}

	// Tạo hoặc cập nhật user trong database
	user, err := h.userService.CreateOrUpdateUser(oauthUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or update user"})
		return
	}

	// Tạo JWT token cho user
	jwtToken, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Redirect về frontend
	frontendURL := fmt.Sprintf("http://localhost:3000?token=%s", jwtToken)
	c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}

// GetMe trả về thông tin người dùng hiện tại từ JWT token
func (h *OAuthHandler) GetMe(c *gin.Context) {
    // Lấy thông tin người dùng từ context (được set trong middleware AuthMiddleware)
    userClaims, _ := c.Get("user")
    claims := userClaims.(*utils.Claims)

    // Kiểm tra và parse UUID
	userID := claims.UserID

	// Tạo một đối tượng người dùng từ claims
	user := models.User{
		ID:    userID,
		Email: claims.Email,
	}

	response.SuccessResponse(c, user)
}