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

// CreateFile is a mocked method
func (m *MockFileRepository) CreateFile(file *models.File) error {
	args := m.Called(file)
	return args.Error(0)
}
