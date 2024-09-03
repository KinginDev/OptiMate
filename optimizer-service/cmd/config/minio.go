package config

import (
	"log"

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
func NewMinioClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool) *minio.Client {
	if endpoint == "" {
		endpoint = "optimate_minio:9000"
	}
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Printf("Failed to connect to Minio %v", err)
		panic(err)
	}
	log.Printf("Connected to Minio")
	return minioClient
}

// Getters

func (m *MinioConfig) GetEndpoint() string {
	return m.Endpoint
}

func (m *MinioConfig) GetAccessKeyID() string {
	return m.AccessKeyID
}

func (m *MinioConfig) GetSecretAccessKey() string {
	return m.SecretAccessKey
}

func (m *MinioConfig) GetUseSSL() bool {
	return m.UseSSL
}
