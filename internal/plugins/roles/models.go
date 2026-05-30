package roles

import "time"

const (
	RoleSuperAdmin = 1
	RoleAdmin      = 2
	RoleUser       = 3
)

// Role represents a role in the system.
type Role struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserRole represents a user-role assignment.
type UserRole struct {
	UserID     string    `json:"user_id"`
	RoleID     int       `json:"role_id"`
	AssignedAt time.Time `json:"assigned_at"`
}

// AssignRoleRequest is the request body for assigning a role.
type AssignRoleRequest struct {
	UserID string `json:"user_id"`
}
