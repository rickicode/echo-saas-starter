package plans

import (
	"echo-saas-starter/internal/plugins/auth"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Handlers holds dependencies for plan HTTP handlers.
type Handlers struct {
	repo *Repository
}

// NewHandlers creates a new plan handler set.
func NewHandlers(repo *Repository) *Handlers {
	return &Handlers{repo: repo}
}

// ListPlans handles GET /plans (admin).
func (h *Handlers) ListPlans(c echo.Context) error {
	plans, err := h.repo.GetAllPlans(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("PLN_001", "Failed to list plans"))
	}
	return c.JSON(http.StatusOK, plans)
}

// GetPlan handles GET /plans/:id (admin).
func (h *Handlers) GetPlan(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("PLN_002", "Invalid plan ID"))
	}

	plan, err := h.repo.GetPlanByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, errorResponse("PLN_003", "Plan not found"))
	}
	return c.JSON(http.StatusOK, plan)
}

// CreatePlan handles POST /plans (admin).
func (h *Handlers) CreatePlan(c echo.Context) error {
	var req CreatePlanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("PLN_004", "Invalid request body"))
	}

	if req.Name == "" || req.Slug == "" {
		return c.JSON(http.StatusBadRequest, errorResponse("PLN_005", "Name and slug are required"))
	}

	if err := h.repo.CreatePlan(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("PLN_006", "Failed to create plan"))
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Plan created"})
}

// UpdatePlan handles PUT /plans/:id (admin).
func (h *Handlers) UpdatePlan(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("PLN_002", "Invalid plan ID"))
	}

	var req UpdatePlanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("PLN_007", "Invalid request body"))
	}

	if err := h.repo.UpdatePlan(c.Request().Context(), id, req); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("PLN_008", "Failed to update plan"))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Plan updated"})
}

// DeletePlan handles DELETE /plans/:id (admin).
func (h *Handlers) DeletePlan(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("PLN_002", "Invalid plan ID"))
	}

	if err := h.repo.DeletePlan(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("PLN_009", "Failed to delete plan"))
	}

	return c.NoContent(http.StatusNoContent)
}

// AssignPlan handles POST /plans/:id/assign (admin).
func (h *Handlers) AssignPlan(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("PLN_002", "Invalid plan ID"))
	}

	var req AssignPlanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("PLN_010", "Invalid request body"))
	}

	if req.UserID == "" {
		return c.JSON(http.StatusBadRequest, errorResponse("PLN_011", "User ID is required"))
	}

	if err := h.repo.AssignPlanToUser(c.Request().Context(), req.UserID, id); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("PLN_012", "Failed to assign plan"))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Plan assigned"})
}

// GetCurrentPlan handles GET /plans/current (authenticated user).
func (h *Handlers) GetCurrentPlan(c echo.Context) error {
	user := auth.GetUserFromContext(c)
	if user == nil {
		return c.JSON(http.StatusUnauthorized, errorResponse("PLN_020", "Not authenticated"))
	}

	up, plan, err := h.repo.GetUserActivePlan(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"user_plan": nil,
			"plan":      nil,
		})
	}

	return c.JSON(http.StatusOK, UserPlanWithDetails{
		UserPlan: *up,
		Plan:     *plan,
	})
}

// GetAvailablePlans handles GET /plans/available (public).
func (h *Handlers) GetAvailablePlans(c echo.Context) error {
	plans, err := h.repo.GetActivePlans(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("PLN_030", "Failed to list plans"))
	}
	return c.JSON(http.StatusOK, plans)
}

func errorResponse(code, message string) map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
}
