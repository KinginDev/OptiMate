// Package types
package types

import (
	"user-service/cmd/internal/app/service"
	"user-service/cmd/internal/utils"

	"gorm.io/gorm"
)

type AppContainer struct {
	Utils       *utils.Utils
	DB          *gorm.DB
	UserService *service.UserService
	JWTService  *service.JWTService
}
