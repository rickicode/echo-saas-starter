export interface Role {
  id: number
  name: string
  description: string
  created_at: string
}

export const ROLE_SUPER_ADMIN = 'super_admin'
export const ROLE_ADMIN = 'admin'
export const ROLE_USER = 'user'
