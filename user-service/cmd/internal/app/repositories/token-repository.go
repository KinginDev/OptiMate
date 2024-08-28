package repositories

import (
	"user-service/cmd/internal/models"

	"gorm.io/gorm"
)

type JWTRepository struct {
	DB *gorm.DB
}

func NewJWTToken(db *gorm.DB) *JWTRepository {
	return &JWTRepository{DB: db}
}

func (repo *JWTRepository) StoreToken(token *models.PersonalToken) error {
	return repo.DB.Create(token).Error
}

func (repo *JWTRepository) RevokeToken(token *models.PersonalToken) error {
	return repo.DB.Model(&models.PersonalToken{}).Where("token = ?", token).Update("revoked", true).Error
}

func (repo *JWTRepository) CheckTokenRevocation(token string) (bool, error) {
	var t models.PersonalToken
	if err := repo.DB.Where("token = ? AND revoked =?", token, true).First(&t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
