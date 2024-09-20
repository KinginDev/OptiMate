package repositories

import (
	"gorm.io/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

// Should return user type or error
func (r *AuthRepository) Login(email, password string) (string, error) {
	return "", nil
}

// Would validate the token and return a boolean or error
func (r *AuthRepository) ValidateToken(token string) (bool, error) {
	return false, nil
}
