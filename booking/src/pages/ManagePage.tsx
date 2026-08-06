import { ArrowLeft, CalendarDays, Check, Clock, KeyRound, Trash2 } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent, ReactNode } from 'react'
import { useNavigate, useParams } from 'react-router-dom'

import {
  cancelPublicAppointment,
  fetchAvailability,
  fetchPublicAppointment,
  fetchPublicProfessionals,
  fetchPublicServices,
  reschedulePublicAppointment,
} from '../lib/bookingApi'
import type { Professional, PublicAppointment, Service, Slot } from '../lib/bookingApi'
import { ApiError } from '../lib/apiError'
import { formatDateLong, formatTime, nextDateKeys, toDateKey } from '../lib/format'
import { Button } from '../shared/components/Button'
import { Field, Form, Input } from '../shared/components/FormControls'
import { Spinner } from '../shared/components/Spinner'

type ManageStep = 'code' | 'details' | 'reschedule'

export function ManagePage() {
  const { slug = '' } = useParams<{ slug: string }>()
  const navigate = useNavigate()

  const [step, setStep] = useState<ManageStep>('code')
  const [code, setCode] = useState('')
  const [managementCode, setManagementCode] = useState('')
  const [appointment, setAppointment] = useState<PublicAppointment | null>(null)
  const [service, setService] = useState<Service | null>(null)
  const [professional, setProfessional] = useState<Professional | null>(null)
  const [services, setServices] = useState<Service[]>([])
  const [professionals, setProfessionals] = useState<Professional[]>([])
  const [loading, setLoading] = useState(false)
  const [loadingData, setLoadingData] = useState(true)
  const [error, setError] = useState('')
  const [notice, setNotice] = useState('')
  const [showCancelConfirm, setShowCancelConfirm] = useState(false)
  const [isCancelling, setIsCancelling] = useState(false)

  const [dateKeys] = useState(() => nextDateKeys(14))
  const [selectedDateKey, setSelectedDateKey] = useState('')
  const [slots, setSlots] = useState<Slot[]>([])
  const [loadingSlots, setLoadingSlots] = useState(false)
  const [selectedSlot, setSelectedSlot] = useState<Slot | null>(null)
  const [isRescheduling, setIsRescheduling] = useState(false)

  useEffect(() => {
    let active = true
    setLoadingData(true)
    Promise.all([fetchPublicServices(slug), fetchPublicProfessionals(slug)])
      .then(([servicesData, professionalsData]) => {
        if (!active) {
          return
        }
        setServices(servicesData)
        setProfessionals(professionalsData)
      })
      .catch(() => {
        if (active) {
          setError('Não foi possível carregar os dados do estabelecimento.')
        }
      })
      .finally(() => {
        if (active) {
          setLoadingData(false)
        }
      })
    return () => {
      active = false
    }
  }, [slug])

  async function handleLookup(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const trimmed = code.trim().toUpperCase()
    if (!trimmed) {
      setError('Informe o código do agendamento.')
      return
    }

    setLoading(true)
    setError('')
    setNotice('')
    try {
      const appointmentData = await fetchPublicAppointment(slug, trimmed)
      if (appointmentData.status === 'cancelled' || appointmentData.status === 'no_show' || appointmentData.status === 'completed') {
        setNotice('Este agendamento já foi encerrado e não pode mais ser gerenciado online.')
        return
      }
      setManagementCode(trimmed)
      setAppointment(appointmentData)
      setService(services.find((item) => item.id === appointmentData.serviceId) ?? null)
      setProfessional(professionals.find((item) => item.id === appointmentData.professionalId) ?? null)
      setStep('details')
    } catch (lookupError) {
      if (lookupError instanceof ApiError && lookupError.code === 'appointment_done') {
        setNotice((lookupError as Error).message)
      } else {
        setError((lookupError as Error).message)
      }
    } finally {
      setLoading(false)
    }
  }

  async function handleCancel() {
    if (!appointment) {
      return
    }
    setIsCancelling(true)
    setError('')
    try {
      await cancelPublicAppointment(slug, managementCode, 'Cancelado pelo cliente')
      setNotice('Seu agendamento foi cancelado. Até a próxima!')
      setStep('code')
      setAppointment(null)
      setManagementCode('')
      setService(null)
      setProfessional(null)
      setShowCancelConfirm(false)
    } catch (cancelError) {
      setError((cancelError as Error).message)
    } finally {
      setIsCancelling(false)
    }
  }

  function startReschedule() {
    setSelectedDateKey('')
    setSlots([])
    setSelectedSlot(null)
    setStep('reschedule')
    setError('')
  }

  useEffect(() => {
    if (step !== 'reschedule' || !appointment || !selectedDateKey) {
      return
    }
    let cancelled = false
    setLoadingSlots(true)
    setSlots([])
    setSelectedSlot(null)
    fetchAvailability(slug, selectedDateKey, appointment.serviceId, appointment.professionalId)
      .then((data) => {
        if (cancelled) {
          return
        }
        if ('professionals' in data) {
          const group = data.professionals.find((item) => item.professionalId === appointment.professionalId)
          setSlots(group?.slots ?? [])
        } else {
          setSlots(data.slots)
        }
      })
      .catch((slotError: unknown) => {
        if (!cancelled) {
          setError((slotError as Error).message)
        }
      })
      .finally(() => {
        if (!cancelled) {
          setLoadingSlots(false)
        }
      })
    return () => {
      cancelled = true
    }
  }, [step, appointment, selectedDateKey, slug])

  async function handleReschedule() {
    if (!appointment || !selectedSlot) {
      return
    }
    setIsRescheduling(true)
    setError('')
    try {
      const updated = await reschedulePublicAppointment(slug, managementCode, selectedSlot.startsAt)
      setAppointment(updated)
      setNotice('Seu horário foi remarcado com sucesso.')
      setStep('details')
    } catch (rescheduleError) {
      if (rescheduleError instanceof ApiError && rescheduleError.code === 'time_conflict') {
        setSelectedSlot(null)
        setError('Este horário acabou de ser ocupado. Escolha outro horário disponível.')
        return
      }
      setError((rescheduleError as Error).message)
    } finally {
      setIsRescheduling(false)
    }
  }

  return (
    <main className="flex min-h-dvh justify-center px-4 py-8">
      <div className="w-full max-w-2xl">
        <header className="mb-6 text-center">
          <span className="mx-auto flex h-12 w-12 items-center justify-center rounded-2xl bg-[var(--color-primary)] text-white">
            <KeyRound className="h-6 w-6" />
          </span>
          <h1 className="mt-3 font-headline text-2xl font-extrabold text-[var(--color-ink)]">Meu agendamento</h1>
          <p className="mt-1 text-sm font-semibold text-[var(--color-ink-soft)]">
            Informe o código recebido na confirmação para cancelar ou reagendar.
          </p>
        </header>

        <section className="rounded-3xl border border-[var(--color-line)] bg-white p-6">
          {error ? (
            <p className="mb-4 rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p>
          ) : null}
          {notice ? (
            <p className="mb-4 rounded-xl border border-emerald-200 bg-emerald-50 p-3 text-sm font-bold text-emerald-700">{notice}</p>
          ) : null}

          {loadingData ? (
            <Spinner label="Carregando estabelecimento..." />
          ) : step === 'code' ? (
            <Form onSubmit={handleLookup}>
              <Field hint="O código está na tela de confirmação do seu agendamento." label="Código do agendamento">
                <Input
                  autoFocus
                  className="mt-2 font-headline tracking-widest"
                  onChange={(event) => {
                    setCode(event.target.value)
                    setError('')
                    setNotice('')
                  }}
                  placeholder="Ex.: AB-CDE-12"
                  value={code}
                />
              </Field>
              <Button icon={<KeyRound className="h-4 w-4" />} loading={loading} size="lg" type="submit">
                Consultar agendamento
              </Button>
              <Button icon={<ArrowLeft className="h-4 w-4" />} onClick={() => navigate('/')} variant="ghost">
                Voltar ao início
              </Button>
            </Form>
          ) : step === 'reschedule' && appointment ? (
            <ReschedulePicker
              dateKeys={dateKeys}
              loadingSlots={loadingSlots}
              onBack={() => {
                setStep('details')
                setError('')
              }}
              onConfirm={() => void handleReschedule()}
              onSelectDate={(dateKey) => {
                setSelectedDateKey(dateKey)
                setError('')
              }}
              onSelectSlot={(slot) => setSelectedSlot(slot)}
              professionalName={professional?.name ?? 'Profissional'}
              saving={isRescheduling}
              selectedDateKey={selectedDateKey}
              selectedSlot={selectedSlot}
              serviceName={service?.name ?? 'Serviço'}
              slots={slots}
            />
          ) : appointment ? (
            <AppointmentDetails
              appointment={appointment}
              onCancel={() => setShowCancelConfirm(true)}
              onReschedule={startReschedule}
              professionalName={professional?.name ?? 'Profissional'}
              serviceName={service?.name ?? 'Serviço'}
            />
          ) : null}
        </section>
      </div>

      {showCancelConfirm ? (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4">
          <div className="w-full max-w-md rounded-3xl border border-[var(--color-line)] bg-white p-6 shadow-2xl">
            <h2 className="font-headline text-xl font-extrabold text-[var(--color-ink)]">Cancelar agendamento?</h2>
            <p className="mt-2 text-sm font-semibold text-[var(--color-ink-soft)]">
              Esta ação não pode ser desfeita. O horário voltará a ficar disponível para outras pessoas.
            </p>
            <div className="mt-6 grid gap-3 sm:grid-cols-2">
              <Button disabled={isCancelling} onClick={() => setShowCancelConfirm(false)} variant="secondary">
                Voltar
              </Button>
              <Button loading={isCancelling} onClick={() => void handleCancel()} variant="danger">
                Cancelar agendamento
              </Button>
            </div>
          </div>
        </div>
      ) : null}
    </main>
  )
}

function AppointmentDetails({
  appointment,
  serviceName,
  professionalName,
  onCancel,
  onReschedule,
}: {
  appointment: PublicAppointment
  serviceName: string
  professionalName: string
  onCancel: () => void
  onReschedule: () => void
}) {
  return (
    <div>
      <h2 className="font-headline text-lg font-extrabold text-[var(--color-ink)]">Dados do agendamento</h2>

      <dl className="mt-4 space-y-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4 text-sm">
        <div className="flex items-center justify-between gap-3">
          <dt className="font-semibold text-[var(--color-ink-soft)]">Serviço</dt>
          <dd className="truncate font-bold text-[var(--color-ink)]">{serviceName}</dd>
        </div>
        <div className="flex items-center justify-between gap-3">
          <dt className="font-semibold text-[var(--color-ink-soft)]">Profissional</dt>
          <dd className="truncate font-bold text-[var(--color-ink)]">{professionalName}</dd>
        </div>
        <div className="flex items-center justify-between gap-3">
          <dt className="font-semibold text-[var(--color-ink-soft)]">Data</dt>
          <dd className="font-bold text-[var(--color-ink)]">{formatDateLong(appointment.startsAt.slice(0, 10))}</dd>
        </div>
        <div className="flex items-center justify-between gap-3">
          <dt className="font-semibold text-[var(--color-ink-soft)]">Horário</dt>
          <dd className="font-bold text-[var(--color-ink)]">
            {formatTime(appointment.startsAt)} às {formatTime(appointment.endsAt)}
          </dd>
        </div>
        <div className="flex items-center justify-between gap-3">
          <dt className="font-semibold text-[var(--color-ink-soft)]">Status</dt>
          <dd className="font-bold text-[var(--color-success)]">Confirmado</dd>
        </div>
      </dl>

      <div className="mt-6 grid gap-3 sm:grid-cols-2">
        <Button icon={<CalendarDays className="h-4 w-4" />} onClick={onReschedule} size="lg">
          Reagendar horário
        </Button>
        <Button icon={<Trash2 className="h-4 w-4" />} onClick={onCancel} size="lg" variant="danger">
          Cancelar agendamento
        </Button>
      </div>
    </div>
  )
}

function ReschedulePicker({
  dateKeys,
  selectedDateKey,
  slots,
  loadingSlots,
  selectedSlot,
  saving,
  serviceName,
  professionalName,
  onSelectDate,
  onSelectSlot,
  onConfirm,
  onBack,
}: {
  dateKeys: string[]
  selectedDateKey: string
  slots: Slot[]
  loadingSlots: boolean
  selectedSlot: Slot | null
  saving: boolean
  serviceName: string
  professionalName: string
  onSelectDate: (dateKey: string) => void
  onSelectSlot: (slot: Slot) => void
  onConfirm: () => void
  onBack: () => void
}) {
  const todayKey = toDateKey(new Date())
  return (
    <div>
      <div className="flex items-center justify-between gap-3">
        <h2 className="font-headline text-lg font-extrabold text-[var(--color-ink)]">Escolha o novo horário</h2>
        <Button icon={<ArrowLeft className="h-4 w-4" />} onClick={onBack} variant="ghost">
          Voltar
        </Button>
      </div>
      <p className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-sm font-semibold text-[var(--color-ink-soft)]">
        <span>{serviceName}</span>
        <span className="h-1 w-1 rounded-full bg-[var(--color-ink-soft)]" />
        <span>{professionalName}</span>
      </p>

      <p className="mt-5 text-sm font-bold text-[var(--color-ink)]">Dia</p>
      <div className="mt-2 flex gap-2 overflow-x-auto pb-2">
        {dateKeys.map((dateKey) => {
          const date = new Date(`${dateKey}T00:00:00`)
          const isToday = dateKey === todayKey
          const isSelected = dateKey === selectedDateKey
          return (
            <button
              className={`flex min-w-16 shrink-0 flex-col items-center gap-1 rounded-2xl border px-3 py-3 transition ${
                isSelected
                  ? 'border-[var(--color-primary)] bg-[var(--color-primary)] text-white'
                  : 'border-[var(--color-line)] bg-[var(--color-bg-muted)] text-[var(--color-ink)] hover:border-[var(--color-primary)]'
              }`}
              key={dateKey}
              onClick={() => onSelectDate(dateKey)}
              type="button"
            >
              <span className="text-xs font-bold opacity-80">{date.toLocaleDateString('pt-BR', { weekday: 'short' })}</span>
              <span className="font-headline text-lg font-extrabold leading-none">{date.getDate()}</span>
              <span className="text-xs font-semibold opacity-80">
                {isToday ? 'Hoje' : date.toLocaleDateString('pt-BR', { month: 'short' })}
              </span>
            </button>
          )
        })}
      </div>

      {selectedDateKey ? (
        <>
          <p className="mt-5 text-sm font-bold text-[var(--color-ink)]">
            Horários disponíveis em {formatDateLong(selectedDateKey)}
          </p>
          <div className="mt-2">
            {loadingSlots ? (
              <Spinner label="Buscando horários..." />
            ) : slots.length === 0 ? (
              <p className="rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-6 text-center text-sm font-bold text-[var(--color-ink-soft)]">
                Não há horários disponíveis nesta data. Escolha outro dia.
              </p>
            ) : (
              <ul className="grid grid-cols-2 gap-2 sm:grid-cols-3 md:grid-cols-4">
                {slots.map((slot) => (
                  <li key={slot.startsAt}>
                    <button
                      className={`flex w-full items-center justify-center gap-2 rounded-xl border px-3 py-3 text-sm font-bold transition ${
                        selectedSlot?.startsAt === slot.startsAt
                          ? 'border-[var(--color-primary)] bg-[var(--color-primary)] text-white'
                          : 'border-[var(--color-line)] bg-white text-[var(--color-ink)] hover:border-[var(--color-primary)]'
                      }`}
                      onClick={() => onSelectSlot(slot)}
                      type="button"
                    >
                      <Clock className="h-4 w-4" />
                      {formatTime(slot.startsAt)}
                    </button>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </>
      ) : null}

      <Button className="mt-6" disabled={!selectedSlot} loading={saving} onClick={onConfirm} size="lg">
        <Check className="h-4 w-4" />
        Confirmar novo horário
      </Button>
    </div>
  )
}
