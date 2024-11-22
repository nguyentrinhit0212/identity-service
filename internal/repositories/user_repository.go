package repositories

import (
	"identity-service/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
    CreateOrUpdateUser(oauthUser *models.OAuthUser) (models.User, error)
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{
        db: db,
    }
}

func (r *userRepository) CreateOrUpdateUser(oauthUser *models.OAuthUser) (models.User, error) {
    tx := r.db.Begin()
    if tx.Error != nil {
        return models.User{}, tx.Error
    }

    var user models.User
    var oauthProvider models.OAuthProvider

    err := tx.Where("email = ?", oauthUser.Email).First(&user).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            user = models.User{
                Email: oauthUser.Email,
            }
            if err := tx.Create(&user).Error; err != nil {
                tx.Rollback()
                return models.User{}, err
            }

            oauthProvider = models.OAuthProvider{
                UserID:         user.ID,
                Provider:       oauthUser.Provider,
                ProviderUserID: oauthUser.ID,
            }
            if err := tx.Create(&oauthProvider).Error; err != nil {
                tx.Rollback()
                return models.User{}, err
            }
        } else {
            tx.Rollback()
            return models.User{}, err
        }
    } else {
        err = tx.Where("provider = ? AND provider_user_id = ?", oauthUser.Provider, oauthUser.ID).FirstOrCreate(&oauthProvider).Error
        if err != nil {
            tx.Rollback()
            return models.User{}, err
        }
    }

    return user, tx.Commit().Error
}