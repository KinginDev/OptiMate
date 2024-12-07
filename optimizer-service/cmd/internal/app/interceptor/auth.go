// Package interceptor
package interceptor

import (
	"log"
	"net/http"
	"optimizer-service/cmd/internal/app/interfaces"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// AuthenticationMiddleware is a middleware that checks if the request has a valid JWT token
// and if the token is not revoked
// uses the ValidateToken method from the auth service to validate the token
func AuthenticationMiddleware(authService interfaces.IAuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				log.Println("Authorization header is required")
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
			}

			// Get the token for the Authorization header
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			if tokenString == authHeader {
				log.Println("Bearer token not found")
				return echo.NewHTTPError(http.StatusUnauthorized, "Bearer token not found")
			}

			// Check if token is valid using the auth service
			token, err := authService.ValidateToken(tokenString)

			if err != nil || !token.Valid {
				log.Printf("Failed to validate token: %v\n", err)
				return echo.NewHTTPError(http.StatusForbidden, "Invalid token")
			}

			// Check the claims for userID
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, "Token claims are not accessible")
			}
			userID, ok := claims["user_id"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, "Token claims are not accessible")
			}

			c.Set("userID", userID)
			return next(c)
		}
	}
}
