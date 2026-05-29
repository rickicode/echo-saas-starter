import { createContext, useContext } from 'react'
import type { User } from './types'

export interface AuthState {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
}

export type AuthAction =
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_USER'; payload: User }
  | { type: 'LOGOUT' }

export const initialAuthState: AuthState = {
  user: null,
  isAuthenticated: false,
  isLoading: true,
}

export function authReducer(state: AuthState, action: AuthAction): AuthState {
  switch (action.type) {
    case 'SET_LOADING':
      return { ...state, isLoading: action.payload }
    case 'SET_USER':
      return { user: action.payload, isAuthenticated: true, isLoading: false }
    case 'LOGOUT':
      return { user: null, isAuthenticated: false, isLoading: false }
    default:
      return state
  }
}

export interface AuthContextType extends AuthState {
  loginUser: (accessToken: string, refreshToken: string, user: User) => void
  logoutUser: () => void
}

export const AuthContext = createContext<AuthContextType | null>(null)

export function useAuth(): AuthContextType {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
