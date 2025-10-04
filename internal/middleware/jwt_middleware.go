package middleware

import (
	"net/http"
	"strings"

	"api-key-generator/internal/utils"

	"github.com/labstack/echo/v4"
)

// OptionalJWT parses Authorization: Bearer <token> if present and sets user_id in context.
func OptionalJWT() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
				token := strings.TrimSpace(auth[7:])
				if uid, err := utils.ParseJWT(token); err == nil && uid != "" {
					c.Set("user_id", uid)
				}
			}
			return next(c)
		}
	}
}

// RequireJWT enforces a valid JWT and sets user_id.
func RequireJWT() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "JWT required")
			}
			token := strings.TrimSpace(auth[7:])
			uid, err := utils.ParseJWT(token)
			if err != nil || uid == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}
			c.Set("user_id", uid)
			return next(c)
		}
	}
}
