// Package utils
package utils

import (
	"bytes"
	"fmt"
	"image"
	"log"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type IUtils interface {
	WriteErrorResponse(c echo.Context, status int, message string) error
	WriteSuccessResponse(c echo.Context, status int, message string, data interface{}) error
	CheckFileType(f []byte) (string, error)
}

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

func (u *Utils) CheckFileType(f []byte) (string, error) {
	// Decode the file to get the file type
	_, format, err := image.Decode(bytes.NewReader(f))
	if err != nil {
		log.Printf("Error decoding the image:- %v", err)
		return "", err
	}

	switch format {
	case "jpeg", "jpg":
		return "jpeg", nil
	case "png":
		return "png", nil
	case "webp":
		return "webp", nil
	case "gif":
		return "gif", nil
	default:
		return "", fmt.Errorf("unsupported file type")
	}

}
