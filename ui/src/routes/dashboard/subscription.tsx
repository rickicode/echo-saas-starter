import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { getCurrentPlan, getAvailablePlans } from '../../plugins/plans/api'
import type { Plan, UserPlanWithDetails } from '../../plugins/plans/types'

export const Route = createFileRoute('/dashboard/subscription')({
  component: SubscriptionPage,
})

function SubscriptionPage() {
  const [current, setCurrent] = useState<UserPlanWithDetails | null>(null)
  const [available, setAvailable] = useState<Plan[]>([])

  useEffect(() => {
    getCurrentPlan().then(setCurrent).catch(() => {})
    getAvailablePlans().then(setAvailable).catch(() => {})
  }, [])

  function parseFeatures(features: string): string[] {
    try {
      return JSON.parse(features)
    } catch {
      return []
    }
  }

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Subscription</h1>

      {current && current.plan && (
        <div className="rounded-lg border p-4 mb-8 max-w-md">
          <h2 className="text-lg font-semibold">
            Current Plan: {current.plan.name}
          </h2>
          <p className="text-sm text-muted-foreground mt-1">
            {current.plan.description}
          </p>
          <p className="text-sm mt-2">
            Status:{' '}
            <span className="font-medium">{current.user_plan.status}</span>
          </p>
          <div className="mt-3">
            <p className="text-xs font-medium text-muted-foreground mb-1">
              Features:
            </p>
            <ul className="list-disc list-inside text-sm space-y-0.5">
              {parseFeatures(current.plan.features).map((f) => (
                <li key={f}>{f}</li>
              ))}
            </ul>
          </div>
        </div>
      )}

      {!current && (
        <p className="text-muted-foreground mb-8">
          You do not have an active plan.
        </p>
      )}

      <h2 className="text-lg font-semibold mb-4">Available Plans</h2>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {available.map((plan) => (
          <div key={plan.id} className="rounded-lg border p-4">
            <h3 className="font-semibold">{plan.name}</h3>
            <p className="text-sm text-muted-foreground mt-1">
              {plan.description}
            </p>
            <p className="text-xl font-bold mt-2">
              ${plan.price_monthly}/mo
            </p>
            <p className="text-xs text-muted-foreground">
              ${plan.price_yearly}/yr
            </p>
            <div className="mt-3">
              <ul className="list-disc list-inside text-sm space-y-0.5">
                {parseFeatures(plan.features).map((f) => (
                  <li key={f}>{f}</li>
                ))}
              </ul>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
