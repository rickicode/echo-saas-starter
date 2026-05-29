export interface User {
  id: string
  email: string
  name: string
  avatar: string
  created_at: string
}

export interface AuthResponse {
  access_token: string
  refresh_token: string
  user: User
}

export interface RegisterRequest {
  name: string
  email: string
  password: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RefreshRequest {
  refresh_token: string
}
