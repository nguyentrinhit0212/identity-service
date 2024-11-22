package services

import (
	"identity-service/internal/models"
	"identity-service/internal/repositories"
)

type UserService interface {
    CreateOrUpdateUser(oauthUser *models.OAuthUser) (models.User, error)
}

type userService struct {
    userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
    return &userService{
        userRepo: userRepo,
    }
}

func (s *userService) CreateOrUpdateUser(oauthUser *models.OAuthUser) (models.User, error) {
    return s.userRepo.CreateOrUpdateUser(oauthUser)
}