package utils

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Config struct {
	DB *gorm.DB
}

// create a response struct to add the response data, status, message,headers etc
type JsonResponse struct {
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
	Success bool        `json:"status"`
	Message string      `json:"message"`
}

func (app *Config) WriteErrorResponse(c echo.Context, status int, message string) error {
	response := &JsonResponse{
		Data:    nil,
		Message: message,
		Success: false,
		Status:  status,
	}
	return c.JSON(status, response)
}

func (app *Config) WriteSuccessResponse(c echo.Context, status int, message string, data interface{}) error {
	response := &JsonResponse{
		Data:    data,
		Message: message,
		Success: true,
		Status:  status,
	}
	return c.JSON(status, response)
}
