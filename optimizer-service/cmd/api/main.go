// Package main
package main

import (
	"optimizer-service/cmd/config"
	"optimizer-service/cmd/internal/app/handler"
	"optimizer-service/cmd/internal/app/repositories"
	"optimizer-service/cmd/internal/app/service"
	"optimizer-service/cmd/internal/types"
	"optimizer-service/cmd/internal/utils"
	"optimizer-service/cmd/lib/optimizer"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "optimizer-service/cmd/api/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	// Init database
	app := config.NewConfig()
	db := app.InitDB()
	storage := app.InitStorage()
	utils := utils.NewUtils(db)
	// Setup Repositories
	fileRepo := repositories.NewFileRepository(db)

	// Setup Services
	fileService := service.NewFileService(fileRepo, storage)

	//Setup Optimizer
	optimizer := optimizer.InitOptimizer(storage, utils)
	// Init App Container
	container := &types.AppContainer{
		DB:          db,
		Utils:       utils,
		FileService: fileService,
		Optimizer:   optimizer,
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
	e.GET("/docs/*", echoSwagger.WrapHandler)

	e.POST("/upload", h.PostUploadFile)

	optimizerServicePort := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + optimizerServicePort))
}
