package storage

import "io"

type Storage interface {
	Save(filePath string, data io.Reader) error
	Retrive(filePath string) (io.ReadCloser, error)
	Delete(filePath string) error
}
