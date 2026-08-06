import type {
  AdminInviteResponse,
  AdminUser,
  Appointment,
  AppointmentAnalytics,
  AppointmentSummary,
  Configuration,
  Customer,
  GiftCard,
  MemberPlan,
  Membership,
  NotificationLog,
  Package,
  PackageCredit,
  Payment,
  Professional,
  Service,
  Tenant,
  WaitingListEntry,
} from '../types'
import { adminFetch, cachedAdminJson, parseJson, invalidateAdminGetCache } from './adminHttp'

function invalidateMutations(...pathPrefixes: string[]) {
  invalidateAdminGetCache('/api/v1/admin/appointments', ...pathPrefixes)
}

export async function fetchProfessionals(): Promise<Professional[]> {
  const data = await cachedAdminJson<{ professionals: Professional[] }>('/api/v1/admin/professionals')
  return data.professionals
}

export async function createProfessional(input: Partial<Professional>): Promise<Professional> {
  const response = await adminFetch('/api/v1/admin/professionals', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Professional>(response)
  invalidateAdminGetCache('/api/v1/admin/professionals')
  return data
}

export async function updateProfessional(id: string, input: Partial<Professional>): Promise<Professional> {
  const response = await adminFetch(`/api/v1/admin/professionals/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Professional>(response)
  invalidateAdminGetCache('/api/v1/admin/professionals')
  return data
}

export async function setProfessionalActive(id: string, active: boolean): Promise<Professional> {
  const response = await adminFetch(`/api/v1/admin/professionals/${id}/active`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ active }),
  })
  const data = await parseJson<Professional>(response)
  invalidateAdminGetCache('/api/v1/admin/professionals')
  return data
}

export async function deleteProfessional(id: string): Promise<void> {
  const response = await adminFetch(`/api/v1/admin/professionals/${id}`, { method: 'DELETE' })
  await parseJson<{ status?: string } | null>(response)
  invalidateAdminGetCache('/api/v1/admin/professionals')
}

export async function fetchServices(): Promise<Service[]> {
  const data = await cachedAdminJson<{ services: Service[] }>('/api/v1/admin/services')
  return data.services
}

export async function createService(input: Partial<Service>): Promise<Service> {
  const response = await adminFetch('/api/v1/admin/services', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Service>(response)
  invalidateAdminGetCache('/api/v1/admin/services')
  return data
}

export async function updateService(id: string, input: Partial<Service>): Promise<Service> {
  const response = await adminFetch(`/api/v1/admin/services/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Service>(response)
  invalidateAdminGetCache('/api/v1/admin/services')
  return data
}

export async function setServiceActive(id: string, active: boolean): Promise<Service> {
  const response = await adminFetch(`/api/v1/admin/services/${id}/active`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ active }),
  })
  const data = await parseJson<Service>(response)
  invalidateAdminGetCache('/api/v1/admin/services')
  return data
}

export async function deleteService(id: string): Promise<void> {
  const response = await adminFetch(`/api/v1/admin/services/${id}`, { method: 'DELETE' })
  await parseJson<{ status?: string } | null>(response)
  invalidateAdminGetCache('/api/v1/admin/services')
}

export async function fetchCustomers(search = ''): Promise<Customer[]> {
  const query = search.trim() ? `?search=${encodeURIComponent(search.trim())}` : ''
  const data = await cachedAdminJson<{ customers: Customer[] }>(`/api/v1/admin/customers${query}`)
  return data.customers
}

export async function createCustomer(input: Partial<Customer>): Promise<Customer> {
  const response = await adminFetch('/api/v1/admin/customers', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Customer>(response)
  invalidateAdminGetCache('/api/v1/admin/customers')
  return data
}

export async function updateCustomer(id: string, input: Partial<Customer>): Promise<Customer> {
  const response = await adminFetch(`/api/v1/admin/customers/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Customer>(response)
  invalidateAdminGetCache('/api/v1/admin/customers')
  return data
}

export async function deleteCustomer(id: string): Promise<void> {
  const response = await adminFetch(`/api/v1/admin/customers/${id}`, { method: 'DELETE' })
  await parseJson<{ status?: string } | null>(response)
  invalidateAdminGetCache('/api/v1/admin/customers')
}

export type AppointmentFilters = {
  from: string
  to: string
  professionalId?: string
  customerId?: string
  status?: string
}

export async function fetchAppointments(filters: AppointmentFilters): Promise<Appointment[]> {
  const params = new URLSearchParams({ from: filters.from, to: filters.to })
  if (filters.professionalId) params.set('professionalId', filters.professionalId)
  if (filters.customerId) params.set('customerId', filters.customerId)
  if (filters.status) params.set('status', filters.status)
  const data = await cachedAdminJson<{ appointments: Appointment[] }>(`/api/v1/admin/appointments?${params.toString()}`)
  return data.appointments
}

export async function fetchAppointmentSummary(date: string): Promise<AppointmentSummary> {
  return cachedAdminJson<AppointmentSummary>(`/api/v1/admin/appointments/summary?date=${date}`)
}

export async function createAppointment(input: Partial<Appointment>): Promise<Appointment> {
  const response = await adminFetch('/api/v1/admin/appointments', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Appointment>(response)
  invalidateMutations()
  return data
}

export async function updateAppointment(id: string, input: Partial<Appointment>): Promise<Appointment> {
  const response = await adminFetch(`/api/v1/admin/appointments/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Appointment>(response)
  invalidateMutations()
  return data
}

export async function updateAppointmentStatus(id: string, status: string): Promise<Appointment> {
  const response = await adminFetch(`/api/v1/admin/appointments/${id}/status`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ status }),
  })
  const data = await parseJson<Appointment>(response)
  invalidateMutations()
  return data
}

export async function checkInAppointment(id: string): Promise<Appointment> {
  const response = await adminFetch(`/api/v1/admin/appointments/${id}/check-in`, { method: 'POST' })
  const data = await parseJson<Appointment>(response)
  invalidateMutations()
  return data
}

export async function markAppointmentNoShow(id: string, feeCents = 0): Promise<Appointment> {
  const response = await adminFetch(`/api/v1/admin/appointments/${id}/no-show`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ feeCents }),
  })
  const data = await parseJson<Appointment>(response)
  invalidateMutations()
  return data
}

export async function fetchAppointmentAnalytics(from: string, to: string): Promise<AppointmentAnalytics> {
  const params = new URLSearchParams({ from, to })
  return cachedAdminJson<AppointmentAnalytics>(`/api/v1/admin/appointments/analytics?${params.toString()}`)
}

export async function fetchNotifications(limit = 50): Promise<NotificationLog[]> {
  const data = await cachedAdminJson<{ notifications: NotificationLog[] }>(`/api/v1/admin/notifications?limit=${limit}`)
  return data.notifications
}

export async function sendTestNotification(input: {
  channel: string
  recipient: string
  subject: string
  body: string
}): Promise<NotificationLog> {
  const response = await adminFetch('/api/v1/admin/notifications/test', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<NotificationLog>(response)
  invalidateAdminGetCache('/api/v1/admin/notifications')
  return data
}

export async function fetchPayments(limit = 50): Promise<Payment[]> {
  const data = await cachedAdminJson<{ payments: Payment[] }>(`/api/v1/admin/payments?limit=${limit}`)
  return data.payments
}

export async function createPayment(input: {
  appointmentId?: string
  amountCents: number
  method: string
  description?: string
}): Promise<Payment> {
  const response = await adminFetch('/api/v1/admin/payments', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Payment>(response)
  invalidateAdminGetCache('/api/v1/admin/payments')
  return data
}

export async function markPaymentPaid(id: string): Promise<Payment> {
  const response = await adminFetch(`/api/v1/admin/payments/${id}/mark-paid`, { method: 'POST' })
  const data = await parseJson<Payment>(response)
  invalidateAdminGetCache('/api/v1/admin/payments')
  return data
}

export async function refundPayment(id: string): Promise<Payment> {
  const response = await adminFetch(`/api/v1/admin/payments/${id}/refund`, { method: 'POST' })
  const data = await parseJson<Payment>(response)
  invalidateAdminGetCache('/api/v1/admin/payments')
  return data
}

export async function fetchWaitingList(status = ''): Promise<WaitingListEntry[]> {
  const query = status ? `?status=${encodeURIComponent(status)}` : ''
  const data = await cachedAdminJson<{ entries: WaitingListEntry[] }>(`/api/v1/admin/waiting-list${query}`)
  return data.entries
}

export async function updateWaitingListStatus(id: string, status: string): Promise<WaitingListEntry> {
  const response = await adminFetch(`/api/v1/admin/waiting-list/${id}/status`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ status }),
  })
  const data = await parseJson<WaitingListEntry>(response)
  invalidateAdminGetCache('/api/v1/admin/waiting-list')
  return data
}

export async function fetchPackages(): Promise<Package[]> {
  const data = await cachedAdminJson<{ packages: Package[] }>('/api/v1/admin/packages')
  return data.packages
}

export async function createPackage(input: Partial<Package>): Promise<Package> {
  const response = await adminFetch('/api/v1/admin/packages', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Package>(response)
  invalidateAdminGetCache('/api/v1/admin/packages')
  return data
}

export async function updatePackage(id: string, input: Partial<Package>): Promise<Package> {
  const response = await adminFetch(`/api/v1/admin/packages/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Package>(response)
  invalidateAdminGetCache('/api/v1/admin/packages')
  return data
}

export async function setPackageActive(id: string, active: boolean): Promise<Package> {
  const response = await adminFetch(`/api/v1/admin/packages/${id}/active`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ active }),
  })
  const data = await parseJson<Package>(response)
  invalidateAdminGetCache('/api/v1/admin/packages')
  return data
}

export async function deletePackage(id: string): Promise<void> {
  const response = await adminFetch(`/api/v1/admin/packages/${id}`, { method: 'DELETE' })
  await parseJson<{ status?: string } | null>(response)
  invalidateAdminGetCache('/api/v1/admin/packages')
}

export async function fetchGiftCards(): Promise<GiftCard[]> {
  const data = await cachedAdminJson<{ giftCards: GiftCard[] }>('/api/v1/admin/gift-cards')
  return data.giftCards
}

export async function createGiftCard(input: Partial<GiftCard>): Promise<GiftCard> {
  const response = await adminFetch('/api/v1/admin/gift-cards', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<GiftCard>(response)
  invalidateAdminGetCache('/api/v1/admin/gift-cards')
  return data
}

export async function updateGiftCard(id: string, input: Partial<GiftCard>): Promise<GiftCard> {
  const response = await adminFetch(`/api/v1/admin/gift-cards/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<GiftCard>(response)
  invalidateAdminGetCache('/api/v1/admin/gift-cards')
  return data
}

export async function deleteGiftCard(id: string): Promise<void> {
  const response = await adminFetch(`/api/v1/admin/gift-cards/${id}`, { method: 'DELETE' })
  await parseJson<{ status?: string } | null>(response)
  invalidateAdminGetCache('/api/v1/admin/gift-cards')
}

export async function fetchMemberships(): Promise<Membership[]> {
  const data = await cachedAdminJson<{ memberships: Membership[] }>('/api/v1/admin/memberships')
  return data.memberships
}

export async function createMembership(input: Partial<Membership>): Promise<Membership> {
  const response = await adminFetch('/api/v1/admin/memberships', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Membership>(response)
  invalidateAdminGetCache('/api/v1/admin/memberships')
  return data
}

export async function updateMembership(id: string, input: Partial<Membership>): Promise<Membership> {
  const response = await adminFetch(`/api/v1/admin/memberships/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Membership>(response)
  invalidateAdminGetCache('/api/v1/admin/memberships')
  return data
}

export async function setMembershipActive(id: string, active: boolean): Promise<Membership> {
  const response = await adminFetch(`/api/v1/admin/memberships/${id}/active`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ active }),
  })
  const data = await parseJson<Membership>(response)
  invalidateAdminGetCache('/api/v1/admin/memberships')
  return data
}

export async function deleteMembership(id: string): Promise<void> {
  const response = await adminFetch(`/api/v1/admin/memberships/${id}`, { method: 'DELETE' })
  await parseJson<{ status?: string } | null>(response)
  invalidateAdminGetCache('/api/v1/admin/memberships')
}

export async function fetchMemberPlans(): Promise<MemberPlan[]> {
  const data = await cachedAdminJson<{ memberPlans: MemberPlan[] }>('/api/v1/admin/member-plans')
  return data.memberPlans
}

export async function createMemberPlan(input: Partial<MemberPlan>): Promise<MemberPlan> {
  const response = await adminFetch('/api/v1/admin/member-plans', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<MemberPlan>(response)
  invalidateAdminGetCache('/api/v1/admin/member-plans')
  return data
}

export async function updateMemberPlanStatus(id: string, status: string): Promise<MemberPlan> {
  const response = await adminFetch(`/api/v1/admin/member-plans/${id}/status`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ status }),
  })
  const data = await parseJson<MemberPlan>(response)
  invalidateAdminGetCache('/api/v1/admin/member-plans')
  return data
}

export async function fetchPackageCredits(): Promise<PackageCredit[]> {
  const data = await cachedAdminJson<{ packageCredits: PackageCredit[] }>('/api/v1/admin/package-credits')
  return data.packageCredits
}

export async function createPackageCredit(input: {
  packageId: string
  customerId: string
  sessionsTotal: number
}): Promise<PackageCredit> {
  const response = await adminFetch('/api/v1/admin/package-credits', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<PackageCredit>(response)
  invalidateAdminGetCache('/api/v1/admin/package-credits')
  return data
}

export async function deleteAppointment(id: string): Promise<void> {
  const response = await adminFetch(`/api/v1/admin/appointments/${id}`, { method: 'DELETE' })
  await parseJson<{ status?: string } | null>(response)
  invalidateMutations()
}

export async function fetchConfiguration(): Promise<Configuration> {
  return cachedAdminJson<Configuration>('/api/v1/admin/configuration')
}

export async function updateConfiguration(input: Partial<Configuration>): Promise<Configuration> {
  const response = await adminFetch('/api/v1/admin/configuration', {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Configuration>(response)
  invalidateAdminGetCache('/api/v1/admin/configuration')
  return data
}

export async function fetchAdminUsers(): Promise<AdminUser[]> {
  const data = await cachedAdminJson<{ users: AdminUser[] }>('/api/v1/admin/users')
  return data.users
}

export async function createAdminUser(input: {
  name: string
  email: string
  phone?: string
  role: string
}): Promise<AdminInviteResponse> {
  const response = await adminFetch('/api/v1/admin/users', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<AdminInviteResponse>(response)
  invalidateAdminGetCache('/api/v1/admin/users')
  return data
}

export async function updateAdminUserRole(id: string, role: string): Promise<AdminUser> {
  const response = await adminFetch(`/api/v1/admin/users/${id}/role`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ role }),
  })
  const data = await parseJson<AdminUser>(response)
  invalidateAdminGetCache('/api/v1/admin/users')
  return data
}

export async function deleteAdminUser(id: string): Promise<void> {
  const response = await adminFetch(`/api/v1/admin/users/${id}`, { method: 'DELETE' })
  await parseJson<{ status?: string } | null>(response)
  invalidateAdminGetCache('/api/v1/admin/users')
}

export async function fetchTenants(): Promise<Tenant[]> {
  const data = await cachedAdminJson<{ tenants: Tenant[] }>('/api/v1/admin/tenants')
  return data.tenants
}

export async function createTenant(input: Partial<Tenant>): Promise<Tenant> {
  const response = await adminFetch('/api/v1/admin/tenants', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Tenant>(response)
  invalidateAdminGetCache('/api/v1/admin/tenants')
  return data
}

export async function updateTenant(id: string, input: Partial<Tenant>): Promise<Tenant> {
  const response = await adminFetch(`/api/v1/admin/tenants/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Tenant>(response)
  invalidateAdminGetCache('/api/v1/admin/tenants')
  return data
}

export async function setTenantActive(id: string, active: boolean): Promise<Tenant> {
  const response = await adminFetch(`/api/v1/admin/tenants/${id}/active`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ active }),
  })
  const data = await parseJson<Tenant>(response)
  invalidateAdminGetCache('/api/v1/admin/tenants')
  return data
}

export async function deleteTenant(id: string): Promise<void> {
  const response = await adminFetch(`/api/v1/admin/tenants/${id}`, { method: 'DELETE' })
  await parseJson<{ status?: string } | null>(response)
  invalidateAdminGetCache('/api/v1/admin/tenants')
}
