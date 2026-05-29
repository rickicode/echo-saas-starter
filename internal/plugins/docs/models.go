package docs

import "time"

// Post represents a documentation/blog post.
type Post struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	ContentMD   string     `json:"content_md"`
	ContentHTML string     `json:"content_html"`
	Excerpt     string     `json:"excerpt"`
	Status      string     `json:"status"`
	CategoryID  *int       `json:"category_id"`
	AuthorID    string     `json:"author_id"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Category represents a post category.
type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Tag represents a post tag.
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// PostTag represents the many-to-many relationship between posts and tags.
type PostTag struct {
	PostID string `json:"post_id"`
	TagID  int    `json:"tag_id"`
}

// CreatePostRequest is the request body for creating a post.
type CreatePostRequest struct {
	Title      string `json:"title"`
	ContentMD  string `json:"content_md"`
	Excerpt    string `json:"excerpt"`
	CategoryID *int   `json:"category_id"`
	TagIDs     []int  `json:"tag_ids"`
}

// UpdatePostRequest is the request body for updating a post.
type UpdatePostRequest struct {
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	ContentMD  string `json:"content_md"`
	Excerpt    string `json:"excerpt"`
	CategoryID *int   `json:"category_id"`
	TagIDs     []int  `json:"tag_ids"`
}

// PostResponse is the response for a single post including tags.
type PostResponse struct {
	Post
	Tags []Tag `json:"tags"`
}

// PostListResponse is the paginated response for listing posts.
type PostListResponse struct {
	Posts      []PostResponse `json:"posts"`
	Total     int            `json:"total"`
	Page      int            `json:"page"`
	PerPage   int            `json:"per_page"`
	TotalPages int           `json:"total_pages"`
}

// CreateCategoryRequest is the request body for creating a category.
type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateCategoryRequest is the request body for updating a category.
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// CreateTagRequest is the request body for creating a tag.
type CreateTagRequest struct {
	Name string `json:"name"`
}
