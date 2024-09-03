package models

import (
	"time"
)

type FileStatus string

const (
	StatusPending     FileStatus = "pending"
	StatusUploaded    FileStatus = "uploaded"
	StatusProcessing  FileStatus = "processing"
	StatusCompleleted FileStatus = "completed"
	StatusFailed      FileStatus = "failed"
)

type File struct {
	ID                string     `json:"id" gorm:"type:uuid;primary_key"`
	UserID            string     `json:"user_id" gorm:"type:uuid;not null"`
	OriginalName      string     `json:"original_name" gorm:"type:varchar(255);not null"`
	OptimizedPath     *string    `json:"optimized_path" gorm:"type:varchar(255)"`
	OptimizedName     *string    `json:"optimized_name" gorm:"type:varchar(255)"`
	OptimizedSize     *int64     `json:"optimized_size" gorm:"type:bigint"`
	OptimizationLevel *string    `json:"optimization_level" gorm:"type:varchar(255)"`
	Size              int64      `json:"size" gorm:"not null"`
	OriginalPath      string     `json:"original_path" gorm:"type:varchar(255);not null"`
	Type              string     `json:"type" gorm:"type:varchar(255);not null"`
	Status            FileStatus `json:"status" gorm:"type:varchar(255);not null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

type OptimizationSettings struct {
	ID                string `json:"id" gorm:"type:uuid;primary_key"`
	FileType          string `json:"file_type" gorm:"type:varchar(255)"`
	OptimizationLevel string `json:"optimization_level" gorm:"type:varchar(255)"`
	Description       string `json:"description" gorm:"type:varchar(255)"`
	SettingsDetails   string `json:"settings_details" gorm:"type:text"`
}
