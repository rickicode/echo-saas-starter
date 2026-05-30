import { useReducer, useEffect, useCallback, type ReactNode } from 'react'
import { AuthContext, authReducer, initialAuthState } from './store'
import { getMe } from './api'
import type { User } from './types'

interface AuthProviderProps {
  children: ReactNode
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [state, dispatch] = useReducer(authReducer, initialAuthState)

  useEffect(() => {
    const token = localStorage.getItem('access_token')
    if (!token) {
      dispatch({ type: 'LOGOUT' })
      return
    }

    getMe()
      .then((user) => {
        dispatch({ type: 'SET_USER', payload: user })
      })
      .catch(() => {
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        dispatch({ type: 'LOGOUT' })
      })
  }, [])

  const loginUser = useCallback((accessToken: string, refreshToken: string, user: User) => {
    localStorage.setItem('access_token', accessToken)
    localStorage.setItem('refresh_token', refreshToken)
    dispatch({ type: 'SET_USER', payload: user })
  }, [])

  const logoutUser = useCallback(() => {
    const refreshToken = localStorage.getItem('refresh_token')
    if (refreshToken) {
      fetch('/api/v1/auth/logout', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: refreshToken }),
      }).catch(() => {})
    }
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    dispatch({ type: 'LOGOUT' })
  }, [])

  return (
    <AuthContext value={{ ...state, loginUser, logoutUser }}>
      {children}
    </AuthContext>
  )
}
