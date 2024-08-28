// package: handler Test cases for the handler package
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"user-service/cmd/config"
	"user-service/cmd/internal/app/repositories"
	"user-service/cmd/internal/app/service"
	"user-service/cmd/internal/models"
	"user-service/cmd/internal/utils"
	"user-service/cmd/internal/validators"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Set up Echo and database for testing
func setUpTest() (*echo.Echo, *gorm.DB, *service.UserService) {
	e := echo.New()
	e.Validator = &validators.CustomValidator{}
	// use an in-memory SQLite database for testing
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	// Migrate the schema for the test database
	err := db.AutoMigrate(&models.User{}, &models.PersonalToken{})
	if err != nil {
		fmt.Println("Error migrating the schema")
		return nil, nil, nil
	}

	repo := repositories.NewUserRepository(db) // Create a new instance of UserRepository
	service := service.NewUserService(repo)

	return e, db, service
}

func TestIndexSuccess(t *testing.T) {
	e, db, service := setUpTest()

	config := &config.Config{
		DB: db,
	}

	// Create a new Http Request
	req := httptest.NewRequest(http.MethodPost, "/", nil)

	/**
		Creates a new ResponseRecorder.
		A ResponseRecorder is an implementation of the http.ResponseWriter interface.
		When passed to an HTTP handler, it records the response that the handler writes,
		 which allows you to inspect the response in your our tests.
	**/
	rec := httptest.NewRecorder()

	/**
		Create a new Echo Context.
		Can be retrieved later using `c.Get("config")`.
	**/
	c := e.NewContext(req, rec)
	/**
		Sets a value in the Echo context.
	**/
	c.Set("config", config)

	// Create a new handler instance
	utils := &utils.Utils{
		DB: db,
	}
	h := NewHandler(db, utils, service)

	// Execute the handler
	if assert.NoError(t, h.Index(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Welcome to the User Service!")
	}

}

func RegisterSuccessTest(t *testing.T) {
	e, db, service := setUpTest()

	config := &config.Config{
		DB: db,
	}

	// Create a new Http Request
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{
		"email": "admin@admin.com",
		"password": "password",
		"firstname": "John",
		"lastname": "Doe"

	}`))

	// set header to application/json
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	/**
		Creates a new ResponseRecorder.
		A ResponseRecorder is an implementation of the http.ResponseWriter interface.
		When passed to an HTTP handler, it records the response that the handler writes,
		which allows you to inspect the response in your our tests.
	**/
	rec := httptest.NewRecorder()

	/**
		Create a new Echo Context.
		Can be retrieved later using `c.Get("config")`.
	**/
	c := e.NewContext(req, rec)
	/**
		Sets a value in the Echo context.
	**/
	c.Set("config", config)

	// Create a new handler instance
	utils := &utils.Utils{
		DB: db,
	}
	h := NewHandler(db, utils, service)

	// Execute the handler
	if assert.NoError(t, h.Register(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "User Registered Successfully")
	}

}

func TestRegisterInvaidData(t *testing.T) {
	e, db, service := setUpTest()

	config := &config.Config{
		DB: db,
	}

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{
		"email": "",
		"password": ""
	}`))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("config", config)

	utils := &utils.Utils{
		DB: db,
	}
	h := NewHandler(db, utils, service)

	if assert.NoError(t, h.Register(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Email is required;Firstname is required;Lastname is required;Password is required")
	}

}

func TestRegisterWithExistingEmail(t *testing.T) {
	e, db, service := setUpTest()
	config := &config.Config{
		DB: db,
	}

	_, err := service.RegisterUser(
		&models.RegisterInput{
			Email:     "admin@admin.com",
			Password:  "password",
			Firstname: "John",
			LastName:  "Doe",
		},
	)

	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{
		"email": "admin@admin.com",
		"password": "password",
		"firstname": "John",
		"lastname": "Doe"
	}`))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("config", config)

	utils := &utils.Utils{
		DB: db,
	}
	h := NewHandler(db, utils, service)

	if assert.NoError(t, h.Register(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Email already exists")
	}
}

func TestLoginSuccess(t *testing.T) {
	e, db, service := setUpTest()

	config := &config.Config{
		DB: db,
	}

	_, err := service.RegisterUser(
		&models.RegisterInput{
			Email:     "admin@admin.com",
			Password:  "password",
			Firstname: "John",
			LastName:  "Doe",
		},
	)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{
		"email": "admin@admin.com",
		"password": "password"
	}`))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("config", config)

	utils := &utils.Utils{
		DB: db,
	}
	h := NewHandler(db, utils, service)

	if assert.NoError(t, h.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response JSONResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Login successful", response.Message)

	}

}

func TestLoginInvalidData(t *testing.T) {
	e, db, service := setUpTest()

	config := &config.Config{
		DB: db,
	}

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{
		"email": "",
		"password": ""
	}`))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("config", config)

	utils := &utils.Utils{
		DB: db,
	}
	h := NewHandler(db, utils, service)

	if assert.NoError(t, h.Login(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Email is required;Password is required")
	}

}
