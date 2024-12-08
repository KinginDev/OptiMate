// Package mocks contains the mocked objects for the file service
package mocks

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

type MockUtils struct {
	mock.Mock
}

// WriteSuccessResponse godoc
// @Summary WriteSuccessResponse writes a success response
// @Description WriteSuccessResponse is a mocked object that writes a success response
// @Tags utils
// @Accept json
// @Produce json
func (m *MockUtils) WriteSuccessResponse(c echo.Context, status int, message string, data interface{}) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(status)
	return json.NewEncoder(c.Response()).Encode(map[string]interface{}{
		"message": message,
		"data":    data,
	})
}

// WriteErrorResponse godoc
// @Summary WriteErrorResponse writes an error response
// @Description WriteErrorResponse is a mocked object that writes an error response
// @Tags utils
// @Accept json
// @Produce json
func (m *MockUtils) WriteErrorResponse(c echo.Context, status int, message string) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(status)
	return json.NewEncoder(c.Response()).Encode(map[string]interface{}{
		"message": message,
	})
}

// CheckFileType
// CheckFileType is a mocked object that checks the file type
func (m *MockUtils) CheckFileType(f []byte) (string, error) {
	args := m.Called(f)
	return args.String(0), args.Error(1)
}
