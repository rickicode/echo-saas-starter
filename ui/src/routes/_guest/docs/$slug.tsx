import { createFileRoute, Link } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { getPublishedPost } from '../../../plugins/docs/api'
import type { PostResponse } from '../../../plugins/docs/types'

export const Route = createFileRoute('/_guest/docs/$slug')({
  component: DocDetailPage,
})

function DocDetailPage() {
  const { slug } = Route.useParams()
  const [post, setPost] = useState<PostResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    setLoading(true)
    getPublishedPost(slug)
      .then(setPost)
      .catch(() => setError('Post not found'))
      .finally(() => setLoading(false))
  }, [slug])

  if (loading) {
    return (
      <div className="py-12 px-4">
        <div className="container mx-auto max-w-3xl">
          <p className="text-muted-foreground">Loading...</p>
        </div>
      </div>
    )
  }

  if (error || !post) {
    return (
      <div className="py-12 px-4">
        <div className="container mx-auto max-w-3xl">
          <p className="text-destructive">{error || 'Post not found'}</p>
          <Link to="/docs" className="mt-4 inline-block text-primary hover:underline">
            Back to docs
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="py-12 px-4">
      <div className="container mx-auto max-w-3xl">
        <Link to="/docs" className="text-sm text-primary hover:underline mb-4 inline-block">
          &larr; Back to docs
        </Link>
        <h1 className="text-3xl font-bold mt-2">{post.title}</h1>
        {post.published_at && (
          <p className="mt-2 text-sm text-muted-foreground">
            Published {new Date(post.published_at).toLocaleDateString()}
          </p>
        )}
        {post.tags.length > 0 && (
          <div className="mt-3 flex gap-2">
            {post.tags.map((tag) => (
              <span key={tag.id} className="text-xs rounded-full bg-muted px-2 py-1">
                {tag.name}
              </span>
            ))}
          </div>
        )}
        <div
          className="mt-8 prose prose-sm max-w-none"
          dangerouslySetInnerHTML={{ __html: post.content_html }}
        />
      </div>
    </div>
  )
}
