import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { getPlans, createPlan, deletePlan } from '../../plugins/plans/api'
import type { Plan } from '../../plugins/plans/types'

export const Route = createFileRoute('/admin/plans')({
  component: AdminPlansPage,
})

function AdminPlansPage() {
  const [plans, setPlans] = useState<Plan[]>([])
  const [showForm, setShowForm] = useState(false)
  const [name, setName] = useState('')
  const [slug, setSlug] = useState('')
  const [priceMonthly, setPriceMonthly] = useState('')
  const [priceYearly, setPriceYearly] = useState('')
  const [description, setDescription] = useState('')

  useEffect(() => {
    loadPlans()
  }, [])

  async function loadPlans() {
    try {
      const result = await getPlans()
      setPlans(result)
    } catch {
      // handle error silently
    }
  }

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    try {
      await createPlan({
        name,
        slug,
        description,
        price_monthly: parseFloat(priceMonthly) || 0,
        price_yearly: parseFloat(priceYearly) || 0,
        features: '[]',
        is_active: true,
      })
      setShowForm(false)
      setName('')
      setSlug('')
      setPriceMonthly('')
      setPriceYearly('')
      setDescription('')
      loadPlans()
    } catch {
      // handle error silently
    }
  }

  async function handleDelete(id: number) {
    if (!confirm('Delete this plan?')) return
    try {
      await deletePlan(id)
      loadPlans()
    } catch {
      // handle error silently
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-2xl font-bold">Plans</h1>
        <button
          onClick={() => setShowForm(!showForm)}
          className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
        >
          {showForm ? 'Cancel' : 'New Plan'}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleCreate} className="rounded-lg border p-4 mb-4 space-y-3 max-w-md">
          <input
            type="text"
            placeholder="Plan name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full rounded-md border px-3 py-2 text-sm"
            required
          />
          <input
            type="text"
            placeholder="Slug"
            value={slug}
            onChange={(e) => setSlug(e.target.value)}
            className="w-full rounded-md border px-3 py-2 text-sm"
            required
          />
          <input
            type="text"
            placeholder="Description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            className="w-full rounded-md border px-3 py-2 text-sm"
          />
          <div className="flex gap-2">
            <input
              type="number"
              placeholder="Monthly price"
              value={priceMonthly}
              onChange={(e) => setPriceMonthly(e.target.value)}
              className="flex-1 rounded-md border px-3 py-2 text-sm"
            />
            <input
              type="number"
              placeholder="Yearly price"
              value={priceYearly}
              onChange={(e) => setPriceYearly(e.target.value)}
              className="flex-1 rounded-md border px-3 py-2 text-sm"
            />
          </div>
          <button
            type="submit"
            className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground"
          >
            Create
          </button>
        </form>
      )}

      <div className="rounded-lg border overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-muted">
            <tr>
              <th className="px-4 py-2 text-left font-medium">Name</th>
              <th className="px-4 py-2 text-left font-medium">Slug</th>
              <th className="px-4 py-2 text-left font-medium">Monthly</th>
              <th className="px-4 py-2 text-left font-medium">Yearly</th>
              <th className="px-4 py-2 text-left font-medium">Status</th>
              <th className="px-4 py-2 text-left font-medium">Actions</th>
            </tr>
          </thead>
          <tbody>
            {plans.map((plan) => (
              <tr key={plan.id} className="border-t">
                <td className="px-4 py-2">{plan.name}</td>
                <td className="px-4 py-2">{plan.slug}</td>
                <td className="px-4 py-2">${plan.price_monthly}</td>
                <td className="px-4 py-2">${plan.price_yearly}</td>
                <td className="px-4 py-2">
                  <span
                    className={`text-xs px-2 py-0.5 rounded ${plan.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}
                  >
                    {plan.is_active ? 'Active' : 'Inactive'}
                  </span>
                </td>
                <td className="px-4 py-2">
                  <button
                    onClick={() => handleDelete(plan.id)}
                    className="text-xs text-red-600 hover:underline"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
