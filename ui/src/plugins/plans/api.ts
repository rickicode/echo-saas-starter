import { get, post, put, del } from '../../lib/api'
import type {
  Plan,
  UserPlanWithDetails,
  CreatePlanRequest,
  UpdatePlanRequest,
} from './types'

export function getPlans(): Promise<Plan[]> {
  return get<Plan[]>('/plans')
}

export function getPlanById(id: number): Promise<Plan> {
  return get<Plan>(`/plans/${id}`)
}

export function createPlan(data: CreatePlanRequest): Promise<{ message: string }> {
  return post<{ message: string }>('/plans', data)
}

export function updatePlan(id: number, data: UpdatePlanRequest): Promise<{ message: string }> {
  return put<{ message: string }>(`/plans/${id}`, data)
}

export function deletePlan(id: number): Promise<void> {
  return del<void>(`/plans/${id}`)
}

export function assignPlan(planId: number, userId: string): Promise<{ message: string }> {
  return post<{ message: string }>(`/plans/${planId}/assign`, { user_id: userId })
}

export function getCurrentPlan(): Promise<UserPlanWithDetails> {
  return get<UserPlanWithDetails>('/plans/current')
}

export function getAvailablePlans(): Promise<Plan[]> {
  return get<Plan[]>('/plans/available')
}
