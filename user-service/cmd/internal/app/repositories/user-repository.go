// Package repositories
package repositories

import (
	"user-service/cmd/internal/models"

	"gorm.io/gorm"
)

// UserRepository is a repository for the user
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser creates a new user in the database
func (repo *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	if err := repo.DB.Create(user).Error; err != nil {
		return nil, err
	}

	// Preload Tokens after creating the user
	if err := repo.DB.Model(user).Preload("Tokens").First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
// It returns a user and an error
func (repo *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	if err := repo.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user by ID
// It returns a user and an error
func (repo *UserRepository) GetUserByID(db *gorm.DB, id string) (*models.User, error) {
	user := &models.User{}
	err := db.Where("id = ?", id).First(user).Error
	return user, err
}

// GetTokensByUserID retrieves all tokens for a user
// It returns a list of tokens and an error
func (repo *UserRepository) GetTokensByUserID(id string) ([]models.PersonalToken, error) {
	var tokens []models.PersonalToken
	if err := repo.DB.Where("user_id = ?", id).Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}
