package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const contextKeyUser = "user"

// RequireAuth returns middleware that validates PASETO tokens and injects user into context.
func RequireAuth(repo *Repository, secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": map[string]interface{}{
						"code":    "AUTH_001",
						"message": "Missing authorization header",
					},
				})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": map[string]interface{}{
						"code":    "AUTH_002",
						"message": "Invalid authorization header format",
					},
				})
			}

			userID, err := ValidateAccessToken(parts[1], secret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": map[string]interface{}{
						"code":    "AUTH_003",
						"message": "Invalid or expired token",
					},
				})
			}

			user, err := repo.GetUserByID(c.Request().Context(), userID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": map[string]interface{}{
						"code":    "AUTH_004",
						"message": "User not found",
					},
				})
			}

			if !user.IsActive {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": map[string]interface{}{
						"code":    "AUTH_005",
						"message": "Account is deactivated",
					},
				})
			}

			c.Set(contextKeyUser, user)
			return next(c)
		}
	}
}

// GetUserFromContext retrieves the authenticated user from echo context.
func GetUserFromContext(c echo.Context) *User {
	val := c.Get(contextKeyUser)
	if val == nil {
		return nil
	}
	user, ok := val.(*User)
	if !ok {
		return nil
	}
	return user
}
