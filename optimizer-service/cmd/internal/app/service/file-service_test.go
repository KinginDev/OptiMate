package service

import (
	"bytes"
	"errors"
	"io/ioutil"
	"optimizer-service/cmd/lib/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUploadFile_Success(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockRepo := new(mocks.MockFileRepository)
	fileService := NewFileService(mockRepo, mockStorage)

	// Setup mock expectations
	mockStorage.On("Save", mock.Anything, mock.Anything).Return(nil)
	mockStorage.On("Retrieve", mock.Anything).Return(ioutil.NopCloser(bytes.NewReader([]byte("file content"))), nil)
	mockRepo.On("CreateFile", mock.AnythingOfType("*models.File")).Return(nil)

	// Execute the method
	file, err := fileService.UploadFile("user123", bytes.NewReader([]byte("file data")), "testfile.txt")

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, "user123", file.UserID)
	mockStorage.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestUploadFile_Failure(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockRepo := new(mocks.MockFileRepository)
	fileService := NewFileService(mockRepo, mockStorage)

	// Setup failure scenario for storage Save
	mockStorage.On("Save", mock.Anything, mock.Anything).Return(errors.New("failed to save"))

	// Execute the method
	file, err := fileService.UploadFile("user123", bytes.NewReader([]byte("file data")), "testfile.txt")

	// Assert that an error was returned
	assert.Error(t, err)
	assert.Nil(t, file)
	assert.Equal(t, "failed to save", err.Error())
	mockStorage.AssertExpectations(t)
}
