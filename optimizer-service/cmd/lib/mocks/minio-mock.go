// Package mocks
package mocks

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/mock"
)

// MockMinioClient is a mock type for the minio client
type MockMinioClient struct {
	mock.Mock
}

// NewMockMinioClient creates a new mock minio client
func NewMockMinioClient() *MockMinioClient {
	return &MockMinioClient{}
}

// PutObject mocks the PutObject method
// It returns the upload info and an error
func (m *MockMinioClient) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	args := m.Called(ctx, bucketName, objectName, reader, objectSize, opts)
	return args.Get(0).(minio.UploadInfo), args.Error(1)
}

// GetObject mocks the GetObject method
// It returns the object and an error
func (m *MockMinioClient) GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	args := m.Called(ctx, bucketName, objectName, opts)
	return nil, args.Error(0)
}

// RemoveObject mocks the RemoveObject method
// It returns an error
func (m *MockMinioClient) RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error {
	args := m.Called(ctx, bucketName, objectName, opts)
	return args.Error(0)
}

// StatObject mocks the StatObject method
// It returns the object info and an error
func (m *MockMinioClient) StatObject(ctx context.Context, bucketName, objectName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error) {

	args := m.Called(ctx, bucketName, objectName, opts)
	return minio.ObjectInfo{}, args.Error(0)
}
