package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"identity-service/internal/auth"
	jwtmanager "identity-service/internal/auth/jwt"
	"identity-service/internal/models"
	"identity-service/internal/repositories"
	"strings"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx *gin.Context, credentials *models.LoginCredentials) (*models.Session, error)
	Logout(ctx *gin.Context) error
	RefreshToken(ctx *gin.Context, refreshToken string) (*models.Session, error)
	ValidateToken(token string) (*models.Session, error)
	HandleOAuthCallback(ctx *gin.Context, provider string, code string) (*models.OAuthUser, error)
	GetOAuthProvider(provider string) (auth.OAuthProviderInterface, error)
	SwitchTenant(ctx *gin.Context, tenantID uuid.UUID) (*models.Session, error)
	GetSession(ctx *gin.Context) (*models.Session, error)
	SendPasswordResetEmail(email string) error
	ResetPassword(token string, newPassword string) error
	VerifyEmail(token string) error
	ListSessions(ctx *gin.Context) ([]*models.Session, error)
	RevokeSession(ctx *gin.Context, sessionID uuid.UUID) error
	RevokeAllSessions(ctx *gin.Context) error
	GetSecuritySettings(ctx *gin.Context) (*models.SecuritySettings, error)
	UpdateSecuritySettings(ctx *gin.Context, settings models.SecuritySettings) error
	EnableMFA(ctx *gin.Context) (secret string, qrCode string, err error)
	DisableMFA(ctx *gin.Context, password string) error
	VerifyMFA(ctx *gin.Context, token string) error
	CreateSession(ctx *gin.Context, user *models.User, tenantID uuid.UUID) (*models.Session, error)
}

type authService struct {
	userService    UserService
	sessionRepo    repositories.SessionRepository
	keyManager     *jwtmanager.KeyManager
	oauthProviders map[string]auth.OAuthProviderInterface
}

func NewAuthService(userService UserService, sessionRepo repositories.SessionRepository, keyManager *jwtmanager.KeyManager) AuthService {
	providers := map[string]auth.OAuthProviderInterface{
		"google": auth.NewGoogleProvider(),
		// Add more providers here as needed
		// "github": auth.NewGithubProvider(),
		// "microsoft": auth.NewMicrosoftProvider(),
	}

	return &authService{
		userService:    userService,
		sessionRepo:    sessionRepo,
		keyManager:     keyManager,
		oauthProviders: providers,
	}
}

func (s *authService) GetOAuthProvider(provider string) (auth.OAuthProviderInterface, error) {
	if p, exists := s.oauthProviders[provider]; exists {
		return p, nil
	}
	return nil, fmt.Errorf("provider %s not supported", provider)
}

func (s *authService) HandleOAuthCallback(ctx *gin.Context, provider string, code string) (*models.OAuthUser, error) {
	oauthProvider, err := s.GetOAuthProvider(provider)
	if err != nil {
		return nil, err
	}

	// Exchange code for token
	token, err := oauthProvider.ExchangeToken(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %v", err)
	}

	// Get user info from provider
	userInfo, err := oauthProvider.FetchUserInfo(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	return userInfo, nil
}

func (s *authService) Login(ctx *gin.Context, credentials *models.LoginCredentials) (*models.Session, error) {
	// Get user by email
	user, err := s.userService.GetUserByEmail(credentials.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := s.userService.VerifyPassword(user.ID, credentials.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Get user's tenants
	tenants, err := s.userService.GetUserTenants(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tenants: %v", err)
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
		return nil, fmt.Errorf("personal tenant not found")
	}

	// Create session with personal tenant
	return s.CreateSession(ctx, user, personalTenant.ID)
}

func (s *authService) Logout(ctx *gin.Context) error {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		return errors.New("no token provided")
	}

	// Parse token to get session ID
	claims, err := s.parseToken(token)
	if err != nil {
		return err
	}

	// Delete session
	return s.sessionRepo.DeleteSession(claims.SessionID)
}

func (s *authService) RefreshToken(ctx *gin.Context, refreshToken string) (*models.Session, error) {
	// Parse refresh token
	claims, err := s.parseToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Verify it's a refresh token
	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}

	// Get existing session
	session, err := s.sessionRepo.GetSession(claims.SessionID)
	if err != nil {
		return nil, err
	}

	// Get user to verify they still exist
	user, err := s.userService.GetUser(session.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Create new session with same tenant
	return s.CreateSession(ctx, user, session.TenantID)
}

func (s *authService) ValidateToken(token string) (*models.Session, error) {
	// Parse token
	claims, err := s.parseToken(token)
	if err != nil {
		return nil, err
	}

	// Get session
	session, err := s.sessionRepo.GetSession(claims.SessionID)
	if err != nil {
		return nil, err
	}

	// Verify session is not expired
	if time.Now().After(session.ExpiresAt) {
		return nil, errors.New("session expired")
	}

	// Update last used time
	session.LastUsedAt = time.Now()
	if err := s.sessionRepo.UpdateSession(session); err != nil {
		return nil, fmt.Errorf("failed to update session: %v", err)
	}

	return session, nil
}

func (s *authService) SwitchTenant(ctx *gin.Context, tenantID uuid.UUID) (*models.Session, error) {
	token := ctx.GetHeader("Authorization")
	claims, err := s.parseToken(token)
	if err != nil {
		return nil, err
	}

	// Create new session with updated tenant
	sessionID := uuid.New()
	session := &models.Session{
		ID:           sessionID,
		UserID:       claims.UserID,
		TenantID:     tenantID,
		AccessToken:  s.generateToken(claims.UserID, sessionID, "access", tenantID, 15*time.Minute),
		RefreshToken: s.generateToken(claims.UserID, sessionID, "refresh", tenantID, 7*24*time.Hour),
		ExpiresAt:    time.Now().Add(15 * time.Minute),
		CreatedAt:    time.Now(),
		LastUsedAt:   time.Now(),
		IPAddress:    ctx.ClientIP(),
		UserAgent:    ctx.GetHeader("User-Agent"),
	}

	if err := s.sessionRepo.CreateSession(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	return session, nil
}

func (s *authService) GetSession(ctx *gin.Context) (*models.Session, error) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		return nil, errors.New("no token provided")
	}

	return s.ValidateToken(token)
}

func (s *authService) SendPasswordResetEmail(email string) error {
	// Get user by email
	_, err := s.userService.GetUserByEmail(email)
	if err != nil {
		return err
	}

	// TODO: Implement password reset functionality
	return nil
}

func (s *authService) ResetPassword(_ string, _ string) error {
	// TODO: Verify token and get associated user
	// TODO: Hash new password
	// TODO: Update user password
	// TODO: Invalidate token
	return nil
}

func (s *authService) VerifyEmail(_ string) error {
	// TODO: Verify token and get associated email verification
	// TODO: Mark email as verified
	// TODO: Update user's email verified status
	return nil
}

func (s *authService) ListSessions(ctx *gin.Context) ([]*models.Session, error) {
	claims, err := s.parseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		return nil, err
	}
	return s.sessionRepo.ListUserSessions(claims.UserID)
}

func (s *authService) RevokeSession(ctx *gin.Context, sessionID uuid.UUID) error {
	claims, err := s.parseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		return err
	}

	// Verify user owns the session
	session, err := s.sessionRepo.GetSession(sessionID)
	if err != nil {
		return err
	}
	if session.UserID != claims.UserID {
		return errors.New("unauthorized")
	}

	return s.sessionRepo.DeleteSession(sessionID)
}

func (s *authService) RevokeAllSessions(ctx *gin.Context) error {
	claims, err := s.parseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		return err
	}

	sessions, err := s.sessionRepo.ListUserSessions(claims.UserID)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		if session.ID != claims.SessionID { // Don't revoke current session
			if err := s.sessionRepo.DeleteSession(session.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *authService) GetSecuritySettings(ctx *gin.Context) (*models.SecuritySettings, error) {
	claims, err := s.parseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		return nil, err
	}

	user, err := s.userService.GetUser(claims.UserID)
	if err != nil {
		return nil, err
	}

	return &models.SecuritySettings{
		MFAEnabled: user.MFAEnabled,
		LastLogin:  user.LastLoginAt,
	}, nil
}

func (s *authService) UpdateSecuritySettings(ctx *gin.Context, settings models.SecuritySettings) error {
	claims, err := s.parseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		return err
	}

	settingsBytes, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	update := &models.UserUpdate{
		Settings: (*json.RawMessage)(&settingsBytes),
	}

	return s.userService.UpdateUser(claims.UserID, update)
}

func (s *authService) EnableMFA(ctx *gin.Context) (string, string, error) {
	claims, err := s.parseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		return "", "", err
	}

	// Generate MFA secret
	secret := auth.GenerateMFASecret()
	qrCode := auth.GenerateMFAQRCode(secret)

	// Update user's MFA settings
	update := &models.UserUpdate{
		Settings: new(json.RawMessage),
	}
	settings := map[string]interface{}{
		"mfa_enabled": true,
		"mfa_secret":  secret,
	}
	settingsBytes, err := json.Marshal(settings)
	if err != nil {
		return "", "", err
	}
	*update.Settings = json.RawMessage(settingsBytes)

	if err := s.userService.UpdateUser(claims.UserID, update); err != nil {
		return "", "", err
	}

	return secret, qrCode, nil
}

func (s *authService) DisableMFA(ctx *gin.Context, password string) error {
	claims, err := s.parseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		return err
	}

	// Verify password
	if err := s.userService.VerifyPassword(claims.UserID, password); err != nil {
		return errors.New("invalid password")
	}

	// Update user's MFA settings
	update := &models.UserUpdate{
		Settings: new(json.RawMessage),
	}
	settings := map[string]interface{}{
		"mfa_enabled": false,
		"mfa_secret":  nil,
	}
	settingsBytes, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	*update.Settings = json.RawMessage(settingsBytes)

	return s.userService.UpdateUser(claims.UserID, update)
}

func (s *authService) VerifyMFA(ctx *gin.Context, token string) error {
	claims, err := s.parseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		return err
	}

	user, err := s.userService.GetUser(claims.UserID)
	if err != nil {
		return err
	}

	// Get MFA secret from user settings
	var settings struct {
		MFASecret string `json:"mfaSecret"`
	}
	if err := json.Unmarshal(user.Settings, &settings); err != nil {
		return err
	}

	// Verify TOTP token
	if !auth.VerifyTOTP(settings.MFASecret, token) {
		return errors.New("invalid MFA token")
	}

	return nil
}

func (s *authService) CreateSession(ctx *gin.Context, user *models.User, tenantID uuid.UUID) (*models.Session, error) {
	// Create session with a new ID
	sessionID := uuid.New()
	session := &models.Session{
		ID:           sessionID,
		UserID:       user.ID,
		TenantID:     tenantID,
		AccessToken:  s.generateToken(user.ID, sessionID, "access", tenantID, 15*time.Minute),
		RefreshToken: s.generateToken(user.ID, sessionID, "refresh", tenantID, 7*24*time.Hour),
		ExpiresAt:    time.Now().Add(15 * time.Minute),
		CreatedAt:    time.Now(),
		LastUsedAt:   time.Now(),
		IPAddress:    ctx.ClientIP(),
		UserAgent:    ctx.GetHeader("User-Agent"),
	}

	// Save session
	if err := s.sessionRepo.CreateSession(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	return session, nil
}

type Claims struct {
	UserID    uuid.UUID `json:"userId"`
	SessionID uuid.UUID `json:"sessionId"`
	TokenType string    `json:"tokenType"`
	TenantID  uuid.UUID `json:"tenantId"`
	jwt.RegisteredClaims
}

func (c Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}

func (s *authService) generateToken(userID uuid.UUID, sessionID uuid.UUID, tokenType string, tenantID uuid.UUID, expiry time.Duration) string {
	claims := Claims{
		UserID:    userID,
		SessionID: sessionID,
		TokenType: tokenType,
		TenantID:  tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		},
	}

	// Log claims for debugging
	log.Printf("Generating token for user %s, session %s, type %s", userID, sessionID, tokenType)

	token, err := s.keyManager.SignToken(claims)
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return ""
	}

	return token
}

func (s *authService) parseToken(tokenString string) (*Claims, error) {
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := &Claims{}
	err := s.keyManager.VerifyToken(tokenString, claims)
	if err != nil {
		log.Printf("Token parsing error: %v", err)
		return nil, err
	}

	return claims, nil
}
