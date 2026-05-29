import type { ReactNode } from 'react'
import { useAuth } from '../auth/store'

interface RequireRoleProps {
  roles: string[]
  children: ReactNode
  fallback?: ReactNode
}

export function RequireRole({ roles, children, fallback }: RequireRoleProps) {
  const { isAuthenticated, isLoading } = useAuth()

  if (isLoading) {
    return null
  }

  if (!isAuthenticated) {
    return fallback ?? (
      <div className="flex min-h-screen items-center justify-center">
        <p className="text-lg text-muted-foreground">Please log in to access this page.</p>
      </div>
    )
  }

  // Role checking would be done with actual user roles from the backend.
  // For now, if user is authenticated they can access.
  // In a full implementation, you would fetch user roles and compare.
  void roles

  return <>{children}</>
}
