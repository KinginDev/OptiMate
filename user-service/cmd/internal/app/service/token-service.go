package service

import (
	"time"
	"user-service/cmd/internal/app/repositories"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JWTService struct {
	Repo      *repositories.JWTRepository
	SecretKey string
}

func NewJWTService(repo *repositories.JWTRepository, secretKey string) *JWTService {
	return &JWTService{
		Repo:      repo,
		SecretKey: secretKey,
	}
}

func (s *JWTService) GenerateToken(userID string) (string, error) {
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
