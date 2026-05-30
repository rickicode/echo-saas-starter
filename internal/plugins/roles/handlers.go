package roles

import (
	"net/http"
	"strconv"

	"echo-saas-starter/internal/plugins/auth"

	"github.com/labstack/echo/v4"
)

// Handlers holds dependencies for roles HTTP handlers.
type Handlers struct {
	repo *Repository
}

// NewHandlers creates a new roles handler set.
func NewHandlers(repo *Repository) *Handlers {
	return &Handlers{repo: repo}
}

// ListRoles returns all available roles.
func (h *Handlers) ListRoles(c echo.Context) error {
	roles, err := h.repo.GetAllRoles(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "ROLE_010",
				"message": "Failed to fetch roles",
			},
		})
	}
	return c.JSON(http.StatusOK, roles)
}

// AssignRoleHandler assigns a role to a user.
func (h *Handlers) AssignRoleHandler(c echo.Context) error {
	roleIDStr := c.Param("id")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "ROLE_011",
				"message": "Invalid role ID",
			},
		})
	}

	var req AssignRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "ROLE_012",
				"message": "Invalid request body",
			},
		})
	}

	if req.UserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "ROLE_013",
				"message": "user_id is required",
			},
		})
	}

	if err := h.repo.AssignRole(c.Request().Context(), req.UserID, roleID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "ROLE_014",
				"message": "Failed to assign role",
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// RemoveRoleHandler removes a role from a user.
func (h *Handlers) RemoveRoleHandler(c echo.Context) error {
	roleIDStr := c.Param("id")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "ROLE_020",
				"message": "Invalid role ID",
			},
		})
	}

	userID := c.Param("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "ROLE_021",
				"message": "user_id is required",
			},
		})
	}

	if err := h.repo.RemoveRole(c.Request().Context(), userID, roleID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "ROLE_022",
				"message": "Failed to remove role",
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetUserRolesHandler returns all roles for a user.
func (h *Handlers) GetUserRolesHandler(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "ROLE_030",
				"message": "user_id is required",
			},
		})
	}

	roles, err := h.repo.GetUserRoles(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "ROLE_031",
				"message": "Failed to fetch user roles",
			},
		})
	}

	return c.JSON(http.StatusOK, roles)
}

// myRoleEntry is the response shape for the /me endpoint.
type myRoleEntry struct {
	RoleID   int    `json:"role_id"`
	RoleName string `json:"role_name"`
}

// GetMyRolesHandler returns the authenticated user's roles.
func (h *Handlers) GetMyRolesHandler(c echo.Context) error {
	user := auth.GetUserFromContext(c)
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "ROLE_040",
				"message": "Authentication required",
			},
		})
	}

	roles, err := h.repo.GetUserRoles(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "ROLE_041",
				"message": "Failed to fetch user roles",
			},
		})
	}

	entries := make([]myRoleEntry, len(roles))
	for i, r := range roles {
		entries[i] = myRoleEntry{RoleID: r.ID, RoleName: r.Name}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"roles": entries,
	})
}
