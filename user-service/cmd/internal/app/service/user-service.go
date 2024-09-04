// Package service
package service

import (
	"errors"
	"user-service/cmd/internal/app/repositories"
	"user-service/cmd/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserService is a service that handles user operations
type UserService struct {
	Repo *repositories.UserRepository
}

// NewUserService creates a new instance of UserService
// It returns a pointer to the instance
func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

// RegisterUser registers a new user
// It returns a user and an error if the operation fails
func (s *UserService) RegisterUser(input *models.RegisterInput) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:        uuid.New().String(),
		Email:     input.Email,
		Password:  string(hashedPassword),
		Firstname: &input.Firstname,
		Lastname:  &input.Firstname,
	}

	return s.Repo.CreateUser(user)
}

// AuthenticateUser authenticates a user
// It returns a user and an error if the operation fails
func (s *UserService) AuthenticateUser(email, password string) (*models.User, error) {
	user, err := s.Repo.GetUserByEmail(email)

	if err != nil {
		return nil, err
	}

	if !user.ComparePassword(password) {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
