package users

import (
	"crypto/rand"
	"crypto/subtle"
	"echo-saas-starter/internal/plugins/auth"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/argon2"
)

// Handlers holds dependencies for user HTTP handlers.
type Handlers struct {
	repo *Repository
}

// NewHandlers creates a new user handler set.
func NewHandlers(repo *Repository) *Handlers {
	return &Handlers{repo: repo}
}

// --- Admin routes ---

// ListUsers handles GET /users with query params.
func (h *Handlers) ListUsers(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}

	params := ListUsersParams{
		Search:       c.QueryParam("search"),
		RoleFilter:   c.QueryParam("role"),
		StatusFilter: c.QueryParam("status"),
		Page:         page,
		PerPage:      perPage,
	}

	result, err := h.repo.GetUsers(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("USR_001", "Failed to list users"))
	}

	return c.JSON(http.StatusOK, result)
}

// GetUser handles GET /users/:id.
func (h *Handlers) GetUser(c echo.Context) error {
	id := c.Param("id")
	user, err := h.repo.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, errorResponse("USR_002", "User not found"))
	}
	return c.JSON(http.StatusOK, user)
}

// UpdateUser handles PUT /users/:id.
func (h *Handlers) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("USR_003", "Invalid request body"))
	}

	if err := h.repo.UpdateUser(c.Request().Context(), id, req); err != nil {
		if strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			return c.JSON(http.StatusConflict, errorResponse("USR_007", "Email already in use"))
		}
		return c.JSON(http.StatusInternalServerError, errorResponse("USR_004", "Failed to update user"))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User updated"})
}

// DeleteUser handles DELETE /users/:id (soft delete).
func (h *Handlers) DeleteUser(c echo.Context) error {
	id := c.Param("id")

	// Prevent admin from deactivating themselves
	currentUser := auth.GetUserFromContext(c)
	if currentUser != nil && currentUser.ID == id {
		return c.JSON(http.StatusBadRequest, errorResponse("USR_006", "Cannot deactivate your own account"))
	}

	if err := h.repo.SoftDeleteUser(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("USR_005", "Failed to delete user"))
	}
	return c.NoContent(http.StatusNoContent)
}

// --- User self-service routes ---

// GetMyProfile handles GET /users/me.
func (h *Handlers) GetMyProfile(c echo.Context) error {
	user := auth.GetUserFromContext(c)
	if user == nil {
		return c.JSON(http.StatusUnauthorized, errorResponse("USR_010", "Not authenticated"))
	}

	result, err := h.repo.GetUserByID(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, errorResponse("USR_011", "User not found"))
	}

	return c.JSON(http.StatusOK, result)
}

// UpdateMyProfile handles PUT /users/me.
func (h *Handlers) UpdateMyProfile(c echo.Context) error {
	user := auth.GetUserFromContext(c)
	if user == nil {
		return c.JSON(http.StatusUnauthorized, errorResponse("USR_010", "Not authenticated"))
	}

	var req UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("USR_012", "Invalid request body"))
	}

	if err := h.repo.UpdateUserProfile(c.Request().Context(), user.ID, req); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("USR_013", "Failed to update profile"))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Profile updated"})
}

// UploadAvatar handles POST /users/me/avatar.
func (h *Handlers) UploadAvatar(c echo.Context) error {
	user := auth.GetUserFromContext(c)
	if user == nil {
		return c.JSON(http.StatusUnauthorized, errorResponse("USR_010", "Not authenticated"))
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("USR_014", "No file uploaded"))
	}

	// Validate file size (max 5MB)
	const maxFileSize = 5 * 1024 * 1024
	if file.Size > maxFileSize {
		return c.JSON(http.StatusBadRequest, errorResponse("USR_030", "File size exceeds 5MB limit"))
	}

	// Validate content type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return c.JSON(http.StatusBadRequest, errorResponse("USR_031", "Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed"))
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("USR_015", "Failed to read file"))
	}
	defer src.Close()

	// Ensure uploads directory exists
	uploadDir := "./uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("USR_016", "Failed to create upload directory"))
	}

	// Map content type to extension for safety (ignore user-provided extension)
	extMap := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/gif":  ".gif",
		"image/webp": ".webp",
	}
	ext := extMap[contentType]
	filename := fmt.Sprintf("%s%s", user.ID, ext)
	dstPath := filepath.Join(uploadDir, filename)

	// Remove old avatar file if it exists
	if user.Avatar != "" {
		oldPath := "." + user.Avatar
		os.Remove(oldPath) // ignore error - file might not exist
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("USR_017", "Failed to save file"))
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("USR_018", "Failed to write file"))
	}

	avatarURL := "/uploads/avatars/" + filename
	if err := h.repo.UpdateAvatar(c.Request().Context(), user.ID, avatarURL); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("USR_019", "Failed to update avatar"))
	}

	return c.JSON(http.StatusOK, map[string]string{"avatar": avatarURL})
}

// ChangePassword handles PUT /users/me/password.
func (h *Handlers) ChangePassword(c echo.Context) error {
	user := auth.GetUserFromContext(c)
	if user == nil {
		return c.JSON(http.StatusUnauthorized, errorResponse("USR_010", "Not authenticated"))
	}

	var req ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("USR_020", "Invalid request body"))
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		return c.JSON(http.StatusBadRequest, errorResponse("USR_021", "Current and new password are required"))
	}

	if len(req.NewPassword) < 8 {
		return c.JSON(http.StatusBadRequest, errorResponse("USR_022", "New password must be at least 8 characters"))
	}

	// Verify current password
	currentHash, err := h.repo.GetPasswordHash(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("USR_023", "Failed to verify password"))
	}

	if !verifyPassword(req.CurrentPassword, currentHash) {
		return c.JSON(http.StatusBadRequest, errorResponse("USR_024", "Current password is incorrect"))
	}

	// Hash new password
	newHash, err := hashPassword(req.NewPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("USR_025", "Failed to hash password"))
	}

	if err := h.repo.UpdatePassword(c.Request().Context(), user.ID, newHash); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("USR_026", "Failed to update password"))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password changed"})
}

// GetStats handles GET /users/stats - returns basic user counts.
func (h *Handlers) GetStats(c echo.Context) error {
	stats, err := h.repo.GetStats(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("USR_040", "Failed to get stats"))
	}
	return c.JSON(http.StatusOK, stats)
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
