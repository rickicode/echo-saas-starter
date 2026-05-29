import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect, type FormEvent } from 'react'
import { listPosts, createPost, publishPost, unpublishPost, deletePost, listCategories } from '../../plugins/docs/api'
import type { PostListResponse, Category } from '../../plugins/docs/types'

export const Route = createFileRoute('/admin/content')({
  component: ContentPage,
})

function ContentPage() {
  const [data, setData] = useState<PostListResponse | null>(null)
  const [categories, setCategories] = useState<Category[]>([])
  const [showForm, setShowForm] = useState(false)
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [excerpt, setExcerpt] = useState('')
  const [categoryId, setCategoryId] = useState<string>('')
  const [statusFilter, setStatusFilter] = useState('')
  const [loading, setLoading] = useState(true)

  function loadPosts() {
    setLoading(true)
    listPosts({ status: statusFilter || undefined })
      .then(setData)
      .catch(() => setData(null))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    loadPosts()
    listCategories().then(setCategories).catch(() => {})
  }, [statusFilter])

  async function handleCreate(e: FormEvent) {
    e.preventDefault()
    try {
      await createPost({
        title,
        content_md: content,
        excerpt,
        category_id: categoryId ? Number(categoryId) : null,
      })
      setTitle('')
      setContent('')
      setExcerpt('')
      setCategoryId('')
      setShowForm(false)
      loadPosts()
    } catch {
      // handle error silently
    }
  }

  async function handlePublish(id: string) {
    await publishPost(id)
    loadPosts()
  }

  async function handleUnpublish(id: string) {
    await unpublishPost(id)
    loadPosts()
  }

  async function handleDelete(id: string) {
    await deletePost(id)
    loadPosts()
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Content Management</h1>
        <button
          onClick={() => setShowForm(!showForm)}
          className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
        >
          {showForm ? 'Cancel' : 'New Post'}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleCreate} className="mb-8 space-y-4 rounded-lg border p-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Title</label>
            <input
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="w-full rounded-md border bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring"
              required
            />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Excerpt</label>
            <input
              value={excerpt}
              onChange={(e) => setExcerpt(e.target.value)}
              className="w-full rounded-md border bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring"
            />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Category</label>
            <select
              value={categoryId}
              onChange={(e) => setCategoryId(e.target.value)}
              className="w-full rounded-md border bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring"
            >
              <option value="">None</option>
              {categories.map((c) => (
                <option key={c.id} value={c.id}>{c.name}</option>
              ))}
            </select>
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Content (Markdown)</label>
            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              rows={8}
              className="w-full rounded-md border bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring font-mono"
            />
          </div>
          <button
            type="submit"
            className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
          >
            Create Post
          </button>
        </form>
      )}

      <div className="mb-4">
        <select
          value={statusFilter}
          onChange={(e) => setStatusFilter(e.target.value)}
          className="rounded-md border bg-background px-3 py-2 text-sm"
        >
          <option value="">All Statuses</option>
          <option value="draft">Draft</option>
          <option value="published">Published</option>
        </select>
      </div>

      {loading && <p className="text-muted-foreground">Loading...</p>}

      {!loading && data && (
        <div className="rounded-lg border">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b bg-muted/50">
                <th className="px-4 py-3 text-left font-medium">Title</th>
                <th className="px-4 py-3 text-left font-medium">Status</th>
                <th className="px-4 py-3 text-left font-medium">Created</th>
                <th className="px-4 py-3 text-right font-medium">Actions</th>
              </tr>
            </thead>
            <tbody>
              {data.posts.map((post) => (
                <tr key={post.id} className="border-b">
                  <td className="px-4 py-3">{post.title}</td>
                  <td className="px-4 py-3">
                    <span className={`inline-block rounded-full px-2 py-0.5 text-xs ${post.status === 'published' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'}`}>
                      {post.status}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-muted-foreground">
                    {new Date(post.created_at).toLocaleDateString()}
                  </td>
                  <td className="px-4 py-3 text-right space-x-2">
                    {post.status === 'draft' ? (
                      <button onClick={() => handlePublish(post.id)} className="text-xs text-primary hover:underline">
                        Publish
                      </button>
                    ) : (
                      <button onClick={() => handleUnpublish(post.id)} className="text-xs text-yellow-600 hover:underline">
                        Unpublish
                      </button>
                    )}
                    <button onClick={() => handleDelete(post.id)} className="text-xs text-destructive hover:underline">
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
              {data.posts.length === 0 && (
                <tr>
                  <td colSpan={4} className="px-4 py-8 text-center text-muted-foreground">
                    No posts found
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
