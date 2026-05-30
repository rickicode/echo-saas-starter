package docs

import (
	"context"
	"echo-saas-starter/internal/core"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Repository handles database operations for docs.
type Repository struct {
	db core.DB
}

// NewRepository creates a new docs repository.
func NewRepository(db core.DB) *Repository {
	return &Repository{db: db}
}

// CreatePost creates a new post.
func (r *Repository) CreatePost(ctx context.Context, post *Post) error {
	return r.db.Exec(ctx,
		`INSERT INTO posts (id, title, slug, content_md, content_html, excerpt, status, category_id, author_id, published_at, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		post.ID, post.Title, post.Slug, post.ContentMD, post.ContentHTML, post.Excerpt,
		post.Status, post.CategoryID, post.AuthorID, post.PublishedAt, post.CreatedAt, post.UpdatedAt,
	)
}

// GetPostByID retrieves a post by ID.
func (r *Repository) GetPostByID(ctx context.Context, id string) (*Post, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, title, slug, content_md, content_html, excerpt, status, category_id, author_id, published_at, created_at, updated_at
		 FROM posts WHERE id = ?`, id)

	p := &Post{}
	if err := row.Scan(&p.ID, &p.Title, &p.Slug, &p.ContentMD, &p.ContentHTML, &p.Excerpt,
		&p.Status, &p.CategoryID, &p.AuthorID, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, err
	}
	return p, nil
}

// GetPostBySlug retrieves a published post by slug.
func (r *Repository) GetPostBySlug(ctx context.Context, slug string) (*Post, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, title, slug, content_md, content_html, excerpt, status, category_id, author_id, published_at, created_at, updated_at
		 FROM posts WHERE slug = ? AND status = 'published'`, slug)

	p := &Post{}
	if err := row.Scan(&p.ID, &p.Title, &p.Slug, &p.ContentMD, &p.ContentHTML, &p.Excerpt,
		&p.Status, &p.CategoryID, &p.AuthorID, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, err
	}
	return p, nil
}

// ListPosts lists posts with filters and pagination.
func (r *Repository) ListPosts(ctx context.Context, status string, categoryID int, search string, page, perPage int) ([]Post, int, error) {
	where := "WHERE 1=1"
	args := []interface{}{}

	if status != "" {
		where += " AND status = ?"
		args = append(args, status)
	}
	if categoryID > 0 {
		where += " AND category_id = ?"
		args = append(args, categoryID)
	}
	if search != "" {
		where += " AND (title LIKE ? OR content_md LIKE ?)"
		s := "%" + search + "%"
		args = append(args, s, s)
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM posts " + where
	row := r.db.QueryRow(ctx, countQuery, args...)
	var total int
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	offset := (page - 1) * perPage
	query := fmt.Sprintf("SELECT id, title, slug, content_md, content_html, excerpt, status, category_id, author_id, published_at, created_at, updated_at FROM posts %s ORDER BY created_at DESC LIMIT ? OFFSET ?", where)
	args = append(args, perPage, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.ContentMD, &p.ContentHTML, &p.Excerpt,
			&p.Status, &p.CategoryID, &p.AuthorID, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		posts = append(posts, p)
	}
	if posts == nil {
		posts = []Post{}
	}
	return posts, total, rows.Err()
}

// UpdatePost updates an existing post.
func (r *Repository) UpdatePost(ctx context.Context, post *Post) error {
	return r.db.Exec(ctx,
		`UPDATE posts SET title = ?, slug = ?, content_md = ?, content_html = ?, excerpt = ?, status = ?, category_id = ?, updated_at = ? WHERE id = ?`,
		post.Title, post.Slug, post.ContentMD, post.ContentHTML, post.Excerpt, post.Status, post.CategoryID, post.UpdatedAt, post.ID,
	)
}

// DeletePost deletes a post by ID.
func (r *Repository) DeletePost(ctx context.Context, id string) error {
	return r.db.Exec(ctx, `DELETE FROM posts WHERE id = ?`, id)
}

// PublishPost sets status to published and sets published_at.
func (r *Repository) PublishPost(ctx context.Context, id string) error {
	now := time.Now().UTC()
	return r.db.Exec(ctx, `UPDATE posts SET status = 'published', published_at = ?, updated_at = ? WHERE id = ?`, now, now, id)
}

// UnpublishPost sets status back to draft.
func (r *Repository) UnpublishPost(ctx context.Context, id string) error {
	now := time.Now().UTC()
	return r.db.Exec(ctx, `UPDATE posts SET status = 'draft', updated_at = ? WHERE id = ?`, now, id)
}

// GetAllCategories returns all categories.
func (r *Repository) GetAllCategories(ctx context.Context) ([]Category, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, slug, description, created_at FROM categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	if result == nil {
		result = []Category{}
	}
	return result, rows.Err()
}

// GetCategoryByID returns a category by ID.
func (r *Repository) GetCategoryByID(ctx context.Context, id int) (*Category, error) {
	row := r.db.QueryRow(ctx, `SELECT id, name, slug, description, created_at FROM categories WHERE id = ?`, id)
	c := &Category{}
	if err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.CreatedAt); err != nil {
		return nil, err
	}
	return c, nil
}

// CreateCategory creates a new category.
func (r *Repository) CreateCategory(ctx context.Context, name, slug, description string) error {
	return r.db.Exec(ctx, `INSERT INTO categories (name, slug, description) VALUES (?, ?, ?)`, name, slug, description)
}

// UpdateCategory updates an existing category.
func (r *Repository) UpdateCategory(ctx context.Context, id int, name, slug, description string) error {
	return r.db.Exec(ctx, `UPDATE categories SET name = ?, slug = ?, description = ? WHERE id = ?`, name, slug, description, id)
}

// DeleteCategory deletes a category.
func (r *Repository) DeleteCategory(ctx context.Context, id int) error {
	return r.db.Exec(ctx, `DELETE FROM categories WHERE id = ?`, id)
}

// GetAllTags returns all tags.
func (r *Repository) GetAllTags(ctx context.Context) ([]Tag, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, slug FROM tags ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	if result == nil {
		result = []Tag{}
	}
	return result, rows.Err()
}

// CreateTag creates a new tag.
func (r *Repository) CreateTag(ctx context.Context, name, slug string) error {
	return r.db.Exec(ctx, `INSERT INTO tags (name, slug) VALUES (?, ?)`, name, slug)
}

// DeleteTag deletes a tag.
func (r *Repository) DeleteTag(ctx context.Context, id int) error {
	return r.db.Exec(ctx, `DELETE FROM tags WHERE id = ?`, id)
}

// SetPostTags replaces all tags for a post.
func (r *Repository) SetPostTags(ctx context.Context, postID string, tagIDs []int) error {
	_ = r.db.Exec(ctx, `DELETE FROM post_tags WHERE post_id = ?`, postID)
	for _, tagID := range tagIDs {
		if err := r.db.Exec(ctx, `INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)`, postID, tagID); err != nil {
			return err
		}
	}
	return nil
}

// GetPostTags returns all tags for a post.
func (r *Repository) GetPostTags(ctx context.Context, postID string) ([]Tag, error) {
	rows, err := r.db.Query(ctx,
		`SELECT t.id, t.name, t.slug FROM tags t
		 INNER JOIN post_tags pt ON pt.tag_id = t.id
		 WHERE pt.post_id = ?`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	if result == nil {
		result = []Tag{}
	}
	return result, rows.Err()
}

// GenerateID generates a new UUID.
func GenerateID() string {
	return uuid.New().String()
}
