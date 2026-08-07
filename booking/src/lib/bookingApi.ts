import { buildApiUrl } from './apiBase'
import { parseApiResponse } from './apiError'

export type Service = {
  id: string
  name: string
  description: string
  durationMinutes: number
  priceCents: number
  color: string
  active: boolean
  professionalIds?: string[]
}

export type Professional = {
  id: string
  name: string
  email: string
  phone: string
  color: string
  active: boolean
}

export type WorkingHours = {
  dayOfWeek: number
  openTime: string
  closeTime: string
  enabled: boolean
}

export type PublicConfiguration = {
  establishmentName: string
  workingHours: WorkingHours[]
  appointmentIntervalMinutes: number
  defaultDurationMinutes: number
  allowOnlineBooking: boolean
  advanceNoticeHours: number
}

export type PublicInfo = {
  tenant: { id: string; name: string; slug: string }
  configuration: PublicConfiguration
}

export type Slot = {
  startsAt: string
  endsAt: string
}

export type AvailabilityResponse =
  | { professionalId: string; slots: Slot[] }
  | { professionals: Array<{ professionalId: string; professionalName: string; slots: Slot[] }> }

export type BookingResult = {
  appointment: {
    id: string
    startsAt: string
    endsAt: string
    status: string
    professionalId: string
    serviceId: string
  }
  customer: {
    id: string
    name: string
    phone: string
    email: string
  }
  managementCode: string
}

export type PublicAppointment = {
  id: string
  tenantId: string
  professionalId: string
  serviceId: string
  customerId: string
  startsAt: string
  endsAt: string
  status: string
  notes: string
  source?: string
  cancellationReason?: string
  createdAt: string
  updatedAt: string
}

async function publicFetch<T>(path: string, init?: RequestInit, fallbackMessage = 'Falha na requisicao.'): Promise<T> {
  const response = await fetch(buildApiUrl(path), {
    ...init,
    headers: { 'Content-Type': 'application/json', ...(init?.headers ?? {}) },
  })
  return parseApiResponse<T>(response, fallbackMessage)
}

export function fetchPublicInfo(slug: string): Promise<PublicInfo> {
  return publicFetch(`/api/v1/public/${encodeURIComponent(slug)}/info`, undefined, 'Nao foi possivel carregar o estabelecimento.')
}

export function fetchPublicServices(slug: string): Promise<Service[]> {
  return publicFetch<{ services: Service[] }>(
    `/api/v1/public/${encodeURIComponent(slug)}/services`,
    undefined,
    'Nao foi possivel carregar os servicos.',
  ).then((data) => data.services)
}

export function fetchPublicProfessionals(slug: string): Promise<Professional[]> {
  return publicFetch<{ professionals: Professional[] }>(
    `/api/v1/public/${encodeURIComponent(slug)}/professionals`,
    undefined,
    'Nao foi possivel carregar os profissionais.',
  ).then((data) => data.professionals)
}

export function fetchAvailability(slug: string, date: string, serviceId: string, professionalId: string): Promise<AvailabilityResponse> {
  const params = new URLSearchParams({ date, serviceId, professionalId })
  return publicFetch<AvailabilityResponse>(
    `/api/v1/public/${encodeURIComponent(slug)}/availability?${params.toString()}`,
    undefined,
    'Nao foi possivel carregar os horarios disponiveis.',
  )
}

export function createPublicBooking(
  slug: string,
  input: {
    serviceId: string
    professionalId: string
    startsAt: string
    name: string
    phone: string
    email: string
    notes: string
  },
): Promise<BookingResult> {
  return publicFetch<BookingResult>(
    `/api/v1/public/${encodeURIComponent(slug)}/appointments`,
    { method: 'POST', body: JSON.stringify(input) },
    'Nao foi possivel concluir o agendamento.',
  )
}

export function fetchPublicAppointment(slug: string, code: string): Promise<PublicAppointment> {
  return publicFetch<PublicAppointment>(
    `/api/v1/public/${encodeURIComponent(slug)}/appointments/${encodeURIComponent(code)}`,
    undefined,
    'Nao foi possivel localizar o agendamento. Confira o codigo.',
  )
}

export function cancelPublicAppointment(slug: string, code: string, reason: string): Promise<PublicAppointment> {
  return publicFetch<PublicAppointment>(
    `/api/v1/public/${encodeURIComponent(slug)}/appointments/${encodeURIComponent(code)}/cancel`,
    { method: 'PATCH', body: JSON.stringify({ reason }) },
    'Nao foi possivel cancelar o agendamento.',
  )
}

export function reschedulePublicAppointment(slug: string, code: string, startsAt: string): Promise<PublicAppointment> {
  return publicFetch<PublicAppointment>(
    `/api/v1/public/${encodeURIComponent(slug)}/appointments/${encodeURIComponent(code)}/reschedule`,
    { method: 'PATCH', body: JSON.stringify({ startsAt }) },
    'Nao foi possivel reagendar. Tente novamente.',
  )
}

export function signUpWaitingList(
  slug: string,
  input: {
    serviceId: string
    professionalId?: string
    name: string
    phone: string
  },
): Promise<{ id: string }> {
  return publicFetch<{ id: string }>(
    `/api/v1/public/${encodeURIComponent(slug)}/waiting-list`,
    { method: 'POST', body: JSON.stringify(input) },
    'Nao foi possivel entrar na lista de espera.',
  )
}

export type PublicReview = {
  id: string
  status: string
  rating: number
  comment: string
  invitedAt: string
  submittedAt?: string
}

export function fetchPublicReview(slug: string, reviewId: string): Promise<PublicReview> {
  return publicFetch<PublicReview>(
    `/api/v1/public/${encodeURIComponent(slug)}/reviews/${encodeURIComponent(reviewId)}`,
    undefined,
    'Nao foi possivel localizar a avaliacao.',
  )
}

export function submitPublicReview(slug: string, reviewId: string, rating: number, comment: string): Promise<PublicReview> {
  return publicFetch<PublicReview>(
    `/api/v1/public/${encodeURIComponent(slug)}/reviews/${encodeURIComponent(reviewId)}`,
    { method: 'POST', body: JSON.stringify({ rating, comment }) },
    'Nao foi possivel enviar a avaliacao.',
  )
}

export type PortalAppointment = {
  id: string
  serviceId: string
  serviceName: string
  professionalId: string
  professionalName: string
  resourceId?: string
  resourceName?: string
  startsAt: string
  endsAt: string
  status: string
  priceCents: number
  cancellationReason?: string
  source?: string
}

export type PortalResult = {
  customer: {
    id: string
    name: string
    phone: string
    email: string
  }
  appointments: PortalAppointment[]
}

export function portalLookup(slug: string, phone: string): Promise<PortalResult> {
  return publicFetch<PortalResult>(
    `/api/v1/public/${encodeURIComponent(slug)}/portal/lookup`,
    { method: 'POST', body: JSON.stringify({ phone }) },
    'Nao foi possivel consultar seus agendamentos.',
  )
}

export function portalCancelAppointment(slug: string, phone: string, appointmentId: string, reason: string): Promise<PublicAppointment> {
  return publicFetch<PublicAppointment>(
    `/api/v1/public/${encodeURIComponent(slug)}/portal/appointments/${encodeURIComponent(appointmentId)}/cancel`,
    { method: 'POST', body: JSON.stringify({ phone, reason }) },
    'Nao foi possivel cancelar o agendamento.',
  )
}

export function portalRescheduleAppointment(slug: string, phone: string, appointmentId: string, startsAt: string): Promise<PublicAppointment> {
  return publicFetch<PublicAppointment>(
    `/api/v1/public/${encodeURIComponent(slug)}/portal/appointments/${encodeURIComponent(appointmentId)}/reschedule`,
    { method: 'POST', body: JSON.stringify({ phone, startsAt }) },
    'Nao foi possivel reagendar. Tente novamente.',
  )
}
