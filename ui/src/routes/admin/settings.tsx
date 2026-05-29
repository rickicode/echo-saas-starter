import { createFileRoute } from '@tanstack/react-router'
import { useState, type FormEvent } from 'react'

export const Route = createFileRoute('/admin/settings')({
  component: SettingsPage,
})

function SettingsPage() {
  const [siteName, setSiteName] = useState('Echo SaaS')
  const [saved, setSaved] = useState(false)

  function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setSaved(true)
    setTimeout(() => setSaved(false), 2000)
  }

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Settings</h1>

      <form onSubmit={handleSubmit} className="max-w-md space-y-4">
        <div className="space-y-2">
          <label htmlFor="siteName" className="text-sm font-medium">
            Site Name
          </label>
          <input
            id="siteName"
            value={siteName}
            onChange={(e) => setSiteName(e.target.value)}
            className="w-full rounded-md border bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring"
          />
        </div>

        <button
          type="submit"
          className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
        >
          Save Settings
        </button>

        {saved && (
          <p className="text-sm text-green-600">Settings saved successfully.</p>
        )}
      </form>
    </div>
  )
}
