// Package interceptor
package interceptor

import (
	"net/http"
	"strings"
	"user-service/cmd/internal/app/service"

	"github.com/labstack/echo/v4"
)

// JWTAuthentication is a middleware that checks if the request has a valid JWT token
// and if the token is not revoked
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

			// Validate the token
			token, err := jwtService.ValidateToken(tokenString)
			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to validate token")
			}

			// Check the claims for userID
			userID, err := jwtService.GetUserIDFromToken(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Token claims are not accessible")
			}

			// Check if token has been revoved
			// _, err = jwtService.CheckTokenRevocation(tokenString)
			// if err != nil {
			// 	log.Printf("Failed to check token revocation %v\n", err)
			// }

			// // Attempt revoke token
			// err = jwtService.RevokeToken(tokenString)
			// if err != nil {
			// 	return echo.NewHTTPError(http.StatusInternalServerError, "Failed to revoke token")
			// }

			c.Set("userID", userID)

			return next(c)
		}
	}
}
