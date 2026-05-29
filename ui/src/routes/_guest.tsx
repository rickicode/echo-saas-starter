import { createFileRoute, Outlet, Link } from '@tanstack/react-router'

export const Route = createFileRoute('/_guest')({
  component: GuestLayout,
})

function GuestLayout() {
  return (
    <div className="min-h-screen flex flex-col">
      <header className="border-b">
        <nav className="container mx-auto flex h-16 items-center justify-between px-4">
          <Link to="/" className="text-xl font-bold">
            Echo SaaS
          </Link>
          <div className="flex items-center gap-4">
            <Link
              to="/login"
              className="text-sm font-medium text-muted-foreground hover:text-foreground"
            >
              Login
            </Link>
            <Link
              to="/register"
              className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
            >
              Register
            </Link>
          </div>
        </nav>
      </header>
      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  )
}
