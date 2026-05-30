export interface Plan {
  id: number
  name: string
  slug: string
  description: string
  price_monthly: number
  price_yearly: number
  features: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface UserPlan {
  id: string
  user_id: string
  plan_id: number
  status: string
  starts_at: string
  expires_at?: string
  created_at: string
}

export interface UserPlanWithDetails {
  user_plan: UserPlan
  plan: Plan
}

export interface CreatePlanRequest {
  name: string
  slug: string
  description: string
  price_monthly: number
  price_yearly: number
  features: string
  is_active: boolean
}

export interface UpdatePlanRequest {
  name: string
  slug: string
  description: string
  price_monthly: number
  price_yearly: number
  features: string
  is_active: boolean
}

export interface AssignPlanRequest {
  user_id: string
}
