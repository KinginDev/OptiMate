// Package utils
package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
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

func (u *Utils) WriteErrorResponse(c echo.Context, status int, message string) error {
	response := &JSONResponse{
		Data:    nil,
		Message: message,
		Success: false,
		Status:  status,
	}
	return c.JSON(status, response)
}

func (u *Utils) WriteSuccessResponse(c echo.Context, status int, message string, data interface{}) error {
	response := &JSONResponse{
		Data:    data,
		Message: message,
		Success: true,
		Status:  status,
	}
	return c.JSON(status, response)
}

func (u *Utils) GenerateJWTToken(userID string) (string, error) {
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

func (u *Utils) ValidateJWTToken(tokenString string) (string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.NewHTTPError(401, "Invalid token")
		}
		return []byte("secret"), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return "", nil
			}
		}
		return "", nil
	}

	if !token.Valid {
		return "", errors.New("Invalid token")
	}

	// Check if token has expired

	// Set the userID in the context
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", errors.New("Failed to validate token")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("Failed to validate token")
	}

	// Check if the token has been revoked
	// if err := u.checkTokenRevocation(tokenString); err != nil {
	// 	return "", err
	// }

	return userID, nil

}

// func (u *Utils) deleteExpiredToken(tokenString string) {
// 	if u.DB != nil {
// 		u.DB.Where("token = ?", tokenString).Delete(&models.PersonalToken{})
// 	}
// }

// func (u *Utils) checkTokenRevocation(tokenString string) error {
// 	if u.DB == nil {
// 		return errors.New("DB instance is not initialized")
// 	}

// 	var token models.PersonalToken
// 	result := u.DB.Where("token = ? AND revoked = ?", tokenString, true).First(&token)
// 	if result.Error != nil {
// 		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 			return nil
// 		}
// 		return result.Error
// 	}

// 	if token.Revoked {
// 		return errors.New("Token has been revoked")
// 	}

// 	return nil
// }
