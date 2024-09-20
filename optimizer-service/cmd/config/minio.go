package config

import (
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

// Constructor
func NewMinioClient(endpoint, rootUser, rootPassword string, useSSL bool) *minio.Client {

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(rootUser, rootPassword, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Printf("Failed to connect to Minio %v", err)
	}
	log.Printf("Connected to Minio")
	return minioClient
}

// GetEndpoint returns the endpoint
// It takes no input
func (m *MinioConfig) GetEndpoint() string {
	return m.Endpoint
}

// GetAccessKeyID returns the access key ID
// It takes no input
func (m *MinioConfig) GetAccessKeyID() string {
	return m.AccessKeyID
}

// GetSecretAccessKey returns the secret access key
func (m *MinioConfig) GetSecretAccessKey() string {
	return m.SecretAccessKey
}

// GetUseSSL returns the use SSL
func (m *MinioConfig) GetUseSSL() bool {
	env := os.Getenv("ENV")
	switch env {
	case "production":
		return true
	default:
		return false
	}
}

//swag init --generalInfo ./cmd/api/main.go  --exclude vendor
