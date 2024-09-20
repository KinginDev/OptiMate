package repositories

import (
	"user-service/cmd/internal/models"

	"gorm.io/gorm"
)

// JWTRepository is a repository for the personal token
// It contains the database connection
type JWTRepository struct {
	DB *gorm.DB
}

// NewJWTTokenRepository creates a new instance of JWTRepository
// It returns a pointer to the instance
func NewJWTTokenRepository(db *gorm.DB) *JWTRepository {
	return &JWTRepository{DB: db}
}

// StoreToken stores a token in the database
// It returns an error if the operation fails
func (repo *JWTRepository) StoreToken(token *models.PersonalToken) error {
	return repo.DB.Create(token).Error
}

// RevokeToken revokes a token by setting the revoked field to true
// It returns an error if the operation fails
func (repo *JWTRepository) RevokeToken(tokenString string) error {
	result := repo.DB.Model(&models.PersonalToken{}).Where("token = ?", tokenString).Update("revoked", true)
	return result.Error
}

// CheckTokenRevocation checks if a token has been revoked
// It returns a boolean and an error
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
