package mocks

import (
	"io"

	"github.com/stretchr/testify/mock"
)

type MockOptimizer struct {
	mock.Mock
}

func (m *MockOptimizer) Optimize(filePath string) (io.ReadCloser, error) {
	args := m.Called(filePath)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockOptimizer) SupportedFormats() []string {
	args := m.Called()
	return args.Get(0).([]string)
}
