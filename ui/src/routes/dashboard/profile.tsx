import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { getMyProfile, updateMyProfile, uploadAvatar } from '../../plugins/users/api'
import type { UpdateProfileRequest } from '../../plugins/users/types'

export const Route = createFileRoute('/dashboard/profile')({
  component: ProfilePage,
})

function ProfilePage() {
  const [form, setForm] = useState<UpdateProfileRequest>({
    bio: '',
    location: '',
    website: '',
    phone: '',
    timezone: 'UTC',
    language: 'en',
  })
  const [message, setMessage] = useState('')
  const [avatar, setAvatar] = useState('')

  useEffect(() => {
    getMyProfile()
      .then((data) => {
        setAvatar(data.avatar || '')
        if (data.profile) {
          setForm({
            bio: data.profile.bio || '',
            location: data.profile.location || '',
            website: data.profile.website || '',
            phone: data.profile.phone || '',
            timezone: data.profile.timezone || 'UTC',
            language: data.profile.language || 'en',
          })
        }
      })
      .catch(() => {})
  }, [])

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    try {
      await updateMyProfile(form)
      setMessage('Profile updated')
      setTimeout(() => setMessage(''), 3000)
    } catch {
      setMessage('Failed to update profile')
    }
  }

  async function handleAvatar(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (!file) return
    try {
      const result = await uploadAvatar(file)
      setAvatar(result.avatar)
    } catch {
      // handle error silently
    }
  }

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Profile</h1>
      <div className="max-w-lg space-y-6">
        <div>
          <p className="text-sm font-medium mb-2">Avatar</p>
          {avatar && (
            <img src={avatar} alt="Avatar" className="w-16 h-16 rounded-full mb-2 object-cover" />
          )}
          <input type="file" accept="image/*" onChange={handleAvatar} className="text-sm" />
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="text-sm font-medium">Bio</label>
            <textarea
              value={form.bio}
              onChange={(e) => setForm({ ...form, bio: e.target.value })}
              className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
              rows={3}
            />
          </div>
          <div>
            <label className="text-sm font-medium">Location</label>
            <input
              type="text"
              value={form.location}
              onChange={(e) => setForm({ ...form, location: e.target.value })}
              className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            />
          </div>
          <div>
            <label className="text-sm font-medium">Website</label>
            <input
              type="text"
              value={form.website}
              onChange={(e) => setForm({ ...form, website: e.target.value })}
              className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            />
          </div>
          <div>
            <label className="text-sm font-medium">Phone</label>
            <input
              type="text"
              value={form.phone}
              onChange={(e) => setForm({ ...form, phone: e.target.value })}
              className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            />
          </div>
          <div>
            <label className="text-sm font-medium">Timezone</label>
            <input
              type="text"
              value={form.timezone}
              onChange={(e) => setForm({ ...form, timezone: e.target.value })}
              className="mt-1 w-full rounded-md border px-3 py-2 text-sm"
            />
          </div>

          {message && <p className="text-sm text-green-600">{message}</p>}

          <button
            type="submit"
            className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
          >
            Save Profile
          </button>
        </form>
      </div>
    </div>
  )
}
