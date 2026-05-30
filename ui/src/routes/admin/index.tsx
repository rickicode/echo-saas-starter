import { createFileRoute } from '@tanstack/react-router'
import { useEffect, useState } from 'react'
import { get } from '../../lib/api'

export const Route = createFileRoute('/admin/')({
  component: AdminDashboard,
})

interface AdminStats {
  total_users: number
  active_users: number
}

function AdminDashboard() {
  const [stats, setStats] = useState<AdminStats | null>(null)

  useEffect(() => {
    get<AdminStats>('/users/stats')
      .then(setStats)
      .catch(() => setStats(null))
  }, [])

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Admin Dashboard</h1>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="rounded-lg border p-4">
          <p className="text-sm text-muted-foreground">Total Users</p>
          <p className="text-2xl font-bold mt-1">{stats?.total_users ?? '-'}</p>
        </div>
        <div className="rounded-lg border p-4">
          <p className="text-sm text-muted-foreground">Active Users</p>
          <p className="text-2xl font-bold mt-1">{stats?.active_users ?? '-'}</p>
        </div>
        <div className="rounded-lg border p-4">
          <p className="text-sm text-muted-foreground">Inactive Users</p>
          <p className="text-2xl font-bold mt-1">
            {stats ? stats.total_users - stats.active_users : '-'}
          </p>
        </div>
      </div>
    </div>
  )
}
