package roles

import (
	"context"
	"echo-saas-starter/internal/core"
)

// Repository handles database operations for roles.
type Repository struct {
	db core.DB
}

// NewRepository creates a new roles repository.
func NewRepository(db core.DB) *Repository {
	return &Repository{db: db}
}

// GetAllRoles returns all roles.
func (r *Repository) GetAllRoles(ctx context.Context) ([]Role, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, description, created_at FROM roles ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, role)
	}
	return result, rows.Err()
}

// GetUserRoles returns all roles for a user.
func (r *Repository) GetUserRoles(ctx context.Context, userID string) ([]Role, error) {
	rows, err := r.db.Query(ctx,
		`SELECT r.id, r.name, r.description, r.created_at
		 FROM roles r
		 INNER JOIN user_roles ur ON r.id = ur.role_id
		 WHERE ur.user_id = ?
		 ORDER BY r.id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, role)
	}
	return result, rows.Err()
}

// AssignRole assigns a role to a user.
func (r *Repository) AssignRole(ctx context.Context, userID string, roleID int) error {
	return r.db.Exec(ctx,
		`INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`, userID, roleID)
}

// RemoveRole removes a role from a user.
func (r *Repository) RemoveRole(ctx context.Context, userID string, roleID int) error {
	return r.db.Exec(ctx,
		`DELETE FROM user_roles WHERE user_id = ? AND role_id = ?`, userID, roleID)
}

// HasRole checks if a user has a specific role by name.
func (r *Repository) HasRole(ctx context.Context, userID string, roleName string) (bool, error) {
	row := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM user_roles ur
		 INNER JOIN roles ro ON ur.role_id = ro.id
		 WHERE ur.user_id = ? AND ro.name = ?`, userID, roleName)

	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
