package storage

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{BasePath: basePath}
}

func (l *LocalStorage) Save(filePath string, data io.Reader) error {

	file, err := os.Create(filepath.Join(l.BasePath, filePath))
	if err != nil {
		log.Printf("Error saving file to %v, %v", l.BasePath, err)
		return err
	}
	defer file.Close()

	// Ensure the file changes are written to disk
	_, err = io.Copy(file, data)
	return err
}

func (l *LocalStorage) Retrieve(filePath string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(l.BasePath, filePath))
}

func (l *LocalStorage) Delete(filePath string) error {
	return os.Remove(l.BasePath + filePath)
}

func (l *LocalStorage) Exists(filePath string) (bool, error) {
	_, err := os.Stat(l.BasePath + filePath)
	if os.IsNotExist(err) {
		log.Printf("Error checking file to %v, %v", l.BasePath, err)
		return false, nil
	}
	return true, err
}
