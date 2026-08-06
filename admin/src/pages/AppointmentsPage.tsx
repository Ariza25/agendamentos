import { CalendarClock, Plus, Search } from 'lucide-react'
import { useEffect, useState } from 'react'

import { deleteAppointment, fetchAppointments, updateAppointmentStatus } from '../lib/adminApi'
import { APPOINTMENT_STATUS_LABELS, APPOINTMENT_STATUS_TONES, nextStatusesFor } from '../lib/appointmentStatus'
import { formatDate, formatTime, todayKey } from '../lib/format'
import { AppointmentFormModal } from '../features/appointments/components/AppointmentFormModal'
import { useReferenceData } from '../shared/hooks/useReferenceData'
import { Badge } from '../shared/components/Badge'
import { Button } from '../shared/components/Button'
import { ConfirmModal } from '../shared/components/ConfirmModal'
import { EmptyState } from '../shared/components/EmptyState'
import { Input, Select } from '../shared/components/FormControls'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { Appointment, AppointmentStatus } from '../types'

export function AppointmentsPage() {
  const { showToast } = useToast()
  const referenceData = useReferenceData()
  const [date, setDate] = useState(todayKey())
  const [professionalId, setProfessionalId] = useState('')
  const [status, setStatus] = useState('')
  const [appointments, setAppointments] = useState<Appointment[]>([])
  const [loading, setLoading] = useState(true)
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [appointmentToDelete, setAppointmentToDelete] = useState<Appointment | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  function loadAppointments() {
    setLoading(true)
    fetchAppointments({
      from: `${date}T00:00:00`,
      to: `${date}T23:59:59`,
      ...(professionalId ? { professionalId } : {}),
      ...(status ? { status } : {}),
    })
      .then(setAppointments)
      .catch(() => {
        showToast({ message: 'Não foi possível carregar os agendamentos.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  useEffect(loadAppointments, [date, professionalId, status])

  async function handleStatusChange(appointment: Appointment, nextStatus: string) {
    try {
      await updateAppointmentStatus(appointment.id, nextStatus)
      showToast({ message: 'Status atualizado.', tone: 'success' })
      loadAppointments()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    }
  }

  async function handleDelete() {
    if (!appointmentToDelete) {
      return
    }
    setIsDeleting(true)
    try {
      await deleteAppointment(appointmentToDelete.id)
      showToast({ message: 'Agendamento removido.', tone: 'success' })
      setAppointmentToDelete(null)
      loadAppointments()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    } finally {
      setIsDeleting(false)
    }
  }

  return (
    <PageContainer>
      <Section
        actions={
          <Button icon={<Plus className="h-4 w-4" />} onClick={() => setIsFormOpen(true)}>
            Novo agendamento
          </Button>
        }
        subtitle="Filtre os agendamentos do dia e gerencie os horários."
        title="Agendamentos"
      >
        <div className="mt-4 grid gap-3 md:grid-cols-4">
          <FieldInput label="Data">
            <Input type="date" value={date} onChange={(event) => setDate(event.target.value)} />
          </FieldInput>
          <FieldInput label="Profissional">
            <Select value={professionalId} onChange={(event) => setProfessionalId(event.target.value)}>
              <option value="">Todos</option>
              {referenceData.professionals.map((professional) => (
                <option key={professional.id} value={professional.id}>
                  {professional.name}
                </option>
              ))}
            </Select>
          </FieldInput>
          <FieldInput label="Status">
            <Select value={status} onChange={(event) => setStatus(event.target.value)}>
              <option value="">Todos</option>
              {Object.entries(APPOINTMENT_STATUS_LABELS).map(([value, label]) => (
                <option key={value} value={value}>
                  {label}
                </option>
              ))}
            </Select>
          </FieldInput>
          <div className="flex items-end">
            <Button onClick={loadAppointments} variant="secondary">
              <Search className="h-4 w-4" />
              Buscar
            </Button>
          </div>
        </div>
      </Section>

      {loading ? (
        <Section>
          <Spinner />
        </Section>
      ) : appointments.length === 0 ? (
        <Section>
          <EmptyState
            icon={<CalendarClock className="h-10 w-10" />}
            description="Não há agendamentos para este filtro. Crie um novo agendamento para este dia."
            title="Nenhum agendamento encontrado"
          />
        </Section>
      ) : (
        <Section>
          <ul className="mt-2 grid gap-3">
            {appointments.map((appointment) => {
              const customer = referenceData.byId.customer.get(appointment.customerId)
              const service = referenceData.byId.service.get(appointment.serviceId)
              const professional = referenceData.byId.professional.get(appointment.professionalId)
              const nextStatuses = nextStatusesFor(appointment.status)
              return (
                <li
                  className="flex flex-col gap-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] px-4 py-3 md:flex-row md:items-center"
                  key={appointment.id}
                >
                  <div className="min-w-0 flex-1">
                    <div className="flex flex-wrap items-center gap-2">
                      <span className="rounded-xl bg-[var(--color-primary)] px-2.5 py-1 font-headline text-sm font-black text-white">
                        {formatDate(appointment.startsAt)}
                      </span>
                      <span className="font-headline text-sm font-black text-[var(--color-ink)]">
                        {formatTime(appointment.startsAt)} – {formatTime(appointment.endsAt)}
                      </span>
                      <Badge tone={APPOINTMENT_STATUS_TONES[appointment.status]}>{APPOINTMENT_STATUS_LABELS[appointment.status]}</Badge>
                    </div>
                    <p className="mt-1.5 truncate font-bold text-[var(--color-ink)]">
                      {customer?.name ?? 'Cliente desconhecido'}
                      {customer?.phone ? ` · ${customer.phone}` : ''}
                    </p>
                    <p className="truncate text-sm text-[var(--color-ink-soft)]">
                      {service?.name ?? 'Serviço não informado'}
                      {service?.durationMinutes ? ` · ${service.durationMinutes} min` : ''}
                      {professional ? ` · ${professional.name}` : ''}
                    </p>
                    {appointment.notes ? <p className="mt-1 truncate text-xs italic text-[var(--color-ink-soft)]">{appointment.notes}</p> : null}
                  </div>
                  <div className="flex shrink-0 flex-wrap items-center gap-2">
                    {nextStatuses.length > 0 ? (
                      <Select
                        className="h-9 w-auto min-w-32 text-xs"
                        onChange={(event) => void handleStatusChange(appointment, event.target.value)}
                        value=""
                      >
                        <option value="">Mudar status...</option>
                        {nextStatuses.map((nextStatus) => (
                          <option key={nextStatus} value={nextStatus}>
                            {APPOINTMENT_STATUS_LABELS[nextStatus]}
                          </option>
                        ))}
                      </Select>
                    ) : null}
                    <Button onClick={() => setAppointmentToDelete(appointment)} size="sm" variant="danger">
                      Remover
                    </Button>
                  </div>
                </li>
              )
            })}
          </ul>
        </Section>
      )}

      <AppointmentFormModal
        customers={referenceData.customers}
        defaultDate={`${date}T09:00`}
        onClose={() => setIsFormOpen(false)}
        onSaved={loadAppointments}
        open={isFormOpen}
        professionals={referenceData.professionals}
        services={referenceData.services}
      />

      <ConfirmModal
        confirmLabel="Remover"
        loading={isDeleting}
        message="Esta ação não pode ser desfeita. O horário ficará livre novamente."
        onCancel={() => setAppointmentToDelete(null)}
        onConfirm={() => void handleDelete()}
        open={Boolean(appointmentToDelete)}
        title="Remover agendamento?"
      />
    </PageContainer>
  )
}

function FieldInput({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <label className="block">
      <span className="text-xs font-bold uppercase tracking-wide text-[var(--color-ink-soft)]">{label}</span>
      <div className="mt-1">{children}</div>
    </label>
  )
}
