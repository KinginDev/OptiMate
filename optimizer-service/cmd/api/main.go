// Package main
package main

import (
	"optimizer-service/cmd/config"
	"optimizer-service/cmd/internal/app/handler"
	"optimizer-service/cmd/internal/app/repositories"
	"optimizer-service/cmd/internal/app/service"
	"optimizer-service/cmd/internal/types"
	"optimizer-service/cmd/internal/utils"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Init database
	app := config.NewConfig()
	db := app.InitDB()
	storage := app.InitStorage()

	// Setup Repositories
	fileRepo := repositories.NewFileRepository(db)

	// Setup Services
	storagePath := "./storage/uploads"

	//if the dir has not been created create it
	if _, err := os.Stat(storagePath); err != nil {
		os.MkdirAll(storagePath, os.ModePerm)
	}

	fileService := service.NewFileService(fileRepo, storage)

	// Init App Container
	container := &types.AppContainer{
		DB:          db,
		Utils:       utils.NewUtils(db),
		FileService: fileService,
	}

	// Start a new handle
	h := handler.NewHandler(container)

	//Set up echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.GET("/", h.HomePage)
	e.POST("/upload", h.PostUploadFile)

	optimizerServicePort := os.Getenv("OPTIMIZER_SERVICE_PORT")
	e.Logger.Fatal(e.Start(":" + optimizerServicePort))
}
