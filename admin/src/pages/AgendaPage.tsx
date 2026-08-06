import { ChevronLeft, ChevronRight } from 'lucide-react'
import { useEffect, useState } from 'react'

import { fetchAppointments, fetchAppointmentSummary } from '../lib/adminApi'
import { APPOINTMENT_STATUS_LABELS, APPOINTMENT_STATUS_TONES } from '../lib/appointmentStatus'
import { formatDate, formatTime, todayKey } from '../lib/format'
import { useReferenceData } from '../shared/hooks/useReferenceData'
import { Badge } from '../shared/components/Badge'
import { Button } from '../shared/components/Button'
import { EmptyState } from '../shared/components/EmptyState'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import type { Appointment, AppointmentSummary, AppointmentStatus, Professional } from '../types'

const SUMMARY_ORDER: AppointmentStatus[] = ['scheduled', 'confirmed', 'in_progress', 'completed', 'cancelled', 'no_show']

export function AgendaPage() {
  const [date, setDate] = useState(todayKey())
  const [appointments, setAppointments] = useState<Appointment[]>([])
  const [summary, setSummary] = useState<AppointmentSummary | null>(null)
  const [loading, setLoading] = useState(true)
  const referenceData = useReferenceData()

  useEffect(() => {
    let active = true
    setLoading(true)

    Promise.all([fetchAppointments({ from: `${date}T00:00:00`, to: `${date}T23:59:59` }), fetchAppointmentSummary(date)])
      .then(([items, summaryData]) => {
        if (!active) {
          return
        }
        setAppointments(items)
        setSummary(summaryData)
      })
      .finally(() => {
        if (active) {
          setLoading(false)
        }
      })

    return () => {
      active = false
    }
  }, [date])

  const groups: Array<{ professional: Professional | null; items: Appointment[] }> = referenceData.professionals
    .map((professional) => ({
      professional,
      items: appointments
        .filter((appointment) => appointment.professionalId === professional.id)
        .sort((a, b) => a.startsAt.localeCompare(b.startsAt)),
    }))
    .filter((group) => group.items.length > 0)

  const unlinkedAppointments = appointments.filter(
    (appointment) => !referenceData.professionals.some((professional) => professional.id === appointment.professionalId),
  )
  if (unlinkedAppointments.length > 0) {
    groups.push({
      professional: null,
      items: unlinkedAppointments.sort((a, b) => a.startsAt.localeCompare(b.startsAt)),
    })
  }

  return (
    <PageContainer>
      <Section
        actions={
          <div className="flex items-center gap-2">
            <Button
              aria-label="Dia anterior"
              onClick={() => setDate((current) => shiftDate(current, -1))}
              size="icon"
              variant="secondary"
            >
              <ChevronLeft className="h-5 w-5" />
            </Button>
            <Button onClick={() => setDate(todayKey())} variant="secondary">
              Hoje
            </Button>
            <Button
              aria-label="Próximo dia"
              onClick={() => setDate((current) => shiftDate(current, 1))}
              size="icon"
              variant="secondary"
            >
              <ChevronRight className="h-5 w-5" />
            </Button>
          </div>
        }
        subtitle={formatDate(date)}
        title="Agenda"
      >
        <div className="mt-4 flex flex-wrap gap-2">
          {summary ? (
            <>
              {SUMMARY_ORDER.filter((status) => (summary[status] ?? 0) > 0).map((status) => (
                <Badge key={status} tone={APPOINTMENT_STATUS_TONES[status]}>
                  {APPOINTMENT_STATUS_LABELS[status]}: {summary[status]}
                </Badge>
              ))}
              <Badge tone="neutral">Total: {summary.total}</Badge>
            </>
          ) : (
            <Badge tone="neutral">Carregando resumo...</Badge>
          )}
        </div>
      </Section>

      {loading ? (
        <Section>
          <Spinner />
        </Section>
      ) : groups.length === 0 ? (
        <Section>
          <EmptyState
            description="Não há agendamentos para este dia. Utilize 'Agendamentos' para registrar um novo horário."
            title="Nenhum agendamento neste dia"
          />
        </Section>
      ) : (
        <div className="grid gap-6 lg:grid-cols-2">
          {groups.map(({ professional, items }) => (
            <Section
              key={professional?.id ?? 'sem-profissional'}
              subtitle={`${items.length} agendamento${items.length === 1 ? '' : 's'}`}
              title={professional?.name ?? 'Sem profissional vinculado'}
            >
              <ul className="mt-4 grid gap-3">
                {items.map((appointment) => {
                  const customer = referenceData.byId.customer.get(appointment.customerId)
                  const service = referenceData.byId.service.get(appointment.serviceId)
                  return (
                    <li
                      className="flex items-center gap-4 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] px-4 py-3"
                      key={appointment.id}
                    >
                      <span className="shrink-0 rounded-xl bg-[var(--color-primary)] px-2.5 py-1 font-headline text-sm font-black text-white">
                        {formatTime(appointment.startsAt)}
                      </span>
                      <div className="min-w-0 flex-1">
                        <p className="truncate font-bold text-[var(--color-ink)]">{customer?.name ?? 'Cliente desconhecido'}</p>
                        <p className="truncate text-sm text-[var(--color-ink-soft)]">
                          {service?.name ?? 'Serviço não informado'}
                          {service?.durationMinutes ? ` · ${service.durationMinutes} min` : ''}
                        </p>
                      </div>
                      <Badge tone={APPOINTMENT_STATUS_TONES[appointment.status]}>{APPOINTMENT_STATUS_LABELS[appointment.status]}</Badge>
                    </li>
                  )
                })}
              </ul>
            </Section>
          ))}
        </div>
      )}
    </PageContainer>
  )
}

function shiftDate(dateKey: string, offsetDays: number): string {
  const [year, month, day] = dateKey.split('-').map(Number)
  const date = new Date(year, month - 1, day)
  date.setDate(date.getDate() + offsetDays)
  const nextYear = date.getFullYear()
  const nextMonth = String(date.getMonth() + 1).padStart(2, '0')
  const nextDay = String(date.getDate()).padStart(2, '0')
  return `${nextYear}-${nextMonth}-${nextDay}`
}
