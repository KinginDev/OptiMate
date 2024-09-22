package types

import (
	"optimizer-service/cmd/internal/app/service"
	"optimizer-service/cmd/internal/utils"
	"optimizer-service/cmd/lib/optimizer"

	"gorm.io/gorm"
)

type AppContainer struct {
	DB          *gorm.DB
	Utils       utils.IUtils
	FileService service.IFileService // interface
	Optimizer   optimizer.IOptimizer // Optimizer interface
}
