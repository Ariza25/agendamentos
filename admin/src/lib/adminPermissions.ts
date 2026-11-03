import type { AdminRole } from './adminAuth'

export type AdminAreaKey =
  | 'agenda'
  | 'appointments'
  | 'customers'
  | 'services'
  | 'resources'
  | 'professionals'
  | 'configuration'
  | 'team'
  | 'tenants'
  | 'analytics'
  | 'notifications'
  | 'payments'
  | 'waiting-list'
  | 'reviews'
  | 'packages'
  | 'gift-cards'
  | 'memberships'
  | 'widget'
  | 'finance'
  | 'inventory'
  | 'commissions'

export type AdminModuleSetting = {
  key: string
  label: string
  enabled: boolean
}

const ADMIN_AREAS_ROLES: Record<AdminAreaKey, AdminRole[]> = {
  agenda: ['owner', 'admin', 'receptionist', 'professional', 'master', 'superadmin'],
  appointments: ['owner', 'admin', 'receptionist', 'professional', 'master', 'superadmin'],
  customers: ['owner', 'admin', 'receptionist', 'master', 'superadmin'],
  services: ['owner', 'admin', 'master', 'superadmin'],
  resources: ['owner', 'admin', 'master', 'superadmin'],
  professionals: ['owner', 'admin', 'master', 'superadmin'],
  configuration: ['owner', 'admin', 'master', 'superadmin'],
  team: ['owner', 'admin', 'master', 'superadmin'],
  tenants: ['master', 'superadmin'],
  analytics: ['owner', 'admin', 'master', 'superadmin'],
  notifications: ['owner', 'admin', 'receptionist', 'master', 'superadmin'],
  payments: ['owner', 'admin', 'receptionist', 'master', 'superadmin'],
  'waiting-list': ['owner', 'admin', 'receptionist', 'master', 'superadmin'],
  reviews: ['owner', 'admin', 'receptionist', 'master', 'superadmin'],
  packages: ['owner', 'admin', 'master', 'superadmin'],
  'gift-cards': ['owner', 'admin', 'master', 'superadmin'],
  memberships: ['owner', 'admin', 'master', 'superadmin'],
  widget: ['owner', 'admin', 'master', 'superadmin'],
  finance: ['owner', 'admin', 'master', 'superadmin'],
  inventory: ['owner', 'admin', 'master', 'superadmin'],
  commissions: ['owner', 'admin', 'master', 'superadmin'],
}

export function canAccessAdminArea(area: AdminAreaKey, role?: AdminRole): boolean {
  if (!role) {
    return false
  }
  return ADMIN_AREAS_ROLES[area]?.includes(role) ?? false
}

export function isMasterAdmin(tenantId?: string, role?: string): boolean {
  return tenantId === 'default' || role === 'master' || role === 'superadmin'
}

export function roleCanManageAdminUsers(role?: AdminRole): boolean {
  return role === 'owner' || role === 'admin' || role === 'master' || role === 'superadmin'
}
