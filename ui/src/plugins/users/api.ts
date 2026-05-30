import { get, put, del } from '../../lib/api'
import type {
  PaginatedResponse,
  ListUsersParams,
  UserWithProfile,
  UpdateUserRequest,
  UpdateProfileRequest,
  ChangePasswordRequest,
} from './types'

function getAuthHeaders(): Record<string, string> {
  const token = localStorage.getItem('access_token')
  if (token) {
    return { Authorization: `Bearer ${token}` }
  }
  return {}
}

export function getUsers(params?: ListUsersParams): Promise<PaginatedResponse> {
  const query = new URLSearchParams()
  if (params?.search) query.set('search', params.search)
  if (params?.role) query.set('role', params.role)
  if (params?.status) query.set('status', params.status)
  if (params?.page) query.set('page', String(params.page))
  if (params?.per_page) query.set('per_page', String(params.per_page))
  const qs = query.toString()
  return get<PaginatedResponse>(`/users${qs ? '?' + qs : ''}`)
}

export function getUserById(id: string): Promise<UserWithProfile> {
  return get<UserWithProfile>(`/users/${id}`)
}

export function updateUser(id: string, data: UpdateUserRequest): Promise<{ message: string }> {
  return put<{ message: string }>(`/users/${id}`, data)
}

export function deleteUser(id: string): Promise<void> {
  return del<void>(`/users/${id}`)
}

export function getMyProfile(): Promise<UserWithProfile> {
  return get<UserWithProfile>('/users/me')
}

export function updateMyProfile(data: UpdateProfileRequest): Promise<{ message: string }> {
  return put<{ message: string }>('/users/me', data)
}

export async function uploadAvatar(file: File): Promise<{ avatar: string }> {
  const formData = new FormData()
  formData.append('avatar', file)
  const response = await fetch('/api/v1/users/me/avatar', {
    method: 'POST',
    headers: getAuthHeaders(),
    body: formData,
  })
  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`)
  }
  return response.json()
}

export function changePassword(data: ChangePasswordRequest): Promise<{ message: string }> {
  return put<{ message: string }>('/users/me/password', data)
}
