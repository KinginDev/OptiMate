// Package interceptor
package interceptor

import (
	"net/http"
	"strings"
	"user-service/cmd/internal/app/service"

	"github.com/labstack/echo/v4"
)

// JWTAuthentication godoc
// @Summary JWT Authentication middleware
// @Description JWT Authentication middleware
// @Param Authorization header string true "
// @Tags JWT
// @Accept json
// @Produce json
// @Success 200 {object} string "success"
// @Failure 401 {object} string "Unauthorized"
// JWTAuthentication is a middleware function that validates the JWT token
// and sets the user ID in the echo context
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
			revoked, err := jwtService.CheckTokenRevocation(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check token revocation")
			}

			if revoked {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token has been revoked")
			}

			// Attempt revoke token
			err = jwtService.RevokeToken(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to revoke token")
			}
			c.Set("userID", userID)

			return next(c)
		}
	}
}
