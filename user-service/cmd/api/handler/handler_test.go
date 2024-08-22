package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/cmd/config"
	"user-service/data"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Set up Echo and database for testing

func setUpTest() (*echo.Echo, *gorm.DB) {
	e := echo.New()

	// use an in-memory SQLite database for testing
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	//Migrate the schema for the test database
	db.AutoMigrate(&data.User{}, &data.PersonalToken{})

	return e, db
}

func TestIndexSuccess(t *testing.T) {
	e, db := setUpTest()

	config := &config.Config{
		Database: db,
	}

	//Create a new Http Request
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
	h := NewHandler(db)

	//Execute the handler
	if assert.NoError(t, h.Index(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Welcome to the User Service!")
	}

}
