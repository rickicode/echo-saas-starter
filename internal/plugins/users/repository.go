package users

import (
	"context"
	"echo-saas-starter/internal/core"
	"fmt"
	"time"
)

// Repository handles database operations for users management.
type Repository struct {
	db core.DB
}

// NewRepository creates a new users repository.
func NewRepository(db core.DB) *Repository {
	return &Repository{db: db}
}

// GetUsers returns a paginated list of users with optional search and filters.
func (r *Repository) GetUsers(ctx context.Context, params ListUsersParams) (*PaginatedResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PerPage < 1 {
		params.PerPage = 20
	}

	baseQuery := `FROM users u`
	where := ` WHERE 1=1`
	args := []interface{}{}

	if params.RoleFilter != "" {
		baseQuery += ` INNER JOIN user_roles ur ON u.id = ur.user_id INNER JOIN roles ro ON ur.role_id = ro.id`
		where += ` AND ro.name = ?`
		args = append(args, params.RoleFilter)
	}

	if params.Search != "" {
		where += ` AND (u.name LIKE ? OR u.email LIKE ?)`
		search := "%" + params.Search + "%"
		args = append(args, search, search)
	}

	if params.StatusFilter == "active" {
		where += ` AND u.is_active = 1`
	} else if params.StatusFilter == "inactive" {
		where += ` AND u.is_active = 0`
	}

	// Count total
	countQuery := `SELECT COUNT(DISTINCT u.id) ` + baseQuery + where
	row := r.db.QueryRow(ctx, countQuery, args...)
	var total int
	if err := row.Scan(&total); err != nil {
		return nil, err
	}

	totalPages := (total + params.PerPage - 1) / params.PerPage

	// Get users
	offset := (params.Page - 1) * params.PerPage
	selectQuery := fmt.Sprintf(
		`SELECT DISTINCT u.id, u.email, u.name, u.avatar, u.is_active, u.created_at, u.updated_at, u.last_seen_at %s%s ORDER BY u.created_at DESC LIMIT ? OFFSET ?`,
		baseQuery, where,
	)
	args = append(args, params.PerPage, offset)

	rows, err := r.db.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []UserWithProfile
	for rows.Next() {
		var u UserWithProfile
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Avatar, &u.IsActive,
			&u.CreatedAt, &u.UpdatedAt, &u.LastSeen); err != nil {
			return nil, err
		}
		items = append(items, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if items == nil {
		items = []UserWithProfile{}
	}

	return &PaginatedResponse{
		Items:      items,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

// GetUserByID returns a user with their profile by ID.
func (r *Repository) GetUserByID(ctx context.Context, id string) (*UserWithProfile, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, email, name, avatar, is_active, created_at, updated_at, last_seen_at
		 FROM users WHERE id = ?`, id)

	u := &UserWithProfile{}
	if err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Avatar, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt, &u.LastSeen); err != nil {
		return nil, err
	}

	// Get profile
	profile, _ := r.GetUserProfile(ctx, id)
	u.Profile = profile

	return u, nil
}

// UpdateUser updates a user's basic info.
func (r *Repository) UpdateUser(ctx context.Context, id string, req UpdateUserRequest) error {
	if req.IsActive != nil {
		return r.db.Exec(ctx,
			`UPDATE users SET name = ?, email = ?, is_active = ?, updated_at = ? WHERE id = ?`,
			req.Name, req.Email, *req.IsActive, time.Now().UTC(), id)
	}
	return r.db.Exec(ctx,
		`UPDATE users SET name = ?, email = ?, updated_at = ? WHERE id = ?`,
		req.Name, req.Email, time.Now().UTC(), id)
}

// SoftDeleteUser deactivates a user.
func (r *Repository) SoftDeleteUser(ctx context.Context, id string) error {
	return r.db.Exec(ctx,
		`UPDATE users SET is_active = 0, updated_at = ? WHERE id = ?`, time.Now().UTC(), id)
}

// GetUserProfile returns the profile for a user.
func (r *Repository) GetUserProfile(ctx context.Context, userID string) (*UserProfile, error) {
	row := r.db.QueryRow(ctx,
		`SELECT user_id, bio, location, website, phone, timezone, language, notification_prefs
		 FROM user_profiles WHERE user_id = ?`, userID)

	p := &UserProfile{}
	if err := row.Scan(&p.UserID, &p.Bio, &p.Location, &p.Website, &p.Phone,
		&p.Timezone, &p.Language, &p.NotificationPrefs); err != nil {
		return nil, err
	}
	return p, nil
}

// UpdateUserProfile upserts a user profile.
func (r *Repository) UpdateUserProfile(ctx context.Context, userID string, req UpdateProfileRequest) error {
	return r.db.Exec(ctx,
		`INSERT INTO user_profiles (user_id, bio, location, website, phone, timezone, language, notification_prefs)
		 VALUES (?, ?, ?, ?, ?, ?, ?, '{}')
		 ON CONFLICT(user_id) DO UPDATE SET
		   bio = excluded.bio,
		   location = excluded.location,
		   website = excluded.website,
		   phone = excluded.phone,
		   timezone = excluded.timezone,
		   language = excluded.language`,
		userID, req.Bio, req.Location, req.Website, req.Phone, req.Timezone, req.Language)
}

// UpdateAvatar updates a user's avatar path.
func (r *Repository) UpdateAvatar(ctx context.Context, userID string, avatarPath string) error {
	return r.db.Exec(ctx,
		`UPDATE users SET avatar = ?, updated_at = ? WHERE id = ?`, avatarPath, time.Now().UTC(), userID)
}

// UpdatePassword updates a user's password hash.
func (r *Repository) UpdatePassword(ctx context.Context, userID string, passwordHash string) error {
	return r.db.Exec(ctx,
		`UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?`, passwordHash, time.Now().UTC(), userID)
}

// GetPasswordHash returns the current password hash for a user.
func (r *Repository) GetPasswordHash(ctx context.Context, userID string) (string, error) {
	row := r.db.QueryRow(ctx, `SELECT password_hash FROM users WHERE id = ?`, userID)
	var hash string
	if err := row.Scan(&hash); err != nil {
		return "", err
	}
	return hash, nil
}
