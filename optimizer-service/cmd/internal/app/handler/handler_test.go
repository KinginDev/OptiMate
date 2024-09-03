package handler

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"optimizer-service/cmd/internal/app/repositories"
	"optimizer-service/cmd/internal/app/service"
	"optimizer-service/cmd/internal/models"
	"optimizer-service/cmd/internal/types"
	"optimizer-service/cmd/internal/utils"
	"optimizer-service/cmd/lib/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setUpTest() (*echo.Echo, *types.AppContainer) {
	// Set up the test
	e := echo.New()

	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	//migrate models
	err := db.AutoMigrate(&models.File{}, &models.OptimizationSettings{})
	if err != nil {
		fmt.Println("Error migrating the schema")
		return nil, nil
	}

	fileRepo := repositories.NewFileRepository(db)
	fileService := service.NewFileService(fileRepo, "./storage/uploads")
	container := &types.AppContainer{
		Utils:       utils.NewUtils(db),
		DB:          db,
		FileService: fileService,
	}

	return e, container

}

func TestIndexSuccess(t *testing.T) {
	e, container := setUpTest()

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHandler(container)

	if assert.NoError(t, h.HomePage(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Index Welcome")
	}
}

func TestPostUploadFileWithSuccess(t *testing.T) {
	//create echo isnatnce
	e := echo.New()
	mockFileService := new(mocks.MockFileService)
	mockUtils := new(mocks.MockUtils)

	// Setup expected file
	expectedFile := &models.File{
		ID:           uuid.New().String(),
		UserID:       uuid.New().String(),
		OriginalName: "test.jpg",
		OriginalPath: "uploads/test.jpg",
		Size:         1000,
		Status:       "uploaded",
		Type:         ".jpg",
	}
	mockFileService.On("UploadFile", mock.Anything, mock.Anything, mock.Anything).Return(expectedFile, nil)
	mockUtils.On("WriteSuccessResponse", mock.Anything, http.StatusOK, "Successfully uploaded the file, optimization starting soon, you will get an email", mock.Anything).Return(nil)

	// Create a handler
	container := &types.AppContainer{
		Utils:       mockUtils,
		FileService: mockFileService,
	}

	handler := NewHandler(container)

	// Create a multipart form file
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("Dummy Data for test"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	// assert that the file was uploaded successfully
	if assert.NoError(t, handler.PostUploadFile(c)) {
		fmt.Println(rec.Body.String())
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Successfully uploaded the file, optimization starting soon, you will get an email")
	}

	mockFileService.AssertExpectations(t)
}

func TestPostUploadFileWithFailure(t *testing.T) {
	//create echo isnatnce
	e := echo.New()
	mockFileService := new(mocks.MockFileService)
	mockUtils := new(mocks.MockUtils)

	// Setup expected file
	expectedFile := &models.File{
		ID:           uuid.New().String(),
		UserID:       uuid.New().String(),
		OriginalName: "test.jpg",
		OriginalPath: "uploads/test.jpg",
		Size:         1000,
		Status:       "uploaded",
		Type:         ".jpg",
	}

	// Mock failed upload
	mockFileService.On("UploadFile", mock.Anything, mock.Anything, mock.Anything).Return(expectedFile, fmt.Errorf("error"))

	// Mock write error response
	mockUtils.On("WriteErrorResponse", mock.Anything, http.StatusBadRequest, "error").Return(nil)

	// Create a handler
	container := &types.AppContainer{
		Utils:       mockUtils,
		FileService: mockFileService,
	}

	handler := NewHandler(container)

	// Create a multipart form file
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("Dummy Data for test"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	// assert that the file was uploaded successfully
	if assert.NoError(t, handler.PostUploadFile(c)) {
		fmt.Println(rec.Body.String())
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "error")
	}

	mockFileService.AssertExpectations(t)
}
