package repositories

import (
	"optimizer-service/cmd/internal/models"

	"gorm.io/gorm"
)

// IFileRepository is an interface for the file repository
type IFileRepository interface {
	CreateFile(file *models.File) error
}

// FileRepository is a struct for the file repository
// It implements the IFileRepository interface
type FileRepository struct {
	DB *gorm.DB
}

// NewFileRepository creates a new file repository
// It returns a pointer to the file repository
// It takes a gorm.DB as input
func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{DB: db}
}

// CreateFile creates a new file
// It takes a file as input
// GetFile retrieves a file by its ID
func (r *FileRepository) CreateFile(file *models.File) error {
	return r.DB.Create(file).Error
}

// GetFile retrieves a file by its ID
// It takes a file ID as input
// It returns a file and an error
func (r *FileRepository) GetFile(id string) (*models.File, error) {
	var file models.File
	result := r.DB.First(&file, id)
	return &file, result.Error
}
