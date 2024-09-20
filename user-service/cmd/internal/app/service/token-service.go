package service

import (
	"errors"
	"time"
	"user-service/cmd/internal/app/repositories"
	"user-service/cmd/internal/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// JWTService is a service that handles JWT token generation, validation, and revocation
type JWTService struct {
	Repo      *repositories.JWTRepository
	SecretKey string
}

// NewJWTService creates a new instance of JWTService
func NewJWTService(repo *repositories.JWTRepository, secretKey string) *JWTService {
	return &JWTService{
		Repo:      repo,
		SecretKey: secretKey,
	}
}

// GenerateJWTToken generates a new JWT token
// It returns a token string and an error if the operation fails
func (s *JWTService) GenerateJWTToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token
// It returns a token and an error if the operation fails
func (s *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
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

// GetUserTokens retrieves all tokens for a user
// It returns a list of tokens and an error
func (s *UserService) GetUserTokens(id string) ([]models.PersonalToken, error) {
	return s.Repo.GetTokensByUserID(id)
}

// StoreToken stores a token in the database
// It returns an error if the operation fails
func (s *JWTService) StoreToken(token *models.PersonalToken) error {
	return s.Repo.StoreToken(token)
}

// RevokeToken revokes a token by setting the revoked field to true
// It returns an error if the operation fails
func (s *JWTService) RevokeToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return nil, nil // don't need the key for just checking expiration
	})

	// Error handling, ignoring signature validation since we provide nil key
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return errors.New("token is already expired")
			}
		}
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Check if the token is already expired
		if exp, ok := claims["exp"].(float64); ok && time.Unix(int64(exp), 0).Before(time.Now()) {
			return errors.New("token is already expired")
		}
	}

	// Proceed to revoke the token if it's not expired
	return s.Repo.RevokeToken(tokenString)
}

// CheckTokenRevocation checks if a token has been revoked
// It returns a boolean and an error
func (s *JWTService) CheckTokenRevocation(token string) (bool, error) {
	return s.Repo.CheckTokenRevocation(token)
}

// GetUserIDFromToken retrieves the user ID from a JWT token
// It returns the user ID and an error if the operation fails
func (s *JWTService) GetUserIDFromToken(tokenString string) (string, error) {
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", echo.NewHTTPError(401, "Token claims are not accessible")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", echo.NewHTTPError(401, "User ID not found in token claims")
	}

	return userID, nil
}
