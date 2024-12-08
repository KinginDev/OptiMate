package types

import (
	"optimizer-service/cmd/internal/app/interfaces"
	"optimizer-service/cmd/internal/utils"

	"gorm.io/gorm"
)

type AppContainer struct {
	DB          *gorm.DB
	Utils       utils.IUtils
	FileService interfaces.IFileService // interface
	AuthService interfaces.IAuthService
	Optimizer   interfaces.IOptimizer
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResponsePayload struct {
	Data interface{} `json:"data"`
}

type GenericInput[T any] struct {
	Data T `json:"data"`
}
