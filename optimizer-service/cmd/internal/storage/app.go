package storage

import "io"

type Storage interface {
	Save(filePath string, data io.Reader) error
	Retrieve(filePath string) (io.ReadCloser, error)
	Delete(filePath string) error
	Exists(filePath string) (bool, error)
}
