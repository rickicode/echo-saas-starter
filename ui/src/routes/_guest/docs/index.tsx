import { createFileRoute, Link } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { listPublishedPosts, searchPosts } from '../../../plugins/docs/api'
import type { PostListResponse } from '../../../plugins/docs/types'

export const Route = createFileRoute('/_guest/docs/')({
  component: DocsIndexPage,
})

function DocsIndexPage() {
  const [data, setData] = useState<PostListResponse | null>(null)
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    setLoading(true)
    const timer = setTimeout(() => {
      const fn = search ? searchPosts(search) : listPublishedPosts(1, 20)
      fn.then(setData).catch(() => setData(null)).finally(() => setLoading(false))
    }, 300)
    return () => clearTimeout(timer)
  }, [search])

  return (
    <div className="py-12 px-4">
      <div className="container mx-auto max-w-4xl">
        <h1 className="text-3xl font-bold mb-6">Documentation</h1>

        <div className="mb-8">
          <input
            type="text"
            placeholder="Search docs..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full max-w-md rounded-md border bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring"
          />
        </div>

        {loading && <p className="text-muted-foreground">Loading...</p>}

        {!loading && data && data.posts.length === 0 && (
          <p className="text-muted-foreground">No posts found.</p>
        )}

        {!loading && data && data.posts.length > 0 && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {data.posts.map((post) => (
              <Link
                key={post.id}
                to="/docs/$slug"
                params={{ slug: post.slug }}
                className="block rounded-lg border p-4 hover:bg-muted/50 transition-colors"
              >
                <h2 className="font-semibold">{post.title}</h2>
                {post.excerpt && (
                  <p className="mt-1 text-sm text-muted-foreground line-clamp-2">
                    {post.excerpt}
                  </p>
                )}
                {post.published_at && (
                  <p className="mt-2 text-xs text-muted-foreground">
                    {new Date(post.published_at).toLocaleDateString()}
                  </p>
                )}
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
