import type { AppointmentStatus } from '../types'

export const APPOINTMENT_STATUS_LABELS: Record<AppointmentStatus, string> = {
  scheduled: 'Agendado',
  confirmed: 'Confirmado',
  in_progress: 'Em andamento',
  completed: 'Concluído',
  cancelled: 'Cancelado',
  no_show: 'Não compareceu',
}

export const APPOINTMENT_STATUS_TONES: Record<AppointmentStatus, 'neutral' | 'success' | 'warning' | 'danger' | 'accent'> = {
  scheduled: 'neutral',
  confirmed: 'success',
  in_progress: 'accent',
  completed: 'success',
  cancelled: 'danger',
  no_show: 'warning',
}

export function nextStatusesFor(status: AppointmentStatus): AppointmentStatus[] {
  if (status === 'scheduled' || status === 'confirmed') {
    return ['confirmed', 'in_progress', 'completed', 'cancelled', 'no_show']
  }
  if (status === 'in_progress') {
    return ['completed', 'cancelled']
  }
  return []
}
