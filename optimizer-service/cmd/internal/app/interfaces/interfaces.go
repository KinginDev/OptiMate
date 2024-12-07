package interfaces

import (
	"io"
	"optimizer-service/cmd/internal/models"

	"github.com/golang-jwt/jwt"
)

// IAuthRepository is an interface for the auth repository
// It defines the methods that the auth repository should implement
type IAuthRepository interface {
	LoginWithREST(email, password string) (interface{}, error)
	ValidateToken(token string) (interface{}, error)
}

// IAuthService is an interface for the auth service
// It defines the methods that the auth service should implement
type IAuthService interface {
	Login(email string, password string) (interface{}, error)
	ValidateToken(token string) (*jwt.Token, error)
}

// IFileService is an interface for the file service
// It defines the methods that the file service should implement
type IFileService interface {
	UploadFile(userID string, fileData io.Reader, fileName string) (*models.File, error)
}

// IFileRepository is an interface for the file repository
type IFileRepository interface {
	CreateFile(file *models.File) error
}
