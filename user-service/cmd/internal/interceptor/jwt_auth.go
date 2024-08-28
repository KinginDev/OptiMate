// Package interceptor
package interceptor

import (
	"net/http"
	"strings"
	"user-service/cmd/internal/app/service"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// JWTAuthentication Middleware function to handle JWT authentication
func JWTAuthentication(jwtService *service.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorizationHeader := c.Request().Header.Get("Authorization")
			if authorizationHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
			}

			tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")

			if tokenString == authorizationHeader {
				return echo.NewHTTPError(http.StatusUnauthorized, "Bearer token not found")
			}

			token, err := jwtService.ValidateToken(tokenString)
			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to validate token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, "Token claims are not accessible")
			}

			userID, ok := claims["user_id"]

			c.Set("userID", userID)

			return next(c)
		}
	}
}
