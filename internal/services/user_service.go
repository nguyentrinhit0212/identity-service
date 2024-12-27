package services

import (
	"identity-service/internal/models"
	"identity-service/internal/repositories"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles user-related business logic
type UserService interface {
	ListUsers(ctx *gin.Context) ([]*models.User, error)
	CreateUser(ctx *gin.Context, user *models.User) error
	GetUser(ctx *gin.Context, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx *gin.Context, email string) (*models.User, error)
	UpdateUser(ctx *gin.Context, id uuid.UUID, update *models.UserUpdate) error
	DeleteUser(ctx *gin.Context, id uuid.UUID) error
	GetUserProfile(ctx *gin.Context, id uuid.UUID) (*models.UserProfile, error)
	UpdateUserProfile(ctx *gin.Context, id uuid.UUID, profile *models.UserProfile) error
	GetUserTenants(ctx *gin.Context, id uuid.UUID) ([]*models.Tenant, error)
	AddUserToTenant(userID uuid.UUID, tenantID uuid.UUID, roles []string) error
	RemoveUserFromTenant(userID uuid.UUID, tenantID uuid.UUID) error
	UpdateUserRole(userID uuid.UUID, tenantID uuid.UUID, role string) error
	CreateOrUpdateUser(oauthUser *models.OAuthUser) (*models.User, error)
	UpdatePassword(ctx *gin.Context, userID uuid.UUID, currentPassword, newPassword string) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) ListUsers(ctx *gin.Context) ([]*models.User, error) {
	users, _, err := s.userRepo.ListUsers(0, 0, "", nil)
	return users, err
}

func (s *userService) CreateUser(ctx *gin.Context, user *models.User) error {
	return s.userRepo.CreateUser(user)
}

func (s *userService) GetUser(ctx *gin.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *userService) GetUserByEmail(ctx *gin.Context, email string) (*models.User, error) {
	return s.userRepo.GetUserByEmail(email)
}

func (s *userService) UpdateUser(ctx *gin.Context, id uuid.UUID, update *models.UserUpdate) error {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}

	// Apply updates
	if update.Email != nil {
		user.Email = *update.Email
	}
	if update.Name != nil {
		user.Name = *update.Name
	}
	if update.Status != nil {
		user.Status = *update.Status
	}
	if update.Role != nil {
		user.Role = *update.Role
	}
	if update.Settings != nil {
		user.Settings = *update.Settings
	}

	return s.userRepo.UpdateUser(user)
}

func (s *userService) DeleteUser(ctx *gin.Context, id uuid.UUID) error {
	return s.userRepo.DeleteUser(id)
}

func (s *userService) GetUserProfile(ctx *gin.Context, id uuid.UUID) (*models.UserProfile, error) {
	return s.userRepo.GetUserProfile(id)
}

func (s *userService) UpdateUserProfile(ctx *gin.Context, id uuid.UUID, profile *models.UserProfile) error {
	return s.userRepo.UpdateUserProfile(id, profile)
}

func (s *userService) GetUserTenants(ctx *gin.Context, id uuid.UUID) ([]*models.Tenant, error) {
	return s.userRepo.GetUserTenants(id)
}

func (s *userService) AddUserToTenant(userID uuid.UUID, tenantID uuid.UUID, roles []string) error {
	return s.userRepo.AddUserToTenant(userID, tenantID, roles)
}

func (s *userService) RemoveUserFromTenant(userID uuid.UUID, tenantID uuid.UUID) error {
	return s.userRepo.RemoveUserFromTenant(userID, tenantID)
}

func (s *userService) UpdateUserRole(userID uuid.UUID, tenantID uuid.UUID, role string) error {
	return s.userRepo.UpdateUserRole(userID, tenantID, role)
}

func (s *userService) UpdatePassword(ctx *gin.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Verify current password
	if !s.verifyPassword(user.Password, currentPassword) {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := s.hashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.UpdateUser(user)
}

func (s *userService) verifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s *userService) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (s *userService) CreateOrUpdateUser(oauthUser *models.OAuthUser) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(oauthUser.Email)
	if err != nil {
		// Create new user
		user = &models.User{
			Email: oauthUser.Email,
			Name:  oauthUser.Name,
		}
		if err := s.userRepo.CreateUser(user); err != nil {
			return nil, err
		}
	}
	return user, nil
}
