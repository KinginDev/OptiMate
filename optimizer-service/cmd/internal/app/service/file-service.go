package service

import (
	"io"
	"log"
	"optimizer-service/cmd/internal/app/repositories"
	"optimizer-service/cmd/internal/models"
	"optimizer-service/cmd/internal/storage"
	"path/filepath"

	"github.com/google/uuid"
)

// IFileService is an interface for the file service
// It defines the methods that the file service should implement
type IFileService interface {
	UploadFile(userID string, fileData io.Reader, fileName string) (*models.File, error)
}

type FileService struct {
	Repo    *repositories.FileRepository
	Storage storage.Storage
}

func NewFileService(r *repositories.FileRepository, storage storage.Storage) *FileService {
	return &FileService{
		Repo:    r,
		Storage: storage,
	}
}

func (s *FileService) UploadFile(userId string, fileData io.Reader, fileName string) (*models.File, error) {
	// Construct file path
	uniqueFileName := uuid.New().String() + filepath.Ext(fileName)
	targetPath := filepath.Join("/", uniqueFileName)
	//save the file
	err := s.Storage.Save(targetPath, fileData)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	info, err := s.Storage.Retrieve(targetPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer info.Close()

	// Retrieve the file size if the storage system doesn't provide it directly
	// This step might need adjustments depending on the storage implementation
	fileSize, err := io.Copy(io.Discard, info)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Create file metadata
	file := &models.File{
		ID:           uuid.New().String(),
		UserID:       userId,
		OriginalName: uniqueFileName,
		OriginalPath: targetPath,
		Type:         filepath.Ext(fileName),
		Status:       models.StatusUploaded,
		Size:         fileSize,
	}

	err = s.Repo.CreateFile(file)
	if err != nil {
		return nil, err
	}

	return file, nil
}
