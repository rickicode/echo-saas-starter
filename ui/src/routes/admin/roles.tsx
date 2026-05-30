import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { get } from '../../lib/api'

export const Route = createFileRoute('/admin/roles')({
  component: RolesPage,
})

interface Role {
  id: number
  name: string
  description: string
  created_at: string
}

function RolesPage() {
  const [roles, setRoles] = useState<Role[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    get<Role[]>('/roles')
      .then(setRoles)
      .catch(() => setRoles([]))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Roles</h1>

      {loading && <p className="text-muted-foreground">Loading...</p>}

      {!loading && (
        <div className="rounded-lg border">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b bg-muted/50">
                <th className="px-4 py-3 text-left font-medium">ID</th>
                <th className="px-4 py-3 text-left font-medium">Name</th>
                <th className="px-4 py-3 text-left font-medium">Description</th>
                <th className="px-4 py-3 text-left font-medium">Created</th>
              </tr>
            </thead>
            <tbody>
              {roles.map((role) => (
                <tr key={role.id} className="border-b">
                  <td className="px-4 py-3">{role.id}</td>
                  <td className="px-4 py-3 font-medium">{role.name}</td>
                  <td className="px-4 py-3 text-muted-foreground">{role.description}</td>
                  <td className="px-4 py-3 text-muted-foreground">
                    {new Date(role.created_at).toLocaleDateString()}
                  </td>
                </tr>
              ))}
              {roles.length === 0 && (
                <tr>
                  <td colSpan={4} className="px-4 py-8 text-center text-muted-foreground">
                    No roles found
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
