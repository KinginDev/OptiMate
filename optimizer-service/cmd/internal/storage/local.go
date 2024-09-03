package storage

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// LocalStorage is a storage implementation that saves files to the local filesystem
type LocalStorage struct {
	BasePath string
}

// NewLocalStorage creates a new LocalStorage instance
// It takes a base path as input
// It returns a pointer to the instance
func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{BasePath: basePath}
}

// Save saves a file to the local storage
// It returns an error if the operation fails
// It takes a file path and data as input
func (l *LocalStorage) Save(filePath string, data io.Reader) error {
	fullPath := filepath.Join(l.BasePath, filePath)
	file, err := os.Create(fullPath) // Make sure fullPath includes the filename
	if err != nil {
		log.Printf("Error saving file to %s: %v", fullPath, err)
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	return err
}

// Retrieve retrieves a file from the local storage
// It returns a reader and an error
// It takes a file path as input
func (l *LocalStorage) Retrieve(filePath string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(l.BasePath, filePath))
}

// Delete deletes a file from the local storage
// It returns an error if the operation fails
// It takes a file path as input
func (l *LocalStorage) Delete(filePath string) error {
	fullPath := filepath.Join(l.BasePath, filePath)
	return os.Remove(fullPath)
}

// Exists checks if a file exists in the local storage
// It returns a boolean and an error
// It takes a file path as input
func (l *LocalStorage) Exists(filePath string) (bool, error) {
	fullPath := filepath.Join(l.BasePath, filePath)
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		log.Printf("Error checking file to %s, %v", l.BasePath, err)
		return false, nil
	}
	return true, err
}
