export interface Post {
  id: string
  title: string
  slug: string
  content_md: string
  content_html: string
  excerpt: string
  status: string
  category_id: number | null
  author_id: string
  published_at: string | null
  created_at: string
  updated_at: string
}

export interface Category {
  id: number
  name: string
  slug: string
  description: string
  created_at: string
}

export interface Tag {
  id: number
  name: string
  slug: string
}

export interface PostResponse {
  id: string
  title: string
  slug: string
  content_md: string
  content_html: string
  excerpt: string
  status: string
  category_id: number | null
  author_id: string
  published_at: string | null
  created_at: string
  updated_at: string
  tags: Tag[]
}

export interface PostListResponse {
  posts: PostResponse[]
  total: number
  page: number
  per_page: number
  total_pages: number
}

export interface CreatePostRequest {
  title: string
  content_md: string
  excerpt: string
  category_id?: number | null
  tag_ids?: number[]
}

export interface UpdatePostRequest {
  title: string
  slug: string
  content_md: string
  excerpt: string
  category_id?: number | null
  tag_ids?: number[]
}

export interface CreateCategoryRequest {
  name: string
  description: string
}

export interface CreateTagRequest {
  name: string
}
