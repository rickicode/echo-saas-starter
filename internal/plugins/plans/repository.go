package plans

import (
	"context"
	"echo-saas-starter/internal/core"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Repository handles database operations for plans.
type Repository struct {
	db core.DB
}

// NewRepository creates a new plans repository.
func NewRepository(db core.DB) *Repository {
	return &Repository{db: db}
}

// GetAllPlans returns all plans.
func (r *Repository) GetAllPlans(ctx context.Context) ([]Plan, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, slug, description, price_monthly, price_yearly, features, is_active, created_at, updated_at FROM plans ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Plan
	for rows.Next() {
		var p Plan
		if err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.PriceMonthly, &p.PriceYearly, &p.Features, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	if result == nil {
		result = []Plan{}
	}
	return result, rows.Err()
}

// GetActivePlans returns only active plans.
func (r *Repository) GetActivePlans(ctx context.Context) ([]Plan, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, slug, description, price_monthly, price_yearly, features, is_active, created_at, updated_at FROM plans WHERE is_active = 1 ORDER BY price_monthly ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Plan
	for rows.Next() {
		var p Plan
		if err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.PriceMonthly, &p.PriceYearly, &p.Features, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	if result == nil {
		result = []Plan{}
	}
	return result, rows.Err()
}

// GetPlanByID returns a plan by its ID.
func (r *Repository) GetPlanByID(ctx context.Context, id int) (*Plan, error) {
	row := r.db.QueryRow(ctx, `SELECT id, name, slug, description, price_monthly, price_yearly, features, is_active, created_at, updated_at FROM plans WHERE id = ?`, id)

	p := &Plan{}
	if err := row.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.PriceMonthly, &p.PriceYearly, &p.Features, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, err
	}
	return p, nil
}

// GetPlanBySlug returns a plan by slug.
func (r *Repository) GetPlanBySlug(ctx context.Context, slug string) (*Plan, error) {
	row := r.db.QueryRow(ctx, `SELECT id, name, slug, description, price_monthly, price_yearly, features, is_active, created_at, updated_at FROM plans WHERE slug = ?`, slug)

	p := &Plan{}
	if err := row.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.PriceMonthly, &p.PriceYearly, &p.Features, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, err
	}
	return p, nil
}

// CreatePlan creates a new plan.
func (r *Repository) CreatePlan(ctx context.Context, req CreatePlanRequest) error {
	now := time.Now().UTC()
	return r.db.Exec(ctx,
		`INSERT INTO plans (name, slug, description, price_monthly, price_yearly, features, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		req.Name, req.Slug, req.Description, req.PriceMonthly, req.PriceYearly, req.Features, req.IsActive, now, now)
}

// UpdatePlan updates an existing plan.
func (r *Repository) UpdatePlan(ctx context.Context, id int, req UpdatePlanRequest) error {
	return r.db.Exec(ctx,
		`UPDATE plans SET name = ?, slug = ?, description = ?, price_monthly = ?, price_yearly = ?, features = ?, is_active = ?, updated_at = ? WHERE id = ?`,
		req.Name, req.Slug, req.Description, req.PriceMonthly, req.PriceYearly, req.Features, req.IsActive, time.Now().UTC(), id)
}

// DeletePlan deletes a plan.
func (r *Repository) DeletePlan(ctx context.Context, id int) error {
	return r.db.Exec(ctx, `DELETE FROM plans WHERE id = ?`, id)
}

// AssignPlanToUser assigns a plan to a user, deactivating any current plan.
func (r *Repository) AssignPlanToUser(ctx context.Context, userID string, planID int) error {
	// Deactivate current plans
	if err := r.db.Exec(ctx, `UPDATE user_plans SET status = 'cancelled' WHERE user_id = ? AND status = 'active'`, userID); err != nil {
		return fmt.Errorf("failed to deactivate current plans: %w", err)
	}

	id := uuid.New().String()
	now := time.Now().UTC()
	return r.db.Exec(ctx,
		`INSERT INTO user_plans (id, user_id, plan_id, status, starts_at, created_at)
		 VALUES (?, ?, ?, 'active', ?, ?)`,
		id, userID, planID, now, now)
}

// GetUserActivePlan returns a user's currently active plan.
func (r *Repository) GetUserActivePlan(ctx context.Context, userID string) (*UserPlan, *Plan, error) {
	row := r.db.QueryRow(ctx,
		`SELECT up.id, up.user_id, up.plan_id, up.status, up.starts_at, up.expires_at, up.created_at,
		        p.id, p.name, p.slug, p.description, p.price_monthly, p.price_yearly, p.features, p.is_active, p.created_at, p.updated_at
		 FROM user_plans up
		 INNER JOIN plans p ON up.plan_id = p.id
		 WHERE up.user_id = ? AND up.status = 'active'
		 ORDER BY up.created_at DESC LIMIT 1`, userID)

	up := &UserPlan{}
	p := &Plan{}
	if err := row.Scan(&up.ID, &up.UserID, &up.PlanID, &up.Status, &up.StartsAt, &up.ExpiresAt, &up.CreatedAt,
		&p.ID, &p.Name, &p.Slug, &p.Description, &p.PriceMonthly, &p.PriceYearly, &p.Features, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, nil, err
	}
	return up, p, nil
}
