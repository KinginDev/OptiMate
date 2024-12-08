package mocks

import (
	"io"
	"optimizer-service/cmd/internal/app/interfaces"
	"optimizer-service/cmd/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockOptimizer struct {
	mock.Mock
}

func (m *MockOptimizer) Optimize(filePath string, file *models.File, oParam *interfaces.OptimizerParams) (io.ReadCloser, error) {
	args := m.Called(filePath, file, oParam)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockOptimizer) SupportedFormats() []string {
	args := m.Called()
	return args.Get(0).([]string)
}
