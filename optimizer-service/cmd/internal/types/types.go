package types

import (
	"optimizer-service/cmd/internal/app/service"
	"optimizer-service/cmd/internal/storage"
	"optimizer-service/cmd/internal/utils"

	"gorm.io/gorm"
)

type AppContainer struct {
	DB          *gorm.DB
	Storage     storage.Storage
	Utils       utils.IUtils
	FileService service.IFileService // interface
}
