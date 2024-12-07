package types

import (
	"optimizer-service/cmd/internal/app/interfaces"
	"optimizer-service/cmd/internal/utils"
	"optimizer-service/cmd/lib/optimizer"

	"gorm.io/gorm"
)

type AppContainer struct {
	DB          *gorm.DB
	Utils       utils.IUtils
	FileService interfaces.IFileService // interface
	AuthService interfaces.IAuthService
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResponsePayload struct {
	Data      interface{}          `json:"data"`
	Optimizer optimizer.IOptimizer // Optimizer interface
}
