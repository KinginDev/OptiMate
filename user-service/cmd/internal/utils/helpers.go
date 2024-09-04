// Package utils
package utils

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Utils is a struct that contains utility functions.
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

// NewUtils creates a new instance of the Utils struct.
func NewUtils(db *gorm.DB) *Utils {
	return &Utils{
		DB: db,
	}
}

// WriteErrorResponse is a helper function to write an error response
// with the given status and message.
func (u *Utils) WriteErrorResponse(c echo.Context, status int, message string) error {
	response := &JSONResponse{
		Data:    nil,
		Message: message,
		Success: false,
		Status:  status,
	}
	return c.JSON(status, response)
}

// WriteSuccessResponse is a helper function to write a success response
// with the given status, message and data.
// It returns an error if the response cannot be written.
func (u *Utils) WriteSuccessResponse(c echo.Context, status int, message string, data interface{}) error {
	response := &JSONResponse{
		Data:    data,
		Message: message,
		Success: true,
		Status:  status,
	}
	return c.JSON(status, response)
}
