import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { getUserById } from '../../plugins/users/api'
import type { UserWithProfile } from '../../plugins/users/types'

export const Route = createFileRoute('/admin/users_/$id')({
  component: AdminUserDetail,
})

function AdminUserDetail() {
  const { id } = Route.useParams()
  const [user, setUser] = useState<UserWithProfile | null>(null)

  useEffect(() => {
    getUserById(id).then(setUser).catch(() => {})
  }, [id])

  if (!user) {
    return <div className="text-muted-foreground">Loading...</div>
  }

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">User Detail</h1>
      <div className="rounded-lg border p-6 max-w-lg">
        <div className="space-y-4">
          <div>
            <p className="text-sm text-muted-foreground">Name</p>
            <p className="font-medium">{user.name}</p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground">Email</p>
            <p className="font-medium">{user.email}</p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground">Status</p>
            <span
              className={`text-xs px-2 py-0.5 rounded ${user.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}
            >
              {user.is_active ? 'Active' : 'Inactive'}
            </span>
          </div>
          <div>
            <p className="text-sm text-muted-foreground">Joined</p>
            <p className="font-medium">{new Date(user.created_at).toLocaleDateString()}</p>
          </div>
          {user.profile && (
            <>
              <div>
                <p className="text-sm text-muted-foreground">Bio</p>
                <p>{user.profile.bio || '-'}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Location</p>
                <p>{user.profile.location || '-'}</p>
              </div>
            </>
          )}
          {user.roles && user.roles.length > 0 && (
            <div>
              <p className="text-sm text-muted-foreground">Roles</p>
              <div className="flex gap-1 mt-1">
                {user.roles.map((role) => (
                  <span key={role} className="text-xs px-2 py-0.5 rounded bg-blue-100 text-blue-800">
                    {role}
                  </span>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
