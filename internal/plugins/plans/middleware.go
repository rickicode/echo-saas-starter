package plans

import (
	"echo-saas-starter/internal/plugins/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RequirePlan returns middleware that checks if the user has an active plan matching one of the allowed slugs.
func RequirePlan(repo *Repository, slugs ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := auth.GetUserFromContext(c)
			if user == nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": map[string]interface{}{
						"code":    "PLN_MW_001",
						"message": "Authentication required",
					},
				})
			}

			_, plan, err := repo.GetUserActivePlan(c.Request().Context(), user.ID)
			if err != nil || plan == nil {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"error": map[string]interface{}{
						"code":    "PLN_MW_002",
						"message": "No active plan",
					},
				})
			}

			for _, slug := range slugs {
				if plan.Slug == slug {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "PLN_MW_003",
					"message": "Plan upgrade required",
				},
			})
		}
	}
}
