import { createFileRoute, Link } from '@tanstack/react-router'
import { Shield, Layers, Database, Users, CreditCard, FileText } from 'lucide-react'

export const Route = createFileRoute('/_guest/')({
  component: LandingPage,
})

const features = [
  { icon: Shield, title: 'Authentication', description: 'Secure PASETO v4 token-based auth with refresh token rotation.' },
  { icon: Layers, title: 'Plugin Architecture', description: 'Modular plugin system for easy feature addition and removal.' },
  { icon: Database, title: 'Multi-DB Support', description: 'SQLite for development, PostgreSQL for production.' },
  { icon: Users, title: 'Role-Based Access', description: 'Granular RBAC with super admin, admin, and user roles.' },
  { icon: CreditCard, title: 'Subscription Plans', description: 'Built-in plan management with user assignment.' },
  { icon: FileText, title: 'Content Management', description: 'Markdown-based docs CMS with categories and tags.' },
]

function LandingPage() {
  return (
    <div>
      <section className="py-20 px-4">
        <div className="container mx-auto text-center max-w-3xl">
          <h1 className="text-5xl font-bold tracking-tight">
            Build SaaS Products Faster
          </h1>
          <p className="mt-6 text-xl text-muted-foreground">
            A production-ready starter kit with Go backend, React frontend, and a plugin-based architecture that scales with your needs.
          </p>
          <div className="mt-8 flex justify-center gap-4">
            <Link
              to="/register"
              className="rounded-md bg-primary px-6 py-3 text-sm font-medium text-primary-foreground hover:bg-primary/90"
            >
              Get Started
            </Link>
            <Link
              to="/docs"
              className="rounded-md border px-6 py-3 text-sm font-medium hover:bg-muted"
            >
              Documentation
            </Link>
          </div>
        </div>
      </section>

      <section className="py-16 px-4 border-t">
        <div className="container mx-auto max-w-5xl">
          <h2 className="text-3xl font-bold text-center mb-12">Everything You Need</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {features.map((feature) => (
              <div key={feature.title} className="rounded-lg border p-6">
                <feature.icon className="h-8 w-8 text-primary mb-3" />
                <h3 className="font-semibold text-lg">{feature.title}</h3>
                <p className="mt-2 text-sm text-muted-foreground">{feature.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      <footer className="border-t py-8 px-4">
        <div className="container mx-auto text-center text-sm text-muted-foreground">
          <p>Echo SaaS Starter - Built with Go, React, and TanStack Router</p>
        </div>
      </footer>
    </div>
  )
}
