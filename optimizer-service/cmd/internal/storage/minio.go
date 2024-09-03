package storage

import (
	"context"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
)

type MinIOStorage struct {
	Client     *minio.Client
	BucketName string
}

func NewMinIOStorage(c *minio.Client, bucketName string) *MinIOStorage {
	return &MinIOStorage{Client: c, BucketName: bucketName}
}

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
		log.Printf("Error saving file to %v, %v", m.BucketName, err)
		return err
	}
	return nil
}

func (m *MinIOStorage) Retrieve(filePath string) (io.ReadCloser, error) {
	file, err := m.Client.GetObject(
		context.Background(),
		m.BucketName,
		filePath,
		minio.GetObjectOptions{},
	)
	if err != nil {
		log.Printf("Error retrieving file to %v, %v", m.BucketName, err)
		return nil, err
	}
	return file, nil
}

func (m *MinIOStorage) Delete(filePath string) error {
	err := m.Client.RemoveObject(
		context.Background(),
		m.BucketName,
		filePath,
		minio.RemoveObjectOptions{})

	if err != nil {
		log.Printf("Error deleting file to %v, %v", m.BucketName, err)
		return err
	}

	return nil
}

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
		log.Printf("Error checking if file exists %v, %v", m.BucketName, err)
		return false, err
	}
	return true, nil
}
