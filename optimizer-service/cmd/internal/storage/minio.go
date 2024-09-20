package storage

import (
	"context"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
)

// MinIOClient is an interface that defines the methods of the minio client
type MinIOClient interface {
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error)
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error
	StatObject(ctx context.Context, bucketName, objectName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error)
}

// MinIOStorage is a struct that implements the Storage interface
type MinIOStorage struct {
	Client     MinIOClient
	BucketName string
}

// NewMinIOStorage creates a new MinIOStorage instance
// It returns a pointer to the instance
func NewMinIOStorage(c MinIOClient, bucketName string) *MinIOStorage {
	return &MinIOStorage{Client: c, BucketName: bucketName}
}

// Save saves a file to the minio storage
// It returns an error if the operation fails
// It takes a file path and data as input
func (m *MinIOStorage) Save(filePath string, data io.Reader) error {
	_, err := m.Client.PutObject(
		context.Background(),
		m.BucketName,
		filePath,
		data,
		-1,
		minio.PutObjectOptions{},
	)
	if err != nil {
		log.Printf("Error saving file to %s, %v", m.BucketName, err)
		return err
	}
	return nil
}

// Retrieve retrieves a file from the minio storage
// It returns a reader and an error
// It takes a file path as input
func (m *MinIOStorage) Retrieve(filePath string) (io.ReadCloser, error) {
	file, err := m.Client.GetObject(
		context.Background(),
		m.BucketName,
		filePath,
		minio.GetObjectOptions{},
	)
	if err != nil {
		log.Printf("Error retrieving file to %s, %v", m.BucketName, err)
		return nil, err
	}
	return file, nil
}

// Delete deletes a file from the minio storage
// It returns an error if the operation fails
// It takes a file path as input
func (m *MinIOStorage) Delete(filePath string) error {
	err := m.Client.RemoveObject(
		context.Background(),
		m.BucketName,
		filePath,
		minio.RemoveObjectOptions{})

	if err != nil {
		log.Printf("Error deleting file to %s, %v", m.BucketName, err)
		return err
	}

	return nil
}

// Exists checks if a file exists in the minio storage
// It returns a boolean and an error
// It takes a file path as input
func (m *MinIOStorage) Exists(filePath string) (bool, error) {
	_, err := m.Client.StatObject(
		context.Background(),
		m.BucketName,
		filePath,
		minio.StatObjectOptions{},
	)
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		log.Printf("Error checking if file exists %s, %v", m.BucketName, err)
		return false, err
	}
	return true, nil
}
