// Package mocks contains the mocked objects for the file service
package mocks

import (
	"io"
	"optimizer-service/cmd/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockFileService is a mock type for the file service
type MockFileService struct {
	mock.Mock
}

// UploadFile is a mocked method
// It expects a userId, fileData and fileName as input
// It returns a file and an error
func (m *MockFileService) UploadFile(userId string, fileData io.Reader, fileName string) (*models.File, error) {
	args := m.Called(userId, fileData, fileName)
	return args.Get(0).(*models.File), args.Error(1)
}
