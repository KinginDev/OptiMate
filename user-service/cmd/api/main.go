package main

import (
	"log"
	"user-service/cmd/api/handler"
	"user-service/cmd/api/interceptor"
	"user-service/cmd/api/validators"
	"user-service/cmd/config"
	"user-service/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	// Initilize app wide dependencies
	app := &config.Config{}
	db := app.InitDB()
	util := &utils.Utils{DB: db}

	if db == nil {
		log.Fatalf("Could not connect to the database.")
		return
	}

	// Create new handler instance with the db instance
	h := handler.NewHandler(db, util)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Add DB to echo context
	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// add validator to middleware
	e.Validator = &validators.CustomValidator{}

	e.GET("/", h.Index)
	e.POST("/register", h.Register)
	e.POST("/login", h.Login)
	// Docs Routes
	e.GET("/docs/*", echoSwagger.WrapHandler)

	authGroup := e.Group("profile")
	// Middleware
	authGroup.Use(interceptor.JWTAuthentication)
	authGroup.GET("/tokens", h.GetUserTokens)

	// Serve docs directly
	// e.GET("/docs/swagger.json", func(c echo.Context) error {
	// 	return c.File("./docs/swagger.json")
	// })

	e.Logger.Fatal(e.Start(":8080"))
}
