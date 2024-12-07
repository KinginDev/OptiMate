// Package service
package service

import (
	"fmt"
	"optimizer-service/cmd/internal/app/interfaces"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// AuthService is a service for the auth repository
// It defines the methods that the auth service should implement
type AuthService struct {
	Repo interfaces.IAuthRepository
}

// NewAuthService creates a new auth service
// It returns a new auth service
func NewAuthService(r interfaces.IAuthRepository) *AuthService {
	return &AuthService{
		Repo: r,
	}
}

// Login logs in a user
// It returns the token if the login is successful
func (a *AuthService) Login(email, password string) (interface{}, error) {
	return a.Repo.LoginWithREST(email, password)
}

// ValidateToken validates the JWT token
// It returns the token if it is valid
func (a *AuthService) ValidateToken(token string) (*jwt.Token, error) {
	fmt.Printf("Validating token: %s\n", token)

	validatedToken, err := validatesJWTToken(token)
	if err != nil {
		return nil, err
	}

	return validatedToken, nil
}

// validatesJWTToken validates the JWT token
// It returns the token if it is valid
func validatesJWTToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.NewHTTPError(401, "Invalid token")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
