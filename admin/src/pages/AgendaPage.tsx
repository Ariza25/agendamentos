import { ChevronLeft, ChevronRight, Receipt } from 'lucide-react'
import { useEffect, useState } from 'react'

import { checkInAppointment, fetchAppointments, fetchAppointmentSummary, markAppointmentNoShow } from '../lib/adminApi'
import { APPOINTMENT_STATUS_LABELS, APPOINTMENT_STATUS_TONES } from '../lib/appointmentStatus'
import { formatDate, formatTime, fromDateKey, toDateKey, todayKey } from '../lib/format'
import { CheckoutComandaModal } from '../features/appointments/components/CheckoutComandaModal'
import { useReferenceData } from '../shared/hooks/useReferenceData'
import { useAdminEvents } from '../shared/hooks/useAdminEvents'
import { Badge } from '../shared/components/Badge'
import { Button } from '../shared/components/Button'
import { EmptyState } from '../shared/components/EmptyState'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { Appointment, AppointmentSummary, AppointmentStatus, Professional } from '../types'

const SUMMARY_ORDER: AppointmentStatus[] = ['scheduled', 'confirmed', 'in_progress', 'completed', 'cancelled', 'no_show']

export function AgendaPage() {
  const { showToast } = useToast()
  const [date, setDate] = useState(todayKey())
  const [appointments, setAppointments] = useState<Appointment[]>([])
  const [summary, setSummary] = useState<AppointmentSummary | null>(null)
  const [loading, setLoading] = useState(true)
  const [checkoutAppointment, setCheckoutAppointment] = useState<Appointment | null>(null)
  const referenceData = useReferenceData()

  const { isConnected } = useAdminEvents((event) => {
    if (event.type.startsWith('appointment.')) {
      load()
    }
  })

  function load() {
    setLoading(true)

    Promise.all([fetchAppointments({ from: `${date}T00:00:00`, to: `${date}T23:59:59` }), fetchAppointmentSummary(date)])
      .then(([items, summaryData]) => {
        setAppointments(items)
        setSummary(summaryData)
      })
      .catch(() => {
        showToast({ message: 'Não foi possível carregar a agenda.', tone: 'error' })
      })
      .finally(() => {
        setLoading(false)
      })
  }

  useEffect(load, [date])

  async function handleCheckIn(appointment: Appointment) {
    try {
      await checkInAppointment(appointment.id)
      showToast({ message: 'Cliente fez check-in.', tone: 'success' })
      load()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    }
  }

  async function handleNoShow(appointment: Appointment) {
    try {
      await markAppointmentNoShow(appointment.id)
      showToast({ message: 'Agendamento marcado como não comparecimento.', tone: 'success' })
      load()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    }
  }

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
        subtitle={
          <span className="flex items-center gap-3">
            <span>{formatDate(date)}</span>
            {isConnected ? (
              <span className="inline-flex items-center gap-1.5 rounded-full bg-emerald-50 px-2 py-0.5 text-xs font-medium text-emerald-700 ring-1 ring-inset ring-emerald-600/20">
                <span className="h-1.5 w-1.5 rounded-full bg-emerald-500 animate-pulse" />
                Tempo real ativo
              </span>
            ) : null}
          </span>
        }
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
                  const resource = appointment.resourceId ? referenceData.byId.resource.get(appointment.resourceId) : undefined
                  return (
                    <li
                      className="flex flex-col gap-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] px-4 py-3 md:flex-row md:items-center"
                      key={appointment.id}
                    >
                      <div className="flex min-w-0 flex-1 items-center gap-4">
                        <span className="shrink-0 rounded-xl bg-[var(--color-primary)] px-2.5 py-1 font-headline text-sm font-black text-white">
                          {formatTime(appointment.startsAt)}
                        </span>
                        <div className="min-w-0 flex-1">
                          <p className="truncate font-bold text-[var(--color-ink)]">{customer?.name ?? 'Cliente desconhecido'}</p>
                          <p className="truncate text-sm text-[var(--color-ink-soft)]">
                            {service?.name ?? 'Serviço não informado'}
                            {service?.durationMinutes ? ` · ${service.durationMinutes} min` : ''}
                            {resource ? ` · ${resource.name}` : ''}
                          </p>
                        </div>
                      </div>
                      <div className="flex shrink-0 flex-wrap items-center gap-2">
                        {appointment.checkInAt ? <Badge tone="success">Check-in realizado</Badge> : null}
                        <Badge tone={APPOINTMENT_STATUS_TONES[appointment.status]}>
                          {APPOINTMENT_STATUS_LABELS[appointment.status]}
                        </Badge>
                        {appointment.status !== 'completed' && appointment.status !== 'cancelled' ? (
                          <Button
                            onClick={() => setCheckoutAppointment(appointment)}
                            size="sm"
                            className="bg-emerald-600 text-white hover:bg-emerald-700"
                            icon={<Receipt className="h-3.5 w-3.5" />}
                          >
                            Fechar Comanda
                          </Button>
                        ) : null}
                        {appointment.status === 'confirmed' || appointment.status === 'in_progress' ? (
                          <Button onClick={() => void handleCheckIn(appointment)} size="sm" variant="secondary">
                            Check-in
                          </Button>
                        ) : null}
                        {appointment.status === 'scheduled' || appointment.status === 'confirmed' ? (
                          <Button onClick={() => void handleNoShow(appointment)} size="sm" variant="danger">
                            No-show
                          </Button>
                        ) : null}
                      </div>
                    </li>
                  )
                })}
              </ul>
            </Section>
          ))}
        </div>
      )}

      <CheckoutComandaModal
        open={Boolean(checkoutAppointment)}
        onClose={() => setCheckoutAppointment(null)}
        appointment={checkoutAppointment}
        customer={checkoutAppointment ? referenceData.byId.customer.get(checkoutAppointment.customerId) : undefined}
        service={checkoutAppointment ? referenceData.byId.service.get(checkoutAppointment.serviceId) : undefined}
        professional={checkoutAppointment ? referenceData.byId.professional.get(checkoutAppointment.professionalId) : undefined}
        onSuccess={load}
      />
    </PageContainer>
  )
}

function shiftDate(dateKey: string, offsetDays: number): string {
  const date = fromDateKey(dateKey)
  date.setDate(date.getDate() + offsetDays)
  return toDateKey(date)
}
