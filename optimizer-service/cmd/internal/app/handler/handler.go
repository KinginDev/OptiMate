// Package handler
package handler

import (
	"fmt"
	"log"
	"net/http"
	"optimizer-service/cmd/internal/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handler struct to hold the db instance
type Handler struct {
	Container *types.AppContainer
}

// NewHandler function to initialize the handler with the given DB instance
func NewHandler(c *types.AppContainer) *Handler {
	return &Handler{Container: c}
}

// HomePage godoc
// @Summary Index Welcome
// @Description Index Welcome
// @Success 200 {object} utils.JSONResponse "success"
// @Router / [get]
func (h *Handler) HomePage(c echo.Context) error {
	response := map[string]interface{}{
		"message": "Index Welcome",
	}

	return h.Container.Utils.WriteSuccessResponse(c, http.StatusOK, "success", response)
}

// PostUploadFile godoc
// @Summary Upload a file
// @Description Upload a file
// @Accept mpfd
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} utils.JSONResponse "Successfully optimized the file"
// @Failure 400 {object} utils.JSONResponse "Error uploading file"
// @Router /upload [post]
func (h *Handler) PostUploadFile(c echo.Context) error {
	userId := uuid.New().String()

	// Get the submitted file
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error getting file from form data %v", err)
		return h.Container.Utils.WriteErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		log.Printf("Error opening the file %v", err)
		return h.Container.Utils.WriteErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Close the file at the end of the function
	defer src.Close()

	// Upload the file via the file service
	uploadedFile, err := h.Container.FileService.UploadFile(userId, src, file.Filename)
	if err != nil {
		log.Printf("Error uploading file %v", err)
		return h.Container.Utils.WriteErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	fmt.Printf("Uploaded file: %+v\n", uploadedFile)

	optimizedFile, err := h.Container.Optimizer.Optimize(uploadedFile.OriginalName)
	if err != nil {
		log.Printf("Error optimizing file %v", err)
	}

	fmt.Printf("Optimized file: %+v\n", optimizedFile)
	return h.Container.Utils.WriteSuccessResponse(c, http.StatusOK, "Successfully optimized the file", uploadedFile)
}

// LoginUser godoc
// @Summary Login a user
// @Description Login a user
// @Accept json
// @Produce json
// @Success 200 {object} utils.JSONResponse "Login successful"
// @Failure 400 {object} utils.JSONResponse "Invalid request payload"
// @Failure 401 {object} utils.JSONResponse "Invalid password"
// @Failure 404 {object} utils.JSONResponse "User not found"
// @Failure 500 {object} utils.JSONResponse "Failed to create token"
// @Router /login [post]
func (h *Handler) LoginUser(c echo.Context) error {
	u := new(types.LoginInput)

	if err := c.Bind(u); err != nil {
		return h.Container.Utils.WriteErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
	}
	email := u.Email
	password := u.Password

	// Call the auth service to login the user
	authResult, err := h.Container.AuthService.Login(email, password)
	if err != nil {
		return h.Container.Utils.WriteErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	return h.Container.Utils.WriteSuccessResponse(c, http.StatusOK, "Login successful", authResult)
}
