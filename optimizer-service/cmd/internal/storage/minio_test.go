// Package storage
package storage

import (
	"errors"
	"optimizer-service/cmd/lib/mocks"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaveSuccess(t *testing.T) {
	client := mocks.NewMockMinioClient()
	client.On("PutObject", mock.Anything, "bucket-name", "file-path", mock.Anything, int64(-1), mock.Anything).Return(minio.UploadInfo{}, nil)

	storage := NewMinIOStorage(client, "bucket-name")
	err := storage.Save("file-path", nil) // nil is used for simplicity

	assert.Nil(t, err)
}

func TestSaveFailure(t *testing.T) {
	client := new(mocks.MockMinioClient)
	client.On("PutObject", mock.Anything, "bucket-name", "file-path", mock.Anything, int64(-1), mock.Anything).Return(minio.UploadInfo{}, errors.New("failed to upload"))

	storage := NewMinIOStorage(client, "bucket-name")
	err := storage.Save("file-path", nil) // nil is used for simplicity

	assert.NotNil(t, err)
	assert.Equal(t, "failed to upload", err.Error())
}

// Write similar tests for Retrieve, Delete, and Exists
