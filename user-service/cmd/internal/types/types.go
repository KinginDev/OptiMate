// Package types
package types

import (
	"user-service/cmd/internal/app/service"
	"user-service/cmd/internal/utils"

	"gorm.io/gorm"
)

// AppContainer is a container for the application
type AppContainer struct {
	Utils       *utils.Utils
	DB          *gorm.DB
	UserService *service.UserService
	JWTService  *service.JWTService
}
