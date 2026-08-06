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
  customerId: string
  startsAt: string
  endsAt: string
  status: AppointmentStatus
  notes: string
  createdBy: string
  createdAt: string
  updatedAt: string
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

export type WorkingHours = {
  dayOfWeek: number
  openTime: string
  closeTime: string
  enabled: boolean
}

export type Configuration = {
  tenantId: string
  establishmentName: string
  workingHours: WorkingHours[]
  appointmentIntervalMinutes: number
  defaultDurationMinutes: number
  allowOnlineBooking: boolean
  advanceNoticeHours: number
  updatedAt: string
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
