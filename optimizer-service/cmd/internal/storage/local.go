package storage

import (
	"io"
	"log"
	"os"
)

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{BasePath: basePath}
}

func (l *LocalStorage) Save(filepath string, data io.Reader) error {
	fullPath := l.BasePath + filepath
	file, err := os.Create(fullPath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	// Ensure the file chages are written to disk
	_, err = io.Copy(file, data)
	return err
}

func (l *LocalStorage) Retrive(filePath string) (io.ReadCloser, error) {
	return os.Open(l.BasePath + filePath)
}

func (l *LocalStorage) Delete(filePath string) error {
	return os.Remove(l.BasePath + filePath)
}
