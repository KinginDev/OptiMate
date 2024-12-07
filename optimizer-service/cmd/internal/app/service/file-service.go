package service

import (
	"io"
	"log"
	"optimizer-service/cmd/internal/app/interfaces"
	"optimizer-service/cmd/internal/models"
	"optimizer-service/cmd/internal/storage"
	"path/filepath"

	"github.com/google/uuid"
)

// FileService is a struct for the file service
// It implements the IFileService interface
type FileService struct {
	Repo    interfaces.IFileRepository
	Storage storage.Storage
}

// NewFileService creates a new file service
// It returns a pointer to the file service
func NewFileService(r interfaces.IFileRepository, storage storage.Storage) *FileService {
	return &FileService{
		Repo:    r,
		Storage: storage,
	}
}

// UploadFile uploads a file to the storage system
// It returns a file and an error
// It takes a userID, fileData and fileName as input
// It saves the file to the storage system and creates a file metadata
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
