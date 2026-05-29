import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: IndexPage,
})

function IndexPage() {
  return (
    <div className="flex min-h-screen items-center justify-center">
      <div className="text-center">
        <h1 className="text-4xl font-bold">Echo SaaS Starter</h1>
        <p className="mt-4 text-lg text-muted-foreground">
          Plugin-based SaaS starter kit
        </p>
      </div>
    </div>
  )
}
