import type {
  AdminInviteResponse,
  AdminUser,
  Appointment,
  AppointmentAnalytics,
  AppointmentSummary,
  CalendarConnection,
  CashRegister,
  CashRegisterEntry,
  CashRegisterEntryType,
  CheckoutInput,
  CheckoutResult,
  CommissionRecord,
  CommissionRule,
  CommissionSummary,
  Configuration,
  CurrentCashRegisterResponse,
  Customer,
  FinancialSummary,
  FinancialTransaction,
  GiftCard,
  MemberPlan,
  Membership,
  NotificationLog,
  Package,
  PackageCredit,
  Payment,
  Product,
  Professional,
  Resource,
  Review,
  Service,
  ServiceSupply,
  StockMovement,
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

export async function fetchResources(): Promise<Resource[]> {
  const data = await cachedAdminJson<{ resources: Resource[] }>('/api/v1/admin/resources')
  return data.resources
}

export async function createResource(input: Partial<Resource>): Promise<Resource> {
  const response = await adminFetch('/api/v1/admin/resources', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Resource>(response)
  invalidateAdminGetCache('/api/v1/admin/resources', '/api/v1/admin/appointments')
  return data
}

export async function updateResource(id: string, input: Partial<Resource>): Promise<Resource> {
  const response = await adminFetch(`/api/v1/admin/resources/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Resource>(response)
  invalidateAdminGetCache('/api/v1/admin/resources', '/api/v1/admin/appointments')
  return data
}

export async function setResourceActive(id: string, active: boolean): Promise<Resource> {
  const response = await adminFetch(`/api/v1/admin/resources/${id}/active`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ active }),
  })
  const data = await parseJson<Resource>(response)
  invalidateAdminGetCache('/api/v1/admin/resources', '/api/v1/admin/appointments')
  return data
}

export async function deleteResource(id: string): Promise<void> {
  const response = await adminFetch(`/api/v1/admin/resources/${id}`, { method: 'DELETE' })
  await parseJson<{ status?: string } | null>(response)
  invalidateAdminGetCache('/api/v1/admin/resources', '/api/v1/admin/appointments')
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

export async function fetchReviews(status = ''): Promise<Review[]> {
  const query = status ? `?status=${encodeURIComponent(status)}` : ''
  const data = await cachedAdminJson<{ reviews: Review[] }>(`/api/v1/admin/reviews${query}`)
  return data.reviews
}

export async function downloadAdminCsv(path: string, fallbackName: string): Promise<void> {
  const response = await adminFetch(path)
  if (!response.ok) {
    await parseJson<{ error?: string }>(response)
    return
  }
  const blob = await response.blob()
  const url = URL.createObjectURL(blob)
  const anchor = document.createElement('a')
  anchor.href = url
  anchor.download = fallbackName
  document.body.appendChild(anchor)
  anchor.click()
  anchor.remove()
  URL.revokeObjectURL(url)
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

export async function fetchCurrentTenant(): Promise<Tenant> {
  const data = await cachedAdminJson<Tenant>('/api/v1/admin/tenant')
  return data
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

// ==========================================
// Stock / Inventory & Supplies
// ==========================================

export async function fetchProducts(activeOnly = false): Promise<Product[]> {
  const query = activeOnly ? '?activeOnly=true' : ''
  const data = await cachedAdminJson<{ products: Product[] }>(`/api/v1/admin/products${query}`)
  return data.products || []
}

export async function createProduct(input: Partial<Product>): Promise<Product> {
  const response = await adminFetch('/api/v1/admin/products', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Product>(response)
  invalidateAdminGetCache('/api/v1/admin/products')
  return data
}

export async function updateProduct(id: string, input: Partial<Product>): Promise<Product> {
  const response = await adminFetch(`/api/v1/admin/products/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<Product>(response)
  invalidateAdminGetCache('/api/v1/admin/products')
  return data
}

export async function deleteProduct(id: string): Promise<void> {
  const response = await adminFetch(`/api/v1/admin/products/${id}`, { method: 'DELETE' })
  await parseJson<{ ok?: boolean } | null>(response)
  invalidateAdminGetCache('/api/v1/admin/products')
}

export async function fetchStockMovements(productId?: string, limit = 50): Promise<StockMovement[]> {
  const params = new URLSearchParams()
  if (productId) params.set('productId', productId)
  if (limit) params.set('limit', String(limit))
  const query = params.toString() ? `?${params.toString()}` : ''
  const data = await cachedAdminJson<{ movements: StockMovement[] }>(`/api/v1/admin/stock/movements${query}`)
  return data.movements || []
}

export async function createStockMovement(input: Partial<StockMovement>): Promise<StockMovement> {
  const response = await adminFetch('/api/v1/admin/stock/movements', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<StockMovement>(response)
  invalidateAdminGetCache('/api/v1/admin/stock/movements', '/api/v1/admin/products')
  return data
}

export async function fetchServiceSupplies(serviceId: string): Promise<ServiceSupply[]> {
  const data = await cachedAdminJson<{ supplies: ServiceSupply[] }>(`/api/v1/admin/services/${serviceId}/supplies`)
  return data.supplies || []
}

export async function setServiceSupplies(serviceId: string, supplies: { productId: string; quantity: number }[]): Promise<void> {
  const response = await adminFetch(`/api/v1/admin/services/${serviceId}/supplies`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(supplies),
  })
  await parseJson<{ ok?: boolean }>(response)
  invalidateAdminGetCache(`/api/v1/admin/services/${serviceId}/supplies`)
}

// ==========================================
// Commissions
// ==========================================

export async function fetchCommissionRules(): Promise<CommissionRule[]> {
  const data = await cachedAdminJson<{ rules: CommissionRule[] }>('/api/v1/admin/commissions/rules')
  return data.rules || []
}

export async function createCommissionRule(input: Partial<CommissionRule>): Promise<CommissionRule> {
  const response = await adminFetch('/api/v1/admin/commissions/rules', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<CommissionRule>(response)
  invalidateAdminGetCache('/api/v1/admin/commissions/rules')
  return data
}

export async function updateCommissionRule(id: string, input: Partial<CommissionRule>): Promise<CommissionRule> {
  const response = await adminFetch(`/api/v1/admin/commissions/rules/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<CommissionRule>(response)
  invalidateAdminGetCache('/api/v1/admin/commissions/rules')
  return data
}

export async function deleteCommissionRule(id: string): Promise<void> {
  const response = await adminFetch(`/api/v1/admin/commissions/rules/${id}`, { method: 'DELETE' })
  await parseJson<{ ok?: boolean } | null>(response)
  invalidateAdminGetCache('/api/v1/admin/commissions/rules')
}

export async function fetchCommissionRecords(professionalId?: string, status?: string): Promise<CommissionRecord[]> {
  const params = new URLSearchParams()
  if (professionalId) params.set('professionalId', professionalId)
  if (status) params.set('status', status)
  const query = params.toString() ? `?${params.toString()}` : ''
  const data = await cachedAdminJson<{ records: CommissionRecord[] }>(`/api/v1/admin/commissions${query}`)
  return data.records || []
}

export async function settleCommissions(professionalId: string, recordIds: string[]): Promise<void> {
  const response = await adminFetch('/api/v1/admin/commissions/settle', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ professionalId, recordIds }),
  })
  await parseJson<{ ok?: boolean }>(response)
  invalidateAdminGetCache('/api/v1/admin/commissions', '/api/v1/admin/commissions/summaries')
}

export async function fetchCommissionSummaries(): Promise<CommissionSummary[]> {
  const data = await cachedAdminJson<{ summaries: CommissionSummary[] }>('/api/v1/admin/commissions/summaries')
  return data.summaries || []
}

// ==========================================
// Finance
// ==========================================

export async function fetchFinancialTransactions(filters?: {
  type?: string
  status?: string
  fromDate?: string
  toDate?: string
}): Promise<FinancialTransaction[]> {
  const params = new URLSearchParams()
  if (filters?.type) params.set('type', filters.type)
  if (filters?.status) params.set('status', filters.status)
  if (filters?.fromDate) params.set('fromDate', filters.fromDate)
  if (filters?.toDate) params.set('toDate', filters.toDate)
  const query = params.toString() ? `?${params.toString()}` : ''
  const data = await cachedAdminJson<{ transactions: FinancialTransaction[] }>(`/api/v1/admin/finance/transactions${query}`)
  return data.transactions || []
}

export async function createFinancialTransaction(input: Partial<FinancialTransaction>): Promise<FinancialTransaction> {
  const response = await adminFetch('/api/v1/admin/finance/transactions', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<FinancialTransaction>(response)
  invalidateAdminGetCache('/api/v1/admin/finance/transactions', '/api/v1/admin/finance/summary')
  return data
}

export async function updateFinancialTransaction(id: string, input: Partial<FinancialTransaction>): Promise<FinancialTransaction> {
  const response = await adminFetch(`/api/v1/admin/finance/transactions/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<FinancialTransaction>(response)
  invalidateAdminGetCache('/api/v1/admin/finance/transactions', '/api/v1/admin/finance/summary')
  return data
}

export async function deleteFinancialTransaction(id: string): Promise<void> {
  const response = await adminFetch(`/api/v1/admin/finance/transactions/${id}`, { method: 'DELETE' })
  await parseJson<{ ok?: boolean } | null>(response)
  invalidateAdminGetCache('/api/v1/admin/finance/transactions', '/api/v1/admin/finance/summary')
}

export async function fetchFinancialSummary(fromDate?: string, toDate?: string): Promise<FinancialSummary> {
  const params = new URLSearchParams()
  if (fromDate) params.set('fromDate', fromDate)
  if (toDate) params.set('toDate', toDate)
  const query = params.toString() ? `?${params.toString()}` : ''
  return cachedAdminJson<FinancialSummary>(`/api/v1/admin/finance/summary${query}`)
}

// ==========================================
// Calendar Sync
// ==========================================

export async function fetchCalendarConnections(): Promise<CalendarConnection[]> {
  const data = await cachedAdminJson<{ connections: CalendarConnection[] }>('/api/v1/admin/calendar/connections')
  return data.connections || []
}

export async function connectCalendar(input: Partial<CalendarConnection>): Promise<CalendarConnection> {
  const response = await adminFetch('/api/v1/admin/calendar/connections', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<CalendarConnection>(response)
  invalidateAdminGetCache('/api/v1/admin/calendar/connections')
  return data
}

export async function disconnectCalendar(id: string): Promise<void> {
  const response = await adminFetch(`/api/v1/admin/calendar/connections/${id}`, { method: 'DELETE' })
  await parseJson<{ ok?: boolean } | null>(response)
  invalidateAdminGetCache('/api/v1/admin/calendar/connections')
}

// ==========================================
// Cash Register & POS
// ==========================================

export async function fetchCurrentCashRegister(): Promise<CurrentCashRegisterResponse> {
  return cachedAdminJson<CurrentCashRegisterResponse>('/api/v1/admin/cash-register/current')
}

export async function openCashRegister(initialBalanceCents: number, notes?: string): Promise<CashRegister> {
  const response = await adminFetch('/api/v1/admin/cash-register/open', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ initialBalanceCents, notes: notes || '' }),
  })
  const data = await parseJson<CashRegister>(response)
  invalidateAdminGetCache('/api/v1/admin/cash-register/current', '/api/v1/admin/cash-register/history')
  return data
}

export async function createCashRegisterEntry(
  type: CashRegisterEntryType,
  amountCents: number,
  reason: string
): Promise<CashRegisterEntry> {
  const response = await adminFetch('/api/v1/admin/cash-register/entry', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ type, amountCents, reason }),
  })
  const data = await parseJson<CashRegisterEntry>(response)
  invalidateAdminGetCache(
    '/api/v1/admin/cash-register/current',
    '/api/v1/admin/finance/transactions',
    '/api/v1/admin/finance/summary'
  )
  return data
}

export async function closeCashRegister(closingBalanceCents: number, notes?: string): Promise<CashRegister> {
  const response = await adminFetch('/api/v1/admin/cash-register/close', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ closingBalanceCents, notes: notes || '' }),
  })
  const data = await parseJson<CashRegister>(response)
  invalidateAdminGetCache('/api/v1/admin/cash-register/current', '/api/v1/admin/cash-register/history')
  return data
}

export async function fetchCashRegisterHistory(): Promise<CashRegister[]> {
  const data = await cachedAdminJson<{ history: CashRegister[] }>('/api/v1/admin/cash-register/history')
  return data.history || []
}

// ==========================================
// Checkout / Fechar Comanda
// ==========================================

export async function checkoutAppointment(input: CheckoutInput): Promise<CheckoutResult> {
  const response = await adminFetch('/api/v1/admin/checkout', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const data = await parseJson<CheckoutResult>(response)
  invalidateMutations(
    '/api/v1/admin/appointments',
    '/api/v1/admin/cash-register/current',
    '/api/v1/admin/finance/transactions',
    '/api/v1/admin/finance/summary',
    '/api/v1/admin/stock/products',
    '/api/v1/admin/commissions/records',
    '/api/v1/admin/commissions/summary'
  )
  return data
}
