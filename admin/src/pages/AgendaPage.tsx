import { Calendar, ChevronLeft, ChevronRight, Receipt, RotateCcw } from 'lucide-react'
import { useEffect, useState } from 'react'

import {
  checkInAppointment,
  fetchAppointments,
  fetchAppointmentSummary,
  markAppointmentNoShow,
  updateAppointmentStatus,
} from '../lib/adminApi'
import { getAdminAccount } from '../lib/adminAuth'
import { APPOINTMENT_STATUS_LABELS, APPOINTMENT_STATUS_TONES, nextStatusesFor } from '../lib/appointmentStatus'
import {
  formatDate,
  formatDateShort,
  formatRelativeDay,
  formatTime,
  fromDateKey,
  SHORT_DAY_LABELS,
  toDateKey,
  todayKey,
} from '../lib/format'
import { CheckoutComandaModal } from '../features/appointments/components/CheckoutComandaModal'
import { useReferenceData } from '../shared/hooks/useReferenceData'
import { useAdminEvents } from '../shared/hooks/useAdminEvents'
import { Badge } from '../shared/components/Badge'
import { Button } from '../shared/components/Button'
import { EmptyState } from '../shared/components/EmptyState'
import { Select } from '../shared/components/FormControls'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { Appointment, AppointmentSummary, AppointmentStatus, Professional } from '../types'

const SUMMARY_ORDER: AppointmentStatus[] = ['scheduled', 'confirmed', 'in_progress', 'completed', 'cancelled', 'no_show']

export function AgendaPage() {
  const { showToast } = useToast()
  const currentUser = getAdminAccount()
  const referenceData = useReferenceData()

  const [date, setDate] = useState(todayKey())
  const [appointments, setAppointments] = useState<Appointment[]>([])
  const [summary, setSummary] = useState<AppointmentSummary | null>(null)
  const [loading, setLoading] = useState(true)
  const [checkoutAppointment, setCheckoutAppointment] = useState<Appointment | null>(null)
  const [selectedProfessionalId, setSelectedProfessionalId] = useState<string>('')
  const [selectedStatus, setSelectedStatus] = useState<string>('')

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
      showToast({ message: 'Check-in realizado com sucesso!', tone: 'success' })
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

  async function handleStatusChange(appointment: Appointment, nextStatus: string) {
    if (!nextStatus) return
    try {
      await updateAppointmentStatus(appointment.id, nextStatus)
      showToast({ message: 'Status do agendamento atualizado.', tone: 'success' })
      load()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    }
  }

  // If a professional account is logged in, restrict to that professional
  const loggedInProfessional =
    currentUser?.role === 'professional'
      ? referenceData.professionals.find((p) => p.email.toLowerCase() === currentUser?.email.toLowerCase())
      : null

  const activeProfessionalFilter = loggedInProfessional ? loggedInProfessional.id : selectedProfessionalId

  // Filter appointments by selected status and professional
  const filteredAppointments = appointments.filter((appointment) => {
    if (selectedStatus && appointment.status !== selectedStatus) {
      return false
    }
    if (activeProfessionalFilter && appointment.professionalId !== activeProfessionalFilter) {
      return false
    }
    return true
  })

  // Group appointments by professional
  const relevantProfessionals = activeProfessionalFilter
    ? referenceData.professionals.filter((p) => p.id === activeProfessionalFilter)
    : referenceData.professionals

  const groups: Array<{ professional: Professional | null; items: Appointment[] }> = relevantProfessionals
    .map((professional) => ({
      professional,
      items: filteredAppointments
        .filter((appointment) => appointment.professionalId === professional.id)
        .sort((a, b) => a.startsAt.localeCompare(b.startsAt)),
    }))
    .filter((group) => group.items.length > 0)

  if (!activeProfessionalFilter) {
    const unlinkedAppointments = filteredAppointments.filter(
      (appointment) => !referenceData.professionals.some((professional) => professional.id === appointment.professionalId),
    )
    if (unlinkedAppointments.length > 0) {
      groups.push({
        professional: null,
        items: unlinkedAppointments.sort((a, b) => a.startsAt.localeCompare(b.startsAt)),
      })
    }
  }

  // Generate 7-day quick strip centered around the selected date
  const daysStrip = [-3, -2, -1, 0, 1, 2, 3].map((offset) => {
    const d = fromDateKey(date)
    d.setDate(d.getDate() + offset)
    const key = toDateKey(d)
    const dayOfWeek = SHORT_DAY_LABELS[d.getDay()]
    const dayNumber = String(d.getDate()).padStart(2, '0')
    const isToday = key === todayKey()
    const isSelected = key === date
    return { key, dayOfWeek, dayNumber, isToday, isSelected }
  })

  const relativeLabel = formatRelativeDay(date)
  const isToday = date === todayKey()

  return (
    <PageContainer>
      <Section
        actions={
          <div className="flex flex-wrap items-center gap-2">
            {!isToday ? (
              <Button
                onClick={() => setDate(todayKey())}
                size="sm"
                variant="secondary"
                className="font-bold text-[var(--color-primary)] border-blue-200 hover:bg-blue-50"
              >
                <RotateCcw className="h-3.5 w-3.5 mr-1" />
                Hoje
              </Button>
            ) : null}

            <div className="flex items-center rounded-xl border border-[var(--color-line)] bg-white p-0.5 shadow-sm">
              <Button
                aria-label="Dia anterior"
                onClick={() => setDate((current) => shiftDate(current, -1))}
                size="icon"
                variant="ghost"
                className="h-9 w-9 text-[var(--color-ink)] hover:bg-slate-100"
              >
                <ChevronLeft className="h-4 w-4" />
              </Button>

              {/* Native Date Picker trigger with styled button */}
              <label
                title="Clique para escolher uma data no calendário"
                className="relative flex cursor-pointer items-center gap-1.5 px-3 py-1 text-sm font-bold text-[var(--color-ink)] hover:text-[var(--color-primary)] transition select-none"
              >
                <Calendar className="h-4 w-4 text-[var(--color-primary)] shrink-0" />
                <span>{formatDateShort(date)}</span>
                <input
                  type="date"
                  value={date}
                  onChange={(e) => {
                    if (e.target.value) setDate(e.target.value)
                  }}
                  className="absolute inset-0 h-full w-full cursor-pointer opacity-0"
                />
              </label>

              <Button
                aria-label="Próximo dia"
                onClick={() => setDate((current) => shiftDate(current, 1))}
                size="icon"
                variant="ghost"
                className="h-9 w-9 text-[var(--color-ink)] hover:bg-slate-100"
              >
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          </div>
        }
        subtitle={
          <span className="flex flex-wrap items-center gap-2.5">
            {relativeLabel ? (
              <span className="rounded-full bg-blue-100 px-2.5 py-0.5 text-xs font-black text-blue-900">
                {relativeLabel}
              </span>
            ) : null}
            <span className="font-bold text-[var(--color-ink)]">{formatDate(date)}</span>
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
        {/* 7-Day Quick Strip for 1-click day selection */}
        <div className="mt-4 flex items-center gap-2 overflow-x-auto pb-1 sm:justify-start">
          {daysStrip.map((day) => (
            <button
              key={day.key}
              type="button"
              onClick={() => setDate(day.key)}
              className={[
                'flex flex-col items-center justify-center min-w-[3.8rem] py-2 px-3 rounded-2xl transition border text-xs font-bold shrink-0 cursor-pointer',
                day.isSelected
                  ? 'bg-[var(--color-primary)] text-white border-[var(--color-primary)] shadow-sm ring-2 ring-[var(--color-primary)]/20'
                  : 'bg-white hover:bg-slate-50 text-[var(--color-ink)] border-[var(--color-line)]',
              ].join(' ')}
            >
              <span
                className={
                  day.isSelected
                    ? 'text-blue-100 uppercase text-[10px]'
                    : 'text-[var(--color-ink-soft)] uppercase text-[10px]'
                }
              >
                {day.dayOfWeek}
              </span>
              <span className="text-base font-black leading-tight mt-0.5">{day.dayNumber}</span>
              {day.isToday ? (
                <span
                  className={[
                    'text-[9px] font-extrabold mt-0.5 px-1 rounded',
                    day.isSelected ? 'bg-white/20 text-white' : 'text-blue-600 bg-blue-50',
                  ].join(' ')}
                >
                  Hoje
                </span>
              ) : null}
            </button>
          ))}
        </div>

        {/* Filters and Status Overview */}
        <div className="mt-4 flex flex-wrap items-center justify-between gap-3 border-t border-[var(--color-line)] pt-3">
          <div className="flex flex-wrap items-center gap-1.5">
            <button
              type="button"
              onClick={() => setSelectedStatus('')}
              className={[
                'px-2.5 py-1 rounded-xl text-xs font-bold transition border cursor-pointer',
                !selectedStatus
                  ? 'bg-[var(--color-primary)] text-white border-[var(--color-primary)]'
                  : 'bg-slate-100 text-slate-700 border-slate-200 hover:bg-slate-200',
              ].join(' ')}
            >
              Todos ({summary?.total ?? 0})
            </button>
            {summary &&
              SUMMARY_ORDER.map((status) => {
                const count = summary[status] ?? 0
                if (count === 0 && selectedStatus !== status) return null
                const isCurrent = selectedStatus === status
                return (
                  <button
                    key={status}
                    type="button"
                    onClick={() => setSelectedStatus(isCurrent ? '' : status)}
                    className={[
                      'rounded-xl text-xs font-bold transition cursor-pointer',
                      isCurrent ? 'ring-2 ring-offset-1 ring-[var(--color-primary)]' : 'opacity-80 hover:opacity-100',
                    ].join(' ')}
                  >
                    <Badge tone={APPOINTMENT_STATUS_TONES[status]}>
                      {APPOINTMENT_STATUS_LABELS[status]}: {count}
                    </Badge>
                  </button>
                )
              })}
          </div>

          {referenceData.professionals.length > 1 && !loggedInProfessional ? (
            <div className="flex items-center gap-2">
              <span className="text-xs font-bold text-[var(--color-ink-soft)] uppercase">Profissional:</span>
              <select
                aria-label="Filtrar por profissional"
                value={selectedProfessionalId}
                onChange={(e) => setSelectedProfessionalId(e.target.value)}
                className="h-8 rounded-lg border border-[var(--color-line)] bg-white px-2.5 text-xs font-bold text-[var(--color-ink)] outline-none"
              >
                <option value="">Todos os profissionais</option>
                {referenceData.professionals.map((p) => (
                  <option key={p.id} value={p.id}>
                    {p.name}
                  </option>
                ))}
              </select>
            </div>
          ) : null}
        </div>
      </Section>

      {loading ? (
        <Section>
          <Spinner />
        </Section>
      ) : groups.length === 0 ? (
        <Section>
          <EmptyState
            description="Não há agendamentos para o filtro selecionado neste dia. Crie um novo agendamento ou navegue para outro dia."
            title="Nenhum agendamento encontrado"
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
                  const nextStatuses = nextStatusesFor(appointment.status)

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
                            className="bg-emerald-600 text-white hover:bg-emerald-700 shadow-sm font-bold"
                            icon={<Receipt className="h-3.5 w-3.5" />}
                          >
                            Fechar Comanda
                          </Button>
                        ) : null}

                        {/* Check-in button available for scheduled/confirmed/in_progress if not yet checked in */}
                        {!appointment.checkInAt &&
                        appointment.status !== 'completed' &&
                        appointment.status !== 'cancelled' &&
                        appointment.status !== 'no_show' ? (
                          <Button onClick={() => void handleCheckIn(appointment)} size="sm" variant="secondary">
                            Check-in
                          </Button>
                        ) : null}

                        {appointment.status === 'scheduled' || appointment.status === 'confirmed' ? (
                          <Button onClick={() => void handleNoShow(appointment)} size="sm" variant="danger">
                            No-show
                          </Button>
                        ) : null}

                        {nextStatuses.length > 0 ? (
                          <Select
                            aria-label="Mudar status do agendamento"
                            className="h-8 w-auto min-w-28 text-xs font-bold"
                            onChange={(event) => void handleStatusChange(appointment, event.target.value)}
                            value=""
                          >
                            <option value="">Status...</option>
                            {nextStatuses.map((nextStatus) => (
                              <option key={nextStatus} value={nextStatus}>
                                {APPOINTMENT_STATUS_LABELS[nextStatus]}
                              </option>
                            ))}
                          </Select>
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
