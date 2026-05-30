package roles

import (
	"echo-saas-starter/internal/plugins/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RequireRole returns middleware that checks if the user has one of the specified roles.
// This middleware must be used AFTER auth.RequireAuth.
func RequireRole(repo *Repository, roleNames ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := auth.GetUserFromContext(c)
			if user == nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": map[string]interface{}{
						"code":    "ROLE_001",
						"message": "Authentication required",
					},
				})
			}

			for _, roleName := range roleNames {
				hasRole, err := repo.HasRole(c.Request().Context(), user.ID, roleName)
				if err != nil {
					continue
				}
				if hasRole {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "ROLE_002",
					"message": "Insufficient permissions",
				},
			})
		}
	}
}
