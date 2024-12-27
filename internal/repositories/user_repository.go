package repositories

import (
	"identity-service/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	ListUsers(page, limit int, search string, filter map[string]string) ([]*models.User, int64, error)
	CreateUser(user *models.User) error
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uuid.UUID) error
	GetUserProfile(id uuid.UUID) (*models.UserProfile, error)
	UpdateUserProfile(id uuid.UUID, profile *models.UserProfile) error
	GetUserTenants(id uuid.UUID) ([]*models.Tenant, error)
	AddUserToTenant(userID, tenantID uuid.UUID, roles []string) error
	RemoveUserFromTenant(userID, tenantID uuid.UUID) error
	UpdateUserRole(userID, tenantID uuid.UUID, role string) error
	GetUserTenantAccess(userID uuid.UUID) ([]models.UserTenantAccess, error)
	GetTenantByID(id uuid.UUID) (*models.Tenant, error)
}

type userRepository struct {
	db GormDB
}

func NewUserRepository(db GormDB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) ListUsers(page, limit int, search string, filter map[string]string) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	query := r.db.Model(&models.User{})

	// Apply search if provided
	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
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
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) DeleteUser(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

func (r *userRepository) GetUserProfile(id uuid.UUID) (*models.UserProfile, error) {
	var profile models.UserProfile
	err := r.db.First(&profile, "user_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *userRepository) UpdateUserProfile(id uuid.UUID, profile *models.UserProfile) error {
	profile.UserID = id
	return r.db.Save(profile).Error
}

func (r *userRepository) GetUserTenants(id uuid.UUID) ([]*models.Tenant, error) {
	var tenants []*models.Tenant
	err := r.db.Joins("JOIN user_tenant_access ON user_tenant_access.tenant_id = tenants.id").
		Where("user_tenant_access.user_id = ?", id).
		Find(&tenants).Error
	return tenants, err
}

func (r *userRepository) AddUserToTenant(userID, tenantID uuid.UUID, roles []string) error {
	access := &models.UserTenantAccess{
		UserID:   userID,
		TenantID: tenantID,
		Roles:    roles,
	}
	return r.db.Create(access).Error
}

func (r *userRepository) RemoveUserFromTenant(userID, tenantID uuid.UUID) error {
	return r.db.Delete(&models.UserTenantAccess{}, "user_id = ? AND tenant_id = ?", userID, tenantID).Error
}

func (r *userRepository) UpdateUserRole(userID, tenantID uuid.UUID, role string) error {
	return r.db.Model(&models.UserTenantAccess{}).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		Update("roles", []string{role}).Error
}

func (r *userRepository) GetUserTenantAccess(userID uuid.UUID) ([]models.UserTenantAccess, error) {
	var accesses []models.UserTenantAccess
	err := r.db.Where("user_id = ?", userID).Find(&accesses).Error
	return accesses, err
}

func (r *userRepository) GetTenantByID(id uuid.UUID) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.Where("id = ?", id).First(&tenant).Error
	return &tenant, err
}
