// Package mocks contains the mocked objects for the file service
package mocks

import (
	"io"
	"optimizer-service/cmd/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockFileService struct {
	mock.Mock
}

func (m *MockFileService) UploadFile(userId string, fileData io.Reader, fileName string) (*models.File, error) {
	args := m.Called(userId, fileData, fileName)
	return args.Get(0).(*models.File), args.Error(1)
}
