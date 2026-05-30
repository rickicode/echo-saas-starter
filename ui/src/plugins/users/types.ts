export interface UserProfile {
  user_id: string
  bio: string
  location: string
  website: string
  phone: string
  timezone: string
  language: string
  notification_prefs: string
}

export interface UserWithProfile {
  id: string
  email: string
  name: string
  avatar: string
  is_active: boolean
  created_at: string
  updated_at: string
  last_seen_at?: string
  profile?: UserProfile
  roles?: string[]
}

export interface ListUsersParams {
  search?: string
  role?: string
  status?: string
  page?: number
  per_page?: number
}

export interface PaginatedResponse {
  items: UserWithProfile[]
  total: number
  page: number
  per_page: number
  total_pages: number
}

export interface UpdateProfileRequest {
  bio: string
  location: string
  website: string
  phone: string
  timezone: string
  language: string
}

export interface ChangePasswordRequest {
  current_password: string
  new_password: string
}

export interface UpdateUserRequest {
  name: string
  email: string
  is_active?: boolean
}
