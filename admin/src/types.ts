export type AppointmentStatus =
  | 'scheduled'
  | 'confirmed'
  | 'in_progress'
  | 'completed'
  | 'cancelled'
  | 'no_show'

export type Appointment = {
  id: string
  tenantId: string
  professionalId: string
  serviceId: string
  resourceId?: string
  customerId: string
  startsAt: string
  endsAt: string
  status: AppointmentStatus
  notes: string
  createdBy: string
  seriesId?: string
  paymentStatus?: string
  managementCode?: string
  noShowFeeCents?: number
  source?: string
  cancellationReason?: string
  checkInAt?: string
  noShowAt?: string
  reminderSentAt?: string
  createdAt: string
  updatedAt: string
}

export type AppointmentAnalytics = {
  totalAppointments: number
  scheduled: number
  confirmed: number
  completed: number
  cancelled: number
  noShow: number
  noShowRate: number
  revenueCents: number
  byService: { serviceId: string; serviceName: string; count: number; revenueCents: number }[]
  byProfessional: { professionalId: string; professionalName: string; count: number }[]
  byDay: { date: string; count: number }[]
}

export type AppointmentSummary = {
  date: string
  total: number
  scheduled?: number
  confirmed?: number
  in_progress?: number
  completed?: number
  cancelled?: number
  no_show?: number
}

export type Customer = {
  id: string
  tenantId: string
  name: string
  phone: string
  email: string
  birthday: string
  notes: string
  createdAt: string
  updatedAt: string
}

export type Service = {
  id: string
  tenantId: string
  name: string
  description: string
  durationMinutes: number
  priceCents: number
  color: string
  active: boolean
  createdAt: string
  updatedAt: string
}

export type Professional = {
  id: string
  tenantId: string
  userId: string
  name: string
  email: string
  phone: string
  color: string
  active: boolean
  createdAt: string
  updatedAt: string
}

export type Resource = {
  id: string
  tenantId: string
  name: string
  type: string
  capacity: number
  active: boolean
  notes: string
  createdAt: string
  updatedAt: string
}

export type WorkingHours = {
  dayOfWeek: number
  openTime: string
  closeTime: string
  enabled: boolean
}

export type ClosedDay = {
  date: string
  name: string
  allDay: boolean
  openTime?: string
  closeTime?: string
}

export type Configuration = {
  tenantId: string
  establishmentName: string
  workingHours: WorkingHours[]
  appointmentIntervalMinutes: number
  defaultDurationMinutes: number
  allowOnlineBooking: boolean
  advanceNoticeHours: number
  bufferMinutes: number
  closedDays: ClosedDay[]
  notificationsEnabled: boolean
  notificationReminderHours: number
  notificationChannel: string
  depositEnabled: boolean
  depositType: string
  depositAmountCents: number
  depositPercentage: number
  noShowFeeCents: number
  cancellationPolicyHours: number
  calendarSyncEnabled: boolean
  updatedAt: string
}

export type NotificationLog = {
  id: string
  tenantId: string
  appointmentId?: string
  type: string
  channel: string
  recipient: string
  subject: string
  body: string
  status: string
  error?: string
  createdAt: string
  sentAt: string
}

export type Payment = {
  id: string
  tenantId: string
  appointmentId?: string
  amountCents: number
  method: string
  provider: string
  providerPaymentId?: string
  status: string
  description?: string
  paidAt?: string
  createdAt: string
}

export type Review = {
  id: string
  tenantId: string
  tenantSlug?: string
  appointmentId?: string
  customerId?: string
  customerName?: string
  status: string
  rating: number
  comment: string
  invitedAt: string
  submittedAt?: string
  createdAt: string
}

export type WaitingListEntry = {
  id: string
  tenantId: string
  serviceId?: string
  professionalId?: string
  customerId?: string
  name: string
  phone: string
  status: string
  createdAt: string
  notifiedAt?: string
}

export type Package = {
  id: string
  tenantId: string
  name: string
  description: string
  priceCents: number
  sessionsCount: number
  durationDays: number
  active: boolean
  createdAt: string
  updatedAt: string
}

export type PackageCredit = {
  id: string
  tenantId: string
  packageId?: string
  packageName?: string
  customerId: string
  customerName?: string
  sessionsTotal: number
  sessionsUsed: number
  expiresAt?: string
  createdAt: string
}

export type GiftCard = {
  id: string
  tenantId: string
  code: string
  valueCents: number
  balanceCents: number
  customerId?: string
  customerName?: string
  status: string
  expiresAt?: string
  createdAt: string
}

export type Membership = {
  id: string
  tenantId: string
  name: string
  description: string
  priceCents: number
  billingPeriodDays: number
  active: boolean
  createdAt: string
  updatedAt: string
}

export type MemberPlan = {
  id: string
  tenantId: string
  membershipId?: string
  membershipName?: string
  customerId: string
  customerName?: string
  status: string
  startedAt: string
  expiresAt?: string
  createdAt: string
}

export type AdminUser = {
  id: string
  name: string
  email: string
  phone: string
  role: string
  tenantId: string
  mfaEnabled: boolean
  createdAt: string
}

export type AdminInviteResponse = {
  user: AdminUser
  inviteUrl: string
  inviteEmailSent: boolean
  temporaryPassword?: string
}

export type Tenant = {
  id: string
  name: string
  slug: string
  active: boolean
  createdAt: string
  updatedAt: string
}
