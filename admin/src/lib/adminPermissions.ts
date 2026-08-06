import type { AdminRole } from './adminAuth'

export type AdminAreaKey =
  | 'agenda'
  | 'appointments'
  | 'customers'
  | 'services'
  | 'professionals'
  | 'configuration'
  | 'team'
  | 'tenants'

export type AdminModuleSetting = {
  key: string
  label: string
  enabled: boolean
}

const ADMIN_AREAS_ROLES: Record<AdminAreaKey, AdminRole[]> = {
  agenda: ['owner', 'admin', 'receptionist', 'professional'],
  appointments: ['owner', 'admin', 'receptionist', 'professional'],
  customers: ['owner', 'admin', 'receptionist'],
  services: ['owner', 'admin'],
  professionals: ['owner', 'admin'],
  configuration: ['owner', 'admin'],
  team: ['owner', 'admin'],
  tenants: ['owner', 'admin'],
}

export function canAccessAdminArea(area: AdminAreaKey, role?: AdminRole): boolean {
  if (!role) {
    return false
  }
  return ADMIN_AREAS_ROLES[area]?.includes(role) ?? false
}

export function isMasterAdmin(tenantId?: string): boolean {
  return tenantId === 'default'
}

export function roleCanManageAdminUsers(role?: AdminRole): boolean {
  return role === 'owner' || role === 'admin'
}
