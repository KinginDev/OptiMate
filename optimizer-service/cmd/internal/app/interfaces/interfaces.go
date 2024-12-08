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
	UpdateFile(file *models.File) error
}

// OptimizerParams is a struct that defines the parameters for the optimizer.
type OptimizerParams struct {
	Level      *string
	Path       *string
	Name       *string
	Size       *int64
	CropParams *CropParams
}

type CropParams struct {
	X, Y, Width, Height int
}

type IOptimizer interface {
	Optimize(filePath string, file *models.File, oParam *OptimizerParams) (io.ReadCloser, error)
	SupportedFormats() []string
}
