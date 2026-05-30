package docs

import (
	"echo-saas-starter/internal/plugins/auth"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// Handlers holds dependencies for docs HTTP handlers.
type Handlers struct {
	repo *Repository
}

// NewHandlers creates a new docs handler set.
func NewHandlers(repo *Repository) *Handlers {
	return &Handlers{repo: repo}
}

// --- Admin Post Handlers ---

// ListPostsAdmin handles GET /posts (admin).
func (h *Handlers) ListPostsAdmin(c echo.Context) error {
	status := c.QueryParam("status")
	categoryID, _ := strconv.Atoi(c.QueryParam("category_id"))
	search := c.QueryParam("search")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	posts, total, err := h.repo.ListPosts(c.Request().Context(), status, categoryID, search, page, perPage)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("DOC_001", "Failed to list posts"))
	}

	responses := make([]PostResponse, 0, len(posts))
	for _, p := range posts {
		tags, _ := h.repo.GetPostTags(c.Request().Context(), p.ID)
		responses = append(responses, PostResponse{Post: p, Tags: tags})
	}

	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	return c.JSON(http.StatusOK, PostListResponse{
		Posts:      responses,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}

// GetPostAdmin handles GET /posts/:id (admin).
func (h *Handlers) GetPostAdmin(c echo.Context) error {
	id := c.Param("id")
	post, err := h.repo.GetPostByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, errorResponse("DOC_002", "Post not found"))
	}
	tags, _ := h.repo.GetPostTags(c.Request().Context(), post.ID)
	return c.JSON(http.StatusOK, PostResponse{Post: *post, Tags: tags})
}

// CreatePost handles POST /posts (admin).
func (h *Handlers) CreatePost(c echo.Context) error {
	user := auth.GetUserFromContext(c)
	if user == nil {
		return c.JSON(http.StatusUnauthorized, errorResponse("DOC_003", "Not authenticated"))
	}

	var req CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("DOC_004", "Invalid request body"))
	}

	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, errorResponse("DOC_005", "Title is required"))
	}

	now := time.Now().UTC()
	postID := GenerateID()
	post := &Post{
		ID:          postID,
		Title:       req.Title,
		Slug:        generateSlug(req.Title) + "-" + postID[:8],
		ContentMD:   req.ContentMD,
		ContentHTML: RenderMarkdown(req.ContentMD),
		Excerpt:     req.Excerpt,
		Status:      "draft",
		CategoryID:  req.CategoryID,
		AuthorID:    user.ID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.repo.CreatePost(c.Request().Context(), post); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("DOC_006", "Failed to create post"))
	}

	if len(req.TagIDs) > 0 {
		_ = h.repo.SetPostTags(c.Request().Context(), post.ID, req.TagIDs)
	}

	tags, _ := h.repo.GetPostTags(c.Request().Context(), post.ID)
	return c.JSON(http.StatusCreated, PostResponse{Post: *post, Tags: tags})
}

// UpdatePost handles PUT /posts/:id (admin).
func (h *Handlers) UpdatePost(c echo.Context) error {
	id := c.Param("id")
	post, err := h.repo.GetPostByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, errorResponse("DOC_002", "Post not found"))
	}

	var req UpdatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("DOC_004", "Invalid request body"))
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Slug != "" {
		post.Slug = req.Slug
	}
	if req.ContentMD != "" {
		post.ContentMD = req.ContentMD
		post.ContentHTML = RenderMarkdown(req.ContentMD)
	}
	if req.Excerpt != "" {
		post.Excerpt = req.Excerpt
	}
	if req.CategoryID != nil {
		post.CategoryID = req.CategoryID
	}
	post.UpdatedAt = time.Now().UTC()

	if err := h.repo.UpdatePost(c.Request().Context(), post); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("DOC_007", "Failed to update post"))
	}

	if req.TagIDs != nil {
		_ = h.repo.SetPostTags(c.Request().Context(), post.ID, req.TagIDs)
	}

	tags, _ := h.repo.GetPostTags(c.Request().Context(), post.ID)
	return c.JSON(http.StatusOK, PostResponse{Post: *post, Tags: tags})
}

// DeletePost handles DELETE /posts/:id (admin).
func (h *Handlers) DeletePost(c echo.Context) error {
	id := c.Param("id")
	if err := h.repo.DeletePost(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("DOC_008", "Failed to delete post"))
	}
	return c.NoContent(http.StatusNoContent)
}

// PublishPost handles POST /posts/:id/publish (admin).
func (h *Handlers) PublishPost(c echo.Context) error {
	id := c.Param("id")
	if err := h.repo.PublishPost(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("DOC_009", "Failed to publish post"))
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Post published"})
}

// UnpublishPost handles POST /posts/:id/unpublish (admin).
func (h *Handlers) UnpublishPost(c echo.Context) error {
	id := c.Param("id")
	if err := h.repo.UnpublishPost(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("DOC_010", "Failed to unpublish post"))
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Post unpublished"})
}

// --- Admin Category Handlers ---

// ListCategories handles GET /categories (admin).
func (h *Handlers) ListCategories(c echo.Context) error {
	categories, err := h.repo.GetAllCategories(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("DOC_011", "Failed to list categories"))
	}
	return c.JSON(http.StatusOK, categories)
}

// CreateCategory handles POST /categories (admin).
func (h *Handlers) CreateCategory(c echo.Context) error {
	var req CreateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("DOC_012", "Invalid request body"))
	}
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, errorResponse("DOC_013", "Name is required"))
	}
	slug := generateSlug(req.Name)
	if err := h.repo.CreateCategory(c.Request().Context(), req.Name, slug, req.Description); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("DOC_014", "Failed to create category"))
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "Category created"})
}

// UpdateCategory handles PUT /categories/:id (admin).
func (h *Handlers) UpdateCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("DOC_015", "Invalid category ID"))
	}
	var req UpdateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("DOC_012", "Invalid request body"))
	}
	if req.Slug == "" {
		req.Slug = generateSlug(req.Name)
	}
	if err := h.repo.UpdateCategory(c.Request().Context(), id, req.Name, req.Slug, req.Description); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("DOC_016", "Failed to update category"))
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Category updated"})
}

// DeleteCategory handles DELETE /categories/:id (admin).
func (h *Handlers) DeleteCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("DOC_015", "Invalid category ID"))
	}
	if err := h.repo.DeleteCategory(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("DOC_017", "Failed to delete category"))
	}
	return c.NoContent(http.StatusNoContent)
}

// --- Admin Tag Handlers ---

// ListTags handles GET /tags (admin).
func (h *Handlers) ListTags(c echo.Context) error {
	tags, err := h.repo.GetAllTags(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("DOC_018", "Failed to list tags"))
	}
	return c.JSON(http.StatusOK, tags)
}

// CreateTag handles POST /tags (admin).
func (h *Handlers) CreateTag(c echo.Context) error {
	var req CreateTagRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("DOC_019", "Invalid request body"))
	}
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, errorResponse("DOC_020", "Name is required"))
	}
	slug := generateSlug(req.Name)
	if err := h.repo.CreateTag(c.Request().Context(), req.Name, slug); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("DOC_021", "Failed to create tag"))
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "Tag created"})
}

// DeleteTag handles DELETE /tags/:id (admin).
func (h *Handlers) DeleteTag(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("DOC_022", "Invalid tag ID"))
	}
	if err := h.repo.DeleteTag(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("DOC_023", "Failed to delete tag"))
	}
	return c.NoContent(http.StatusNoContent)
}

// --- Public Handlers ---

// ListPublishedPosts handles GET /docs (public).
func (h *Handlers) ListPublishedPosts(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	posts, total, err := h.repo.ListPosts(c.Request().Context(), "published", 0, "", page, perPage)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("DOC_030", "Failed to list posts"))
	}

	responses := make([]PostResponse, 0, len(posts))
	for _, p := range posts {
		tags, _ := h.repo.GetPostTags(c.Request().Context(), p.ID)
		responses = append(responses, PostResponse{Post: p, Tags: tags})
	}

	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	return c.JSON(http.StatusOK, PostListResponse{
		Posts:      responses,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}

// GetPublishedPost handles GET /docs/:slug (public).
func (h *Handlers) GetPublishedPost(c echo.Context) error {
	slug := c.Param("slug")
	post, err := h.repo.GetPostBySlug(c.Request().Context(), slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, errorResponse("DOC_031", "Post not found"))
	}
	tags, _ := h.repo.GetPostTags(c.Request().Context(), post.ID)
	return c.JSON(http.StatusOK, PostResponse{Post: *post, Tags: tags})
}

// SearchPosts handles GET /docs/search (public).
func (h *Handlers) SearchPosts(c echo.Context) error {
	q := c.QueryParam("q")
	if q == "" {
		return c.JSON(http.StatusOK, PostListResponse{Posts: []PostResponse{}, Total: 0, Page: 1, PerPage: 20, TotalPages: 0})
	}

	posts, total, err := h.repo.ListPosts(c.Request().Context(), "published", 0, q, 1, 20)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("DOC_032", "Search failed"))
	}

	responses := make([]PostResponse, 0, len(posts))
	for _, p := range posts {
		tags, _ := h.repo.GetPostTags(c.Request().Context(), p.ID)
		responses = append(responses, PostResponse{Post: p, Tags: tags})
	}

	totalPages := total / 20
	if total%20 > 0 {
		totalPages++
	}

	return c.JSON(http.StatusOK, PostListResponse{
		Posts:      responses,
		Total:      total,
		Page:       1,
		PerPage:    20,
		TotalPages: totalPages,
	})
}

// --- Helpers ---

var slugRegexp = regexp.MustCompile(`[^a-z0-9-]+`)

func generateSlug(title string) string {
	slug := strings.ToLower(strings.TrimSpace(title))
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = slugRegexp.ReplaceAllString(slug, "")
	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")
	return slug
}

func errorResponse(code, message string) map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
}
