// Package mocks contains the mocked objects for the file service
package mocks

import (
	"io"
	"optimizer-service/cmd/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/mock"
)

// MockFileService is a mock type for the file service
type MockFileService struct {
	mock.Mock
}
type MockAuthService struct {	
	mock.Mock
}
// UploadFile is a mocked method
// It expects a userId, fileData and fileName as input
// It returns a file and an error
func (m *MockFileService) UploadFile(userId string, fileData io.Reader, fileName string) (*models.File, error) {
	args := m.Called(userId, fileData, fileName)
	return args.Get(0).(*models.File), args.Error(1)
}


func (m *MockAuthService) Login(email string, password string) (interface{}, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(token string) (*jwt.Token,error) {
	args := m.Called(token)
	return args.Get(0).(*jwt.Token), args.Error(1)
}