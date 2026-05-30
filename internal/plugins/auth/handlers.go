package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/argon2"
)

const (
	accessTokenExpiry  = 15 * time.Minute
	refreshTokenExpiry = 7 * 24 * time.Hour
)

// Handlers holds dependencies for auth HTTP handlers.
type Handlers struct {
	repo   *Repository
	secret string
}

// NewHandlers creates a new auth handler set.
func NewHandlers(repo *Repository, secret string) *Handlers {
	return &Handlers{repo: repo, secret: secret}
}

// Register handles user registration.
func (h *Handlers) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("AUTH_010", "Invalid request body"))
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Name = strings.TrimSpace(req.Name)

	if req.Email == "" || req.Password == "" || req.Name == "" {
		return c.JSON(http.StatusBadRequest, errorResponse("AUTH_011", "Name, email, and password are required"))
	}

	if len(req.Password) < 8 {
		return c.JSON(http.StatusBadRequest, errorResponse("AUTH_012", "Password must be at least 8 characters"))
	}

	existing, _ := h.repo.GetUserByEmail(c.Request().Context(), req.Email)
	if existing != nil {
		return c.JSON(http.StatusConflict, errorResponse("AUTH_013", "Email already registered"))
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("AUTH_014", "Failed to process registration"))
	}

	now := time.Now().UTC()
	user := &User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: passwordHash,
		Name:         req.Name,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := h.repo.CreateUser(c.Request().Context(), user); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("AUTH_015", "Failed to create user"))
	}

	// Assign default "user" role (role_id=3)
	_ = h.repo.db.Exec(c.Request().Context(),
		`INSERT INTO user_roles (user_id, role_id) VALUES (?, ?) ON CONFLICT DO NOTHING`, user.ID, 3)

	accessToken, err := GenerateAccessToken(user.ID, h.secret, accessTokenExpiry)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("AUTH_016", "Failed to generate token"))
	}

	refreshToken := GenerateRefreshToken()
	rt := &RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: now.Add(refreshTokenExpiry),
		CreatedAt: now,
	}
	if err := h.repo.CreateRefreshToken(c.Request().Context(), rt); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("AUTH_017", "Failed to store refresh token"))
	}

	return c.JSON(http.StatusCreated, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
	})
}

// Login handles user login.
func (h *Handlers) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("AUTH_020", "Invalid request body"))
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, errorResponse("AUTH_021", "Email and password are required"))
	}

	user, err := h.repo.GetUserByEmail(c.Request().Context(), req.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errorResponse("AUTH_022", "Invalid credentials"))
	}

	if !verifyPassword(req.Password, user.PasswordHash) {
		return c.JSON(http.StatusUnauthorized, errorResponse("AUTH_022", "Invalid credentials"))
	}

	if !user.IsActive {
		return c.JSON(http.StatusForbidden, errorResponse("AUTH_023", "Account is deactivated"))
	}

	_ = h.repo.UpdateLastSeen(c.Request().Context(), user.ID)

	now := time.Now().UTC()
	accessToken, err := GenerateAccessToken(user.ID, h.secret, accessTokenExpiry)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("AUTH_024", "Failed to generate token"))
	}

	refreshToken := GenerateRefreshToken()
	rt := &RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: now.Add(refreshTokenExpiry),
		CreatedAt: now,
	}
	if err := h.repo.CreateRefreshToken(c.Request().Context(), rt); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("AUTH_025", "Failed to store refresh token"))
	}

	return c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
	})
}

// Logout handles user logout by deleting the refresh token.
func (h *Handlers) Logout(c echo.Context) error {
	var req RefreshRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("AUTH_030", "Invalid request body"))
	}

	if req.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, errorResponse("AUTH_031", "Refresh token is required"))
	}

	_ = h.repo.DeleteRefreshToken(c.Request().Context(), req.RefreshToken)
	return c.NoContent(http.StatusNoContent)
}

// Refresh handles token refresh with rotation.
func (h *Handlers) Refresh(c echo.Context) error {
	var req RefreshRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("AUTH_040", "Invalid request body"))
	}

	if req.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, errorResponse("AUTH_041", "Refresh token is required"))
	}

	storedToken, err := h.repo.GetRefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errorResponse("AUTH_042", "Invalid refresh token"))
	}

	if time.Now().UTC().After(storedToken.ExpiresAt) {
		_ = h.repo.DeleteRefreshToken(c.Request().Context(), req.RefreshToken)
		return c.JSON(http.StatusUnauthorized, errorResponse("AUTH_043", "Refresh token expired"))
	}

	// Delete old token (rotation)
	_ = h.repo.DeleteRefreshToken(c.Request().Context(), req.RefreshToken)

	user, err := h.repo.GetUserByID(c.Request().Context(), storedToken.UserID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errorResponse("AUTH_044", "User not found"))
	}

	now := time.Now().UTC()
	accessToken, err := GenerateAccessToken(user.ID, h.secret, accessTokenExpiry)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("AUTH_045", "Failed to generate token"))
	}

	newRefreshToken := GenerateRefreshToken()
	rt := &RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: now.Add(refreshTokenExpiry),
		CreatedAt: now,
	}
	if err := h.repo.CreateRefreshToken(c.Request().Context(), rt); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("AUTH_046", "Failed to store refresh token"))
	}

	return c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user.ToResponse(),
	})
}

// Me returns the current authenticated user.
func (h *Handlers) Me(c echo.Context) error {
	user := GetUserFromContext(c)
	if user == nil {
		return c.JSON(http.StatusUnauthorized, errorResponse("AUTH_050", "Not authenticated"))
	}
	return c.JSON(http.StatusOK, user.ToResponse())
}

func errorResponse(code, message string) map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	saltHex := hex.EncodeToString(salt)
	hashHex := hex.EncodeToString(hash)
	return saltHex + ":" + hashHex, nil
}

func verifyPassword(password, storedHash string) bool {
	parts := strings.SplitN(storedHash, ":", 2)
	if len(parts) != 2 {
		return false
	}

	salt, err := hex.DecodeString(parts[0])
	if err != nil {
		return false
	}

	expectedHash, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	return subtle.ConstantTimeCompare(hash, expectedHash) == 1
}
