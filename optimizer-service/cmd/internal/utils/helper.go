// Package utils
package utils

import (
	"log"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Utils struct {
	DB *gorm.DB
}

// JSONResponse is a type that defines the structure of a JSON response.
type JSONResponse struct {
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
}

func NewUtils(db *gorm.DB) *Utils {
	return &Utils{
		DB: db,
	}
}

func (u *Utils) WriteErrorResponse(c echo.Context, status int, message string) error {
	response := &JSONResponse{
		Data:    nil,
		Message: message,
		Success: false,
		Status:  status,
	}
	log.Println(response)
	return c.JSON(status, response)
}

func (u *Utils) WriteSuccessResponse(c echo.Context, status int, message string, data interface{}) error {
	response := &JSONResponse{
		Data:    data,
		Message: message,
		Success: true,
		Status:  status,
	}
	log.Println(response)
	return c.JSON(status, response)
}
