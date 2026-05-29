import { get, post, put, del } from '../../lib/api'
import type {
  PostListResponse,
  PostResponse,
  CreatePostRequest,
  UpdatePostRequest,
  Category,
  Tag,
} from './types'

// Admin post operations
export function listPosts(params?: {
  status?: string
  category_id?: number
  search?: string
  page?: number
  per_page?: number
}): Promise<PostListResponse> {
  const query = new URLSearchParams()
  if (params?.status) query.set('status', params.status)
  if (params?.category_id) query.set('category_id', String(params.category_id))
  if (params?.search) query.set('search', params.search)
  if (params?.page) query.set('page', String(params.page))
  if (params?.per_page) query.set('per_page', String(params.per_page))
  const qs = query.toString()
  return get<PostListResponse>(`/docs/posts${qs ? `?${qs}` : ''}`)
}

export function getPost(id: string): Promise<PostResponse> {
  return get<PostResponse>(`/docs/posts/${id}`)
}

export function createPost(data: CreatePostRequest): Promise<PostResponse> {
  return post<PostResponse>('/docs/posts', data)
}

export function updatePost(id: string, data: UpdatePostRequest): Promise<PostResponse> {
  return put<PostResponse>(`/docs/posts/${id}`, data)
}

export function deletePost(id: string): Promise<void> {
  return del<void>(`/docs/posts/${id}`)
}

export function publishPost(id: string): Promise<{ message: string }> {
  return post<{ message: string }>(`/docs/posts/${id}/publish`)
}

export function unpublishPost(id: string): Promise<{ message: string }> {
  return post<{ message: string }>(`/docs/posts/${id}/unpublish`)
}

// Admin category operations
export function listCategories(): Promise<Category[]> {
  return get<Category[]>('/docs/categories')
}

export function createCategory(data: { name: string; description: string }): Promise<{ message: string }> {
  return post<{ message: string }>('/docs/categories', data)
}

export function updateCategory(id: number, data: { name: string; slug: string; description: string }): Promise<{ message: string }> {
  return put<{ message: string }>(`/docs/categories/${id}`, data)
}

export function deleteCategory(id: number): Promise<void> {
  return del<void>(`/docs/categories/${id}`)
}

// Admin tag operations
export function listTags(): Promise<Tag[]> {
  return get<Tag[]>('/docs/tags')
}

export function createTag(data: { name: string }): Promise<{ message: string }> {
  return post<{ message: string }>('/docs/tags', data)
}

export function deleteTag(id: number): Promise<void> {
  return del<void>(`/docs/tags/${id}`)
}

// Public operations
export function listPublishedPosts(page?: number, perPage?: number): Promise<PostListResponse> {
  const query = new URLSearchParams()
  if (page) query.set('page', String(page))
  if (perPage) query.set('per_page', String(perPage))
  const qs = query.toString()
  return get<PostListResponse>(`/docs/public${qs ? `?${qs}` : ''}`)
}

export function getPublishedPost(slug: string): Promise<PostResponse> {
  return get<PostResponse>(`/docs/public/${slug}`)
}

export function searchPosts(q: string): Promise<PostListResponse> {
  return get<PostListResponse>(`/docs/public/search?q=${encodeURIComponent(q)}`)
}
