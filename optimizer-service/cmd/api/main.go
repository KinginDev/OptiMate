// Package main
package main

import (
	"optimizer-service/cmd/config"
	"optimizer-service/cmd/internal/app/handler"
	"optimizer-service/cmd/internal/types"
	"optimizer-service/cmd/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	//Init database
	app := config.NewConfig()
	db := app.InitDB()

	//Init App Container
	container := &types.AppContainer{
		DB:    db,
		Utils: utils.NewUtils(db),
	}

	//Start a new handle
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

	e.GET("/", h.Index)
	e.Logger.Fatal(e.Start(":8021"))
}
