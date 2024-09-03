// Package handler
package handler

import (
	"log"
	"net/http"
	"optimizer-service/cmd/internal/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Container *types.AppContainer
}

func NewHandler(c *types.AppContainer) *Handler {
	return &Handler{Container: c}
}

func (h *Handler) HomePage(c echo.Context) error {
	response := map[string]interface{}{
		"message": "Index Welcome",
	}

	return h.Container.Utils.WriteSuccessResponse(c, http.StatusOK, "success", response)
}

func (h *Handler) PostUploadFile(c echo.Context) error {
	userId := uuid.New().String()

	// Get the submitted file
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error getting file from form data %v", err)
		return h.Container.Utils.WriteErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		log.Printf("Error opening the file %v", err)
		return h.Container.Utils.WriteErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	defer src.Close()

	//Uplaod file
	uploadedFile, err := h.Container.FileService.UploadFile(userId, src, file.Filename)
	if err != nil {
		log.Printf("Error uploading file %v", err)
		return h.Container.Utils.WriteErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return h.Container.Utils.WriteSuccessResponse(c, http.StatusOK, "Successfully uploaded the file, optimization starting soon, you will get an email", uploadedFile)
}
