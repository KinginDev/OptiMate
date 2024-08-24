// Package middleware
package interceptor

import (
	"fmt"
	"net/http"
	"strings"
	"user-service/utils"

	"github.com/labstack/echo/v4"
)

// JWTAuthentication Middleware function to handle JWT authentication
func JWTAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		app := &utils.Config{}
		authorizationHeader := c.Request().Header.Get("Authorization")
		if authorizationHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
		}

		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")

		userID, err := app.ValidateJWTToken(tokenString)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to validate token")
		}

		c.Set("userID", userID)
		fmt.Println("User ID: ", userID)

		return next(c)
	}
}
