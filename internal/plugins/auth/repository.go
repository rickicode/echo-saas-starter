package auth

import (
	"context"
	"echo-saas-starter/internal/core"
	"time"
)

// Repository handles database operations for auth.
type Repository struct {
	db core.DB
}

// NewRepository creates a new auth repository.
func NewRepository(db core.DB) *Repository {
	return &Repository{db: db}
}

// CreateUser creates a new user in the database.
func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	return r.db.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, name, avatar, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Email, user.PasswordHash, user.Name, user.Avatar, user.IsActive, user.CreatedAt, user.UpdatedAt,
	)
}

// GetUserByEmail retrieves a user by email.
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, email, password_hash, name, avatar, is_active, created_at, updated_at, last_seen_at
		 FROM users WHERE email = ?`, email)

	user := &User{}
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Avatar,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.LastSeenAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user by ID.
func (r *Repository) GetUserByID(ctx context.Context, id string) (*User, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, email, password_hash, name, avatar, is_active, created_at, updated_at, last_seen_at
		 FROM users WHERE id = ?`, id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Avatar,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.LastSeenAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateLastSeen updates the last_seen_at timestamp for a user.
func (r *Repository) UpdateLastSeen(ctx context.Context, userID string) error {
	return r.db.Exec(ctx,
		`UPDATE users SET last_seen_at = ? WHERE id = ?`, time.Now().UTC(), userID)
}

// CreateRefreshToken stores a new refresh token.
func (r *Repository) CreateRefreshToken(ctx context.Context, token *RefreshToken) error {
	return r.db.Exec(ctx,
		`INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		token.ID, token.UserID, token.Token, token.ExpiresAt, token.CreatedAt,
	)
}

// GetRefreshToken retrieves a refresh token by token string.
func (r *Repository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, user_id, token, expires_at, created_at
		 FROM refresh_tokens WHERE token = ?`, token)

	rt := &RefreshToken{}
	err := row.Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		return nil, err
	}
	return rt, nil
}

// DeleteRefreshToken removes a refresh token.
func (r *Repository) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE token = ?`, token)
}

// DeleteUserRefreshTokens removes all refresh tokens for a user.
func (r *Repository) DeleteUserRefreshTokens(ctx context.Context, userID string) error {
	return r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE user_id = ?`, userID)
}
