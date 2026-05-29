package users

import "time"

// UserProfile represents extended user profile data.
type UserProfile struct {
	UserID            string `json:"user_id"`
	Bio               string `json:"bio"`
	Location          string `json:"location"`
	Website           string `json:"website"`
	Phone             string `json:"phone"`
	Timezone          string `json:"timezone"`
	Language          string `json:"language"`
	NotificationPrefs string `json:"notification_prefs"`
}

// UserWithProfile combines a user with their profile data.
type UserWithProfile struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Avatar    string     `json:"avatar"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	LastSeen  *time.Time `json:"last_seen_at,omitempty"`
	Profile   *UserProfile `json:"profile,omitempty"`
	Roles     []string   `json:"roles,omitempty"`
}

// ListUsersParams defines query parameters for listing users.
type ListUsersParams struct {
	Search       string `json:"search"`
	RoleFilter   string `json:"role_filter"`
	StatusFilter string `json:"status_filter"`
	Page         int    `json:"page"`
	PerPage      int    `json:"per_page"`
}

// PaginatedResponse wraps a paginated result set.
type PaginatedResponse struct {
	Items      []UserWithProfile `json:"items"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
	TotalPages int               `json:"total_pages"`
}

// UpdateProfileRequest is the DTO for updating a user profile.
type UpdateProfileRequest struct {
	Bio      string `json:"bio"`
	Location string `json:"location"`
	Website  string `json:"website"`
	Phone    string `json:"phone"`
	Timezone string `json:"timezone"`
	Language string `json:"language"`
}

// ChangePasswordRequest is the DTO for changing password.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// UpdateUserRequest is the DTO for admin updating a user.
type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	IsActive *bool  `json:"is_active,omitempty"`
}
