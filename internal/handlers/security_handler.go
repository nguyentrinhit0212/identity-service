package handlers

import (
	"identity-service/internal/models"
	"identity-service/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SecurityHandler handles all security-related HTTP requests
type SecurityHandler struct {
	securityService services.SecurityService
}

// NewSecurityHandler creates a new security handler instance
func NewSecurityHandler(securityService services.SecurityService) *SecurityHandler {
	return &SecurityHandler{
		securityService: securityService,
	}
}

// ListWhitelistedIPs returns a paginated list of whitelisted IPs
func (h *SecurityHandler) ListWhitelistedIPs(c *gin.Context) {
	tenantID := h.getTenantID(c)
	ips, err := h.securityService.ListWhitelistedIPs(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ips": ips,
	})
}

// AddWhitelistedIP adds a new IP to the whitelist
func (h *SecurityHandler) AddWhitelistedIP(c *gin.Context) {
	tenantID := h.getTenantID(c)
	ip := c.PostForm("ip")
	if ip == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IP address is required"})
		return
	}

	err := h.securityService.AddWhitelistedIP(tenantID, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "IP added to whitelist successfully"})
}

// RemoveWhitelistedIP removes an IP from the whitelist
func (h *SecurityHandler) RemoveWhitelistedIP(c *gin.Context) {
	tenantID := h.getTenantID(c)
	ip := c.Param("ip")
	if ip == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IP address is required"})
		return
	}

	if err := h.securityService.RemoveWhitelistedIP(tenantID, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "IP removed from whitelist successfully"})
}

// ListAPIKeys returns a list of API keys
func (h *SecurityHandler) ListAPIKeys(c *gin.Context) {
	tenantID := h.getTenantID(c)
	keys, err := h.securityService.ListAPIKeys(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

// CreateAPIKey creates a new API key
func (h *SecurityHandler) CreateAPIKey(c *gin.Context) {
	tenantID := h.getTenantID(c)
	name := c.PostForm("name")
	var expiresAt *time.Time
	if exp := c.PostForm("expires_at"); exp != "" {
		t, err := time.Parse(time.RFC3339, exp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiration date format"})
			return
		}
		expiresAt = &t
	}

	key, err := h.securityService.CreateAPIKey(tenantID, name, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, key)
}

// RevokeAPIKey revokes an API key
func (h *SecurityHandler) RevokeAPIKey(c *gin.Context) {
	tenantID := h.getTenantID(c)
	keyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API key ID"})
		return
	}

	if err := h.securityService.RevokeAPIKey(tenantID, keyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key revoked successfully"})
}

// GetSecurityAuditLogs returns security audit logs
func (h *SecurityHandler) GetSecurityAuditLogs(c *gin.Context) {
	tenantID := h.getTenantID(c)
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}
	filter := c.QueryMap("filter")

	logs, total, err := h.securityService.GetSecurityAuditLogs(tenantID, page, limit, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetSecurityPolicies returns security policies
func (h *SecurityHandler) GetSecurityPolicies(c *gin.Context) {
	tenantID := h.getTenantID(c)
	policies, err := h.securityService.GetSecurityPolicies(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, policies)
}

// UpdateSecurityPolicies updates security policies
func (h *SecurityHandler) UpdateSecurityPolicies(c *gin.Context) {
	tenantID := h.getTenantID(c)
	var policies models.SecurityPolicies
	if err := c.ShouldBindJSON(&policies); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.securityService.UpdateSecurityPolicies(tenantID, &policies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Security policies updated successfully"})
}

// GetSecurityMetrics returns security metrics
func (h *SecurityHandler) GetSecurityMetrics(c *gin.Context) {
	tenantID := h.getTenantID(c)
	metrics, err := h.securityService.GetSecurityMetrics(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetSecurityAlerts returns security alerts
func (h *SecurityHandler) GetSecurityAlerts(c *gin.Context) {
	tenantID := h.getTenantID(c)
	status := c.DefaultQuery("status", "")
	alerts, err := h.securityService.GetSecurityAlerts(tenantID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}

// UpdateSecurityAlert updates a security alert
func (h *SecurityHandler) UpdateSecurityAlert(c *gin.Context) {
	tenantID := h.getTenantID(c)
	alertID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
		return
	}

	status := c.PostForm("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	err = h.securityService.UpdateSecurityAlert(tenantID, alertID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Security alert updated successfully"})
}

func (h *SecurityHandler) getTenantID(c *gin.Context) uuid.UUID {
	return uuid.MustParse(c.GetString("tenantID"))
}

// GetAPIKey returns an API key
func (h *SecurityHandler) GetAPIKey(c *gin.Context) {
	tenantID := h.getTenantID(c)
	keyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API key ID"})
		return
	}

	apiKey, err := h.securityService.GetAPIKey(tenantID, keyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}
	c.JSON(http.StatusOK, apiKey)
}

// UpdateAPIKey updates an API key
func (h *SecurityHandler) UpdateAPIKey(c *gin.Context) {
	tenantID := h.getTenantID(c)
	keyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API key ID"})
		return
	}

	name := c.PostForm("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	apiKey, err := h.securityService.UpdateAPIKey(tenantID, keyID, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, apiKey)
}

// DeleteAPIKey deletes (revokes) an API key
func (h *SecurityHandler) DeleteAPIKey(c *gin.Context) {
	tenantID := h.getTenantID(c)
	keyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API key ID"})
		return
	}

	if err := h.securityService.RevokeAPIKey(tenantID, keyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key deleted successfully"})
}

// ListAuditLogs returns a list of audit logs
func (h *SecurityHandler) ListAuditLogs(c *gin.Context) {
	h.GetSecurityAuditLogs(c)
}

// GetAuditLogEntry returns a single audit log entry
func (h *SecurityHandler) GetAuditLogEntry(c *gin.Context) {
	tenantID := h.getTenantID(c)
	logID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audit log ID"})
		return
	}

	log, err := h.securityService.GetAuditLogEntry(tenantID, logID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Audit log entry not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// ListSecurityPolicies returns security policies
func (h *SecurityHandler) ListSecurityPolicies(c *gin.Context) {
	h.GetSecurityPolicies(c)
}

// TestSecurityPolicy tests a security policy
func (h *SecurityHandler) TestSecurityPolicy(c *gin.Context) {
	tenantID := h.getTenantID(c)
	var policy models.SecurityPolicies
	if err := c.ShouldBindJSON(&policy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.securityService.TestSecurityPolicy(tenantID, &policy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}
