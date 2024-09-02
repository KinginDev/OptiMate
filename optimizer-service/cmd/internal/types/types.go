package types

import (
	"optimizer-service/cmd/internal/app/service"
	"optimizer-service/cmd/internal/utils"

	"gorm.io/gorm"
)

type AppContainer struct {
	DB          *gorm.DB
	Utils       *utils.Utils
	FileService *service.FileService
}
