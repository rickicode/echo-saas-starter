import { createFileRoute, Link } from '@tanstack/react-router'

export const Route = createFileRoute('/_guest/pricing')({
  component: PricingPage,
})

const plans = [
  {
    name: 'Free',
    price: '$0',
    period: '/month',
    description: 'For individuals getting started',
    features: ['1 project', 'Basic support', 'Community access', '1 GB storage'],
    cta: 'Get Started',
    highlighted: false,
  },
  {
    name: 'Pro',
    price: '$29',
    period: '/month',
    description: 'For growing teams and businesses',
    features: ['Unlimited projects', 'Priority support', 'Advanced analytics', '50 GB storage', 'Custom domains', 'API access'],
    cta: 'Start Free Trial',
    highlighted: true,
  },
  {
    name: 'Enterprise',
    price: '$99',
    period: '/month',
    description: 'For large organizations',
    features: ['Everything in Pro', 'Dedicated support', 'SSO/SAML', 'Unlimited storage', 'SLA guarantee', 'Custom integrations'],
    cta: 'Contact Sales',
    highlighted: false,
  },
]

function PricingPage() {
  return (
    <div className="py-16 px-4">
      <div className="container mx-auto max-w-5xl">
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold">Simple, Transparent Pricing</h1>
          <p className="mt-4 text-lg text-muted-foreground">
            Choose the plan that fits your needs
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {plans.map((plan) => (
            <div
              key={plan.name}
              className={`rounded-lg border p-6 flex flex-col ${plan.highlighted ? 'border-primary ring-2 ring-primary' : ''}`}
            >
              {plan.highlighted && (
                <span className="text-xs font-semibold text-primary mb-2">Popular</span>
              )}
              <h3 className="text-xl font-bold">{plan.name}</h3>
              <p className="text-sm text-muted-foreground mt-1">{plan.description}</p>
              <div className="mt-4">
                <span className="text-3xl font-bold">{plan.price}</span>
                <span className="text-muted-foreground">{plan.period}</span>
              </div>
              <ul className="mt-6 space-y-2 flex-1">
                {plan.features.map((feature) => (
                  <li key={feature} className="flex items-center text-sm">
                    <span className="mr-2 text-primary">&#10003;</span>
                    {feature}
                  </li>
                ))}
              </ul>
              <Link
                to="/register"
                className={`mt-6 block text-center rounded-md px-4 py-2 text-sm font-medium ${
                  plan.highlighted
                    ? 'bg-primary text-primary-foreground hover:bg-primary/90'
                    : 'border hover:bg-muted'
                }`}
              >
                {plan.cta}
              </Link>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
