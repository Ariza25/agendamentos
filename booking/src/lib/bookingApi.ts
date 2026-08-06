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
