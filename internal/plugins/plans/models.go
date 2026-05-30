package plans

import "time"

// Plan represents a subscription plan.
type Plan struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description"`
	PriceMonthly float64   `json:"price_monthly"`
	PriceYearly  float64   `json:"price_yearly"`
	Features     string    `json:"features"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserPlan represents a user's plan assignment.
type UserPlan struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	PlanID    int        `json:"plan_id"`
	Status    string     `json:"status"`
	StartsAt  time.Time  `json:"starts_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// CreatePlanRequest is the DTO for creating a plan.
type CreatePlanRequest struct {
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
	Description  string  `json:"description"`
	PriceMonthly float64 `json:"price_monthly"`
	PriceYearly  float64 `json:"price_yearly"`
	Features     string  `json:"features"`
	IsActive     bool    `json:"is_active"`
}

// UpdatePlanRequest is the DTO for updating a plan.
type UpdatePlanRequest struct {
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
	Description  string  `json:"description"`
	PriceMonthly float64 `json:"price_monthly"`
	PriceYearly  float64 `json:"price_yearly"`
	Features     string  `json:"features"`
	IsActive     bool    `json:"is_active"`
}

// AssignPlanRequest is the DTO for assigning a plan to a user.
type AssignPlanRequest struct {
	UserID string `json:"user_id"`
}

// UserPlanWithDetails includes plan details.
type UserPlanWithDetails struct {
	UserPlan UserPlan `json:"user_plan"`
	Plan     Plan     `json:"plan"`
}
