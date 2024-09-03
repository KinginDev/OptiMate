package repositories

import (
	"optimizer-service/cmd/internal/models"

	"gorm.io/gorm"
)

type FileRepository struct {
	DB *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{DB: db}
}

func (r *FileRepository) CreateFile(file *models.File) error {
	return r.DB.Create(file).Error
}

func (r *FileRepository) GetFile(id string) (*models.File, error) {
	var file models.File
	result := r.DB.First(&file, id)
	return &file, result.Error
}
