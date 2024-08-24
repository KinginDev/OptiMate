package main

import (
	"net/http"

	"user-service/config"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// create a response struct to add the response data, status, message,headers etc
type JsonResponse struct {
	Data    string `json:"data"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func main() {
	app := &config.Config{}
	app.InitDB()
	// models := data.New(db)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {

		response := &JsonResponse{
			Data:    "Test Data",
			Message: "Test Message",
			Status:  http.StatusOK,
		}

		return c.JSON(response.Status, response)
	})

	// Routes
	e.GET("/docs/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
