// Package mocks
package mocks

import (
	"optimizer-service/cmd/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockFileRepository is a mock type for the file repository
type MockFileRepository struct {
	mock.Mock
}
type MockAuthRepository struct {
	mock.Mock
}
// CreateFile is a mocked method
func (m *MockFileRepository) CreateFile(file *models.File) error {
	args := m.Called(file)
	return args.Error(0)
}

func (m *MockAuthRepository) LoginWithREST(username string, password string) (interface{}, error) {
	args := m.Called(username, password)
	return args.Get(0), args.Error(1)
}

func (m *MockAuthRepository) ValidateToken(token string) (interface{},error) {
	args := m.Called(token)
	return args.Get(0), args.Error(1)
}
