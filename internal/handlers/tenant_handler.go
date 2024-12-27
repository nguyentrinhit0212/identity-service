package handlers

import (
	"identity-service/internal/models"
	"identity-service/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TenantHandler handles all tenant-related HTTP requests
type TenantHandler struct {
	tenantService services.TenantService
}

// NewTenantHandler creates a new tenant handler instance
func NewTenantHandler(tenantService services.TenantService) *TenantHandler {
	return &TenantHandler{
		tenantService: tenantService,
	}
}

// ListTenants returns a paginated list of tenants
func (h *TenantHandler) ListTenants(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	filter := c.QueryMap("filter")

	tenants, total, err := h.tenantService.ListTenants(c, page, limit, search, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tenants": tenants,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// CreateTenant creates a new tenant
func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var tenant models.Tenant
	if err := c.ShouldBindJSON(&tenant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdTenant, err := h.tenantService.CreateTenant(c, &tenant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTenant)
}

// GetTenant returns a specific tenant by ID
func (h *TenantHandler) GetTenant(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	tenant, err := h.tenantService.GetTenant(c, tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

// UpdateTenant updates a specific tenant
func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var updates models.TenantUpdate
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedTenant, err := h.tenantService.UpdateTenant(c, tenantID, &updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTenant)
}

// DeleteTenant deletes a specific tenant
func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	if err := h.tenantService.DeleteTenant(c, tenantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant deleted successfully"})
}

// GetTenantSettings returns settings for a specific tenant
func (h *TenantHandler) GetTenantSettings(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	settings, err := h.tenantService.GetTenantSettings(c, tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateTenantSettings updates settings for a specific tenant
func (h *TenantHandler) UpdateTenantSettings(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var settings models.TenantSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSettings, err := h.tenantService.UpdateTenantSettings(c, tenantID, &settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedSettings)
}

// GetTenantMembers returns all members of a specific tenant
func (h *TenantHandler) GetTenantMembers(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	filter := c.QueryMap("filter")

	members, total, err := h.tenantService.GetTenantMembers(c, tenantID, page, limit, search, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"members": members,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// ListTenantMembers returns all members of a specific tenant
func (h *TenantHandler) ListTenantMembers(c *gin.Context) {
	h.GetTenantMembers(c)
}

// GetTenantFeatures returns all features for a specific tenant
func (h *TenantHandler) GetTenantFeatures(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	features, err := h.tenantService.GetTenantFeatures(c, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, features)
}

// UpdateTenantFeatures updates features for a specific tenant
func (h *TenantHandler) UpdateTenantFeatures(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var features models.TenantFeatures
	if err := c.ShouldBindJSON(&features); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedFeatures, err := h.tenantService.UpdateTenantFeatures(c, tenantID, &features)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedFeatures)
}

// SwitchTenant switches the active tenant for the current user
func (h *TenantHandler) SwitchTenant(c *gin.Context) {
	userID := h.getUserID(c)
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	if err := h.tenantService.SwitchTenant(c, userID, tenantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant switched successfully"})
}

// UpgradeTenant upgrades the tenant's subscription plan
func (h *TenantHandler) UpgradeTenant(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var upgrade models.TenantUpgrade
	if err := c.ShouldBindJSON(&upgrade); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.tenantService.UpgradeTenant(c, tenantID, &upgrade); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant upgraded successfully"})
}

// CreateTenantInvite creates a new invite for a tenant
func (h *TenantHandler) CreateTenantInvite(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var invite models.TenantInvite
	if err := c.ShouldBindJSON(&invite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdInvite, err := h.tenantService.CreateTenantInvite(c, tenantID, &invite)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdInvite)
}

// DeleteTenantInvite deletes an invite from a tenant
func (h *TenantHandler) DeleteTenantInvite(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	inviteID, err := uuid.Parse(c.Param("inviteId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invite ID"})
		return
	}

	if err := h.tenantService.DeleteTenantInvite(c, tenantID, inviteID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invite deleted successfully"})
}

// ListTenantFeatures returns all features for a specific tenant
func (h *TenantHandler) ListTenantFeatures(c *gin.Context) {
	h.GetTenantFeatures(c)
}

func (h *TenantHandler) getUserID(c *gin.Context) uuid.UUID {
	return uuid.MustParse(c.GetString("userID"))
}
