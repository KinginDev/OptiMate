// Package utils
package utils

import (
	"fmt"
	"time"
	"user-service/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Config struct {
	DB *gorm.DB
}

// JSONResponse is a type that defines the structure of a JSON response.
type JSONResponse struct {
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
}

func (app *Config) WriteErrorResponse(c echo.Context, status int, message string) error {
	response := &JSONResponse{
		Data:    nil,
		Message: message,
		Success: false,
		Status:  status,
	}
	return c.JSON(status, response)
}

func (app *Config) WriteSuccessResponse(c echo.Context, status int, message string, data interface{}) error {
	response := &JSONResponse{
		Data:    data,
		Message: message,
		Success: true,
		Status:  status,
	}
	return c.JSON(status, response)
}

func (app *Config) GenerateJWTToken(userID string) (string, error) {
	// Create a new JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Invalid token")
	}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (app *Config) ValidateJWTToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.NewHTTPError(401, "Invalid token")
		}
		return []byte("secret"), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("Invalid token")
	}

	// Set the userID in the context
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", errors.New("Failed to validate token")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("Failed to validate token")
	}

	// remove expired tokens from the database
	exp := int64(claims["exp"].(float64))
	if exp < time.Now().Unix() {
		// Token is expired, delete it
		app.DB.Where("token = ?", tokenString).Delete(&models.PersonalToken{})
		fmt.Println("token has expired")
	}

	return userID, nil
}
