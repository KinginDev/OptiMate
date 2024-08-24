package main

import (
	"log"
	"user-service/cmd/api/handler"
	"user-service/cmd/api/validators"
	"user-service/cmd/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

var db *gorm.DB

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

	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	//add validator to middleware
	e.Validator = &validators.CustomValidator{}

	e.GET("/", h.Index)
	e.POST("/register", h.Register)
	e.POST("/login", h.Login)

	// Docs Routes
	e.GET("/docs/*", echoSwagger.WrapHandler)

	// Serve docs directly
	// e.GET("/docs/swagger.json", func(c echo.Context) error {
	// 	return c.File("./docs/swagger.json")
	// })

	e.Logger.Fatal(e.Start(":8080"))
}
