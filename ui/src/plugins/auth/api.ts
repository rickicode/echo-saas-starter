import { post, get } from '../../lib/api'
import type { AuthResponse, LoginRequest, RegisterRequest, User } from './types'

export function register(data: RegisterRequest): Promise<AuthResponse> {
  return post<AuthResponse>('/auth/register', data)
}

export function login(data: LoginRequest): Promise<AuthResponse> {
  return post<AuthResponse>('/auth/login', data)
}

export function logout(refreshToken: string): Promise<void> {
  return post<void>('/auth/logout', { refresh_token: refreshToken })
}

export function refresh(refreshToken: string): Promise<AuthResponse> {
  return post<AuthResponse>('/auth/refresh', { refresh_token: refreshToken })
}

export function getMe(): Promise<User> {
  return get<User>('/auth/me')
}
