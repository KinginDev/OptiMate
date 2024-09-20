package storage

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() (*LocalStorage, func()) {
	basePath, _ := ioutil.TempDir("", "localstorage")
	storage := NewLocalStorage(basePath)
	return storage, func() {
		os.RemoveAll(basePath)
	}
}

func TestLocalStorage_Save_Success(t *testing.T) {
	storage, cleanup := setup()
	defer cleanup()

	// Test saving a file
	err := storage.Save("testfile.txt", strings.NewReader("test content"))
	assert.NoError(t, err)

	// Check if the file exists and content is correct
	content, err := ioutil.ReadFile(filepath.Join(storage.BasePath, "testfile.txt"))
	assert.NoError(t, err)
	assert.Equal(t, "test content", string(content))
}

func TestLocalStorage_Save_Failure(t *testing.T) {
	storage, cleanup := setup()
	defer cleanup()

	// Induce failure by using an invalid path
	err := storage.Save(string([]byte{0x00}), strings.NewReader("test content"))
	assert.Error(t, err)
}

func TestLocalStorage_Retrieve_Success(t *testing.T) {
	storage, cleanup := setup()
	defer cleanup()

	// Setup - create a file first
	filePath := filepath.Join(storage.BasePath, "testfile.txt")
	ioutil.WriteFile(filePath, []byte("test content"), 0644)

	// Test retrieving a file
	reader, err := storage.Retrieve("testfile.txt")
	assert.NoError(t, err)

	// Check the content of the file
	buf := new(strings.Builder)
	_, err = io.Copy(buf, reader)
	assert.NoError(t, err)
	assert.Equal(t, "test content", buf.String())
}

func TestLocalStorage_Retrieve_Failure(t *testing.T) {
	storage, cleanup := setup()
	defer cleanup()

	// Test retrieving a non-existent file
	_, err := storage.Retrieve("nonexistent.txt")
	assert.Error(t, err)
}

func TestLocalStorage_Delete_Success(t *testing.T) {
	storage, cleanup := setup()
	defer cleanup()

	// Setup - create a file first
	filePath := filepath.Join(storage.BasePath, "testfile.txt")
	ioutil.WriteFile(filePath, []byte("test content"), 0644)

	// Test deleting a file
	err := storage.Delete("testfile.txt")
	assert.NoError(t, err)

	// Check if the file still exists
	_, err = os.Stat(filePath)
	assert.True(t, os.IsNotExist(err))
}

func TestLocalStorage_Delete_Failure(t *testing.T) {
	storage, cleanup := setup()
	defer cleanup()

	// Test deleting a non-existent file
	err := storage.Delete("nonexistent.txt")
	assert.Error(t, err)
}

func TestLocalStorage_Exists(t *testing.T) {
	storage, cleanup := setup()
	defer cleanup()

	// Creating a file to test existence
	err := storage.Save("testfile.txt", strings.NewReader("test content"))
	assert.NoError(t, err)

	// Test checking if the newly created file exists
	exists, err := storage.Exists("testfile.txt")
	assert.NoError(t, err)
	assert.True(t, exists) // This should now pass

	// Test checking if a non-existent file exists
	exists, err = storage.Exists("nonexistent.txt")
	assert.NoError(t, err)
	assert.False(t, exists)
}
