import { createRootRoute, Outlet } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/router-devtools'
import { AuthProvider } from '../plugins/auth/AuthProvider'

export const Route = createRootRoute({
  component: RootLayout,
})

function RootLayout() {
  return (
    <AuthProvider>
      <div className="min-h-screen bg-background text-foreground">
        <Outlet />
        <TanStackRouterDevtools position="bottom-right" />
      </div>
    </AuthProvider>
  )
}
