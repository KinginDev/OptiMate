package main

import (
	"os"
	"user-service/cmd/config"
	"user-service/cmd/internal/app/handler"
	"user-service/cmd/internal/app/repositories"
	"user-service/cmd/internal/app/service"
	"user-service/cmd/internal/interceptor"
	"user-service/cmd/internal/types"
	"user-service/cmd/internal/utils"
	"user-service/cmd/internal/validators"

	_ "user-service/cmd/api/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	// Initilize
	app := config.NewConfig()
	db := app.InitDB()

	// Set up the repositories and services
	userRepo := repositories.NewUserRepository(db)
	jwtRepo := repositories.NewJWTTokenRepository(db)
	userService := service.NewUserService(userRepo)
	jwtService := service.NewJWTService(jwtRepo, "secret")

	// Create a new container
	container := &types.AppContainer{
		Utils:       utils.NewUtils(db),
		DB:          db,
		UserService: userService,
		JWTService:  jwtService,
	}

	// Create new handler instance with the db instance
	h := handler.NewHandler(container)

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
	e.POST("/validate", h.ValidateUserToken)
	// Docs Routes
	e.GET("/docs/*", echoSwagger.WrapHandler)

	authGroup := e.Group("profile")
	// Middleware
	authGroup.Use(interceptor.JWTAuthentication(jwtService))

	authGroup.GET("/tokens", h.GetUserJWTTokens)

	userServicePort := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + userServicePort))
}

// swag init --generalInfo ./cmd/api/main.go  -o ./cmd/api/docs
