package main

import (
	"log"
	"user-service/cmd/api/handler"
	"user-service/cmd/api/validators"
	"user-service/cmd/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

var db *gorm.DB

// create a response struct to add the response data, status, message,headers etc
type JsonResponse struct {
	Data    string `json:"data"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func main() {
	app := &config.Config{}
	db := app.InitDB()

	if db == nil {
		log.Fatalf("Could not connect to the database.")
		return
	}

	// Create new handler instance with the db instance
	h := handler.NewHandler(db)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//add validator to middleware
	e.Validator = &validators.CustomValidator{}

	e.GET("/", h.Index)
	e.POST("/register", h.Register)
	e.POST("/login", h.Login)

	e.Logger.Fatal(e.Start(":8080"))
}
