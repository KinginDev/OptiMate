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

// @note: We Can establish any kind of communication we want here,
// with the user service or any other service

// Should return user type or error
func (r *AuthRepository) Login(email, password string) (string, error) {
	// userServiceUrl := "http://localhost:8020/login"
	return "", nil
}

// Would validate the token and return a boolean or error
func (r *AuthRepository) ValidateToken(token string) (bool, error) {
	return false, nil
}
