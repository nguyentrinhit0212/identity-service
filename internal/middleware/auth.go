package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"identity-service/internal/auth"
	jwtmanager "identity-service/internal/auth/jwt"
	"identity-service/internal/models"
	"identity-service/internal/response"
)

type JWTAuthMiddleware struct {
	keyManager *jwtmanager.KeyManager
	userRepo   UserRepository
}

type UserRepository interface {
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserTenantAccess(userID uuid.UUID) ([]models.UserTenantAccess, error)
	GetTenantByID(id uuid.UUID) (*models.Tenant, error)
}

func NewJWTAuthMiddleware(keyManager *jwtmanager.KeyManager, userRepo UserRepository) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		keyManager: keyManager,
		userRepo:   userRepo,
	}
}

func (m *JWTAuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "No authorization header", nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			response.Error(c, http.StatusUnauthorized, "Invalid authorization header format", nil)
			c.Abort()
			return
		}

		// Parse token without validation to get the key ID
		parser := jwt.Parser{
			ValidMethods: []string{jwt.SigningMethodRS256.Name},
		}
		token, _ := parser.Parse(tokenString, nil)
		if token == nil {
			response.Error(c, http.StatusUnauthorized, "Invalid token format", nil)
			c.Abort()
			return
		}

		// Get key ID from token header
		keyID, ok := token.Header["kid"].(string)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Missing key ID in token", nil)
			c.Abort()
			return
		}

		// Get the key pair for this key ID
		keyPair := m.keyManager.GetKeyPairByID(keyID)
		if keyPair == nil {
			response.Error(c, http.StatusUnauthorized, "Invalid key ID", nil)
			c.Abort()
			return
		}

		// Parse and validate token with the correct public key
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, auth.ErrInvalidToken
			}
			return keyPair.PublicKey, nil
		})

		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "Invalid token", err)
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Invalid token claims", nil)
			c.Abort()
			return
		}

		// Get user ID from claims
		userIDClaim, ok := claims["userId"]
		if !ok || userIDClaim == nil {
			response.Error(c, http.StatusUnauthorized, "Missing user ID in token", nil)
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDClaim.(string))
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid user ID in token", err)
			c.Abort()
			return
		}

		// Get tenant ID from claims
		tenantID, err := uuid.Parse(claims["tenantId"].(string))
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid tenant ID in token", err)
			c.Abort()
			return
		}

		// Get user from database
		user, err := m.userRepo.GetUserByID(userID)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "User not found", err)
			c.Abort()
			return
		}

		// Get user's tenant access
		tenantAccess, err := m.userRepo.GetUserTenantAccess(userID)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Failed to get user tenant access", err)
			c.Abort()
			return
		}

		// Get current tenant
		currentTenant, err := m.userRepo.GetTenantByID(tenantID)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Tenant not found", err)
			c.Abort()
			return
		}

		// Verify user has access to the tenant
		hasAccess := false
		for _, access := range tenantAccess {
			if access.TenantID == tenantID {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			response.Error(c, http.StatusForbidden, "Access denied to tenant", nil)
			c.Abort()
			return
		}

		// Set context values
		c.Set("user", user)
		c.Set("currentTenant", currentTenant)
		c.Set("tenantAccess", tenantAccess)

		c.Next()
	}
}
