package repositories

import (
	"identity-service/internal/models"

	"github.com/google/uuid"
)

type TenantRepository interface {
	ListTenants(page, limit int, search string, filter map[string]string) ([]*models.Tenant, int64, error)
	CreateTenant(tenant *models.Tenant) error
	GetTenantByID(id uuid.UUID) (*models.Tenant, error)
	GetTenantBySlug(slug string) (*models.Tenant, error)
	UpdateTenant(tenant *models.Tenant) error
	DeleteTenant(id uuid.UUID) error
	GetTenantMembers(tenantID uuid.UUID, page, limit int, search string, filter map[string]string) ([]*models.UserTenantAccess, int64, error)
	GetUserTenantAccess(userID, tenantID uuid.UUID) (*models.UserTenantAccess, error)
	CreateTenantInvite(tenantID uuid.UUID, invite *models.TenantInvite) (*models.TenantInvite, error)
	DeleteTenantInvite(tenantID, inviteID uuid.UUID) error
}

type tenantRepository struct {
	db GormDB
}

func NewTenantRepository(db GormDB) TenantRepository {
	return &tenantRepository{
		db: db,
	}
}

func (r *tenantRepository) ListTenants(page, limit int, search string, filter map[string]string) ([]*models.Tenant, int64, error) {
	var tenants []*models.Tenant
	var total int64

	query := r.db.Model(&models.Tenant{})

	// Apply search if provided
	if search != "" {
		query = query.Where("name ILIKE ? OR slug ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Apply filters
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&tenants).Error; err != nil {
		return nil, 0, err
	}

	return tenants, total, nil
}

func (r *tenantRepository) CreateTenant(tenant *models.Tenant) error {
	return r.db.Create(tenant).Error
}

func (r *tenantRepository) GetTenantByID(id uuid.UUID) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.First(&tenant, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *tenantRepository) GetTenantBySlug(slug string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.First(&tenant, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *tenantRepository) UpdateTenant(tenant *models.Tenant) error {
	return r.db.Save(tenant).Error
}

func (r *tenantRepository) DeleteTenant(id uuid.UUID) error {
	return r.db.Delete(&models.Tenant{}, "id = ?", id).Error
}

func (r *tenantRepository) GetTenantMembers(tenantID uuid.UUID, page, limit int, search string, filter map[string]string) ([]*models.UserTenantAccess, int64, error) {
	var members []*models.UserTenantAccess
	var total int64

	query := r.db.Model(&models.UserTenantAccess{}).
		Preload("User").
		Where("tenant_id = ?", tenantID)

	// Apply search if provided
	if search != "" {
		query = query.Joins("JOIN users ON users.id = user_tenant_access.user_id").
			Where("users.name ILIKE ? OR users.email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Apply filters
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&members).Error

	return members, total, err
}

func (r *tenantRepository) GetUserTenantAccess(userID, tenantID uuid.UUID) (*models.UserTenantAccess, error) {
	var access models.UserTenantAccess
	err := r.db.Where("user_id = ? AND tenant_id = ?", userID, tenantID).First(&access).Error
	if err != nil {
		return nil, err
	}
	return &access, nil
}

func (r *tenantRepository) CreateTenantInvite(tenantID uuid.UUID, invite *models.TenantInvite) (*models.TenantInvite, error) {
	invite.TenantID = tenantID
	if err := r.db.Create(invite).Error; err != nil {
		return nil, err
	}
	return invite, nil
}

func (r *tenantRepository) DeleteTenantInvite(tenantID, inviteID uuid.UUID) error {
	return r.db.Delete(&models.TenantInvite{}, "tenant_id = ? AND id = ?", tenantID, inviteID).Error
}
