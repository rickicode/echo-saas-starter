import { useEffect, useState } from 'react'
import type { ReactNode } from 'react'
import { useAuth } from '../auth/store'
import { get } from '../../lib/api'

interface RequireRoleProps {
  roles: string[]
  children: ReactNode
  fallback?: ReactNode
}

interface UserRole {
  role_name: string
}

export function RequireRole({ roles, children, fallback }: RequireRoleProps) {
  const { isAuthenticated, isLoading, user } = useAuth()
  const [userRoles, setUserRoles] = useState<string[]>([])
  const [rolesLoading, setRolesLoading] = useState(true)

  useEffect(() => {
    if (!isAuthenticated || !user) {
      setRolesLoading(false)
      return
    }

    get<{ roles: UserRole[] }>('/roles/me')
      .then((data) => {
        const roleNames = data.roles.map((r) => r.role_name)
        setUserRoles(roleNames)
      })
      .catch(() => {
        setUserRoles([])
      })
      .finally(() => {
        setRolesLoading(false)
      })
  }, [isAuthenticated, user])

  if (isLoading || rolesLoading) {
    return null
  }

  if (!isAuthenticated) {
    return fallback ?? (
      <div className="flex min-h-screen items-center justify-center">
        <p className="text-lg text-muted-foreground">Please log in to access this page.</p>
      </div>
    )
  }

  const hasRequiredRole = roles.some((role) => userRoles.includes(role))

  if (!hasRequiredRole) {
    return fallback ?? (
      <div className="flex min-h-screen items-center justify-center">
        <p className="text-lg text-muted-foreground">You do not have permission to access this page.</p>
      </div>
    )
  }

  return <>{children}</>
}
