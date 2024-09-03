// Package mocks
package mocks

import (
	"io"

	"github.com/stretchr/testify/mock"
)

// MockStorage is a mock type for the storage
type MockStorage struct {
	mock.Mock
}

// Save is a mocked method
// It expects a filePath and data as input
func (m *MockStorage) Save(filePath string, data io.Reader) error {
	args := m.Called(filePath, data)
	return args.Error(0)
}

// Retrieve is a mocked method
// It returns a reader and an error
func (m *MockStorage) Retrieve(filePath string) (io.ReadCloser, error) {
	args := m.Called(filePath)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// Delete is a mocked method
// It returns an error
func (m *MockStorage) Delete(filePath string) error {
	args := m.Called(filePath)
	return args.Error(0)
}

// Exists is a mocked method
// It returns a boolean and an error
func (m *MockStorage) Exists(filePath string) (bool, error) {
	args := m.Called(filePath)
	return args.Bool(0), args.Error(1)
}
