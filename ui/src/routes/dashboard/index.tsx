import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/dashboard/')({
  component: DashboardOverview,
})

function DashboardOverview() {
  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Welcome back</h1>
      <p className="text-muted-foreground mb-6">
        Here is an overview of your account.
      </p>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="rounded-lg border p-4">
          <p className="text-sm text-muted-foreground">Current Plan</p>
          <p className="text-lg font-semibold mt-1">-</p>
        </div>
        <div className="rounded-lg border p-4">
          <p className="text-sm text-muted-foreground">Member Since</p>
          <p className="text-lg font-semibold mt-1">-</p>
        </div>
      </div>
    </div>
  )
}
