package service

import (
	"io"
	"log"
	"optimizer-service/cmd/internal/app/repositories"
	"optimizer-service/cmd/internal/models"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// IFileService is an interface for the file service
// It defines the methods that the file service should implement
type IFileService interface {
	UploadFile(userID string, fileData io.Reader, fileName string) (*models.File, error)
}

type FileService struct {
	Repo        *repositories.FileRepository
	StoragePath string
}

func NewFileService(r *repositories.FileRepository, storagePath string) *FileService {
	return &FileService{
		Repo:        r,
		StoragePath: storagePath,
	}
}

func (s *FileService) UploadFile(userId string, fileData io.Reader, fileName string) (*models.File, error) {
	// Construct file path
	targetPath := filepath.Join(s.StoragePath, fileName)

	// Set the output file
	outFile, err := os.Create(targetPath)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer outFile.Close()

	//Copy the file data to the target UploadFile
	bytesWritten, err := io.Copy(outFile, fileData)

	if err != nil {
		return nil, err
	}

	// Ensure the file chages are written to disk
	err = outFile.Sync()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Create file metadata
	file := &models.File{
		ID:           uuid.New().String(),
		UserID:       userId,
		OriginalName: fileName,
		OriginalPath: targetPath,
		Type:         filepath.Ext(fileName),
		Status:       models.StatusUploaded,
		Size:         bytesWritten,
	}

	err = s.Repo.CreateFile(file)
	if err != nil {
		return nil, err
	}

	return file, nil
}
