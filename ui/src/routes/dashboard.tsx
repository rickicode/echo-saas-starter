import { createFileRoute, Outlet, Link } from '@tanstack/react-router'
import { useAuth } from '../plugins/auth/store'

export const Route = createFileRoute('/dashboard')({
  component: DashboardLayout,
})

function DashboardLayout() {
  const { isAuthenticated, isLoading } = useAuth()

  if (isLoading) {
    return null
  }

  if (!isAuthenticated) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <p className="text-lg text-muted-foreground mb-4">Please log in to access your dashboard.</p>
          <Link to="/" className="text-sm font-medium underline">
            Go to Home
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen flex">
      <aside className="w-64 border-r bg-muted/30 p-4">
        <div className="mb-8">
          <Link to="/" className="text-xl font-bold">
            Echo SaaS
          </Link>
          <p className="text-xs text-muted-foreground mt-1">Dashboard</p>
        </div>
        <nav className="space-y-1">
          <NavLink to="/dashboard">Overview</NavLink>
          <NavLink to="/dashboard/profile">Profile</NavLink>
          <NavLink to="/dashboard/subscription">Subscription</NavLink>
          <NavLink to="/dashboard/settings">Settings</NavLink>
        </nav>
      </aside>
      <div className="flex-1 flex flex-col">
        <header className="border-b h-14 flex items-center px-6">
          <h2 className="text-sm font-medium text-muted-foreground">My Account</h2>
        </header>
        <main className="flex-1 p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}

function NavLink({ to, children }: { to: string; children: React.ReactNode }) {
  return (
    <Link
      to={to}
      className="block rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-muted hover:text-foreground"
    >
      {children}
    </Link>
  )
}
