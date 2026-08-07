import { ArrowLeft, CalendarDays, Check, Clock, Phone, Trash2, UserRound } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { useNavigate, useParams } from 'react-router-dom'

import {
  fetchAvailability,
  portalCancelAppointment,
  portalLookup,
  portalRescheduleAppointment,
} from '../lib/bookingApi'
import type { PortalAppointment, PortalResult, Slot } from '../lib/bookingApi'
import { ApiError } from '../lib/apiError'
import { formatDateLong, formatMoney, formatTime, nextDateKeys, toDateKey } from '../lib/format'
import { Button } from '../shared/components/Button'
import { Field, Form, Input } from '../shared/components/FormControls'
import { Spinner } from '../shared/components/Spinner'

const STATUS_LABELS: Record<string, string> = {
  scheduled: 'Agendado',
  confirmed: 'Confirmado',
  in_progress: 'Em andamento',
  completed: 'Concluído',
  cancelled: 'Cancelado',
  no_show: 'Não compareceu',
}

const ACTIONABLE_STATUSES = ['scheduled', 'confirmed', 'in_progress']

type PortalStep = 'phone' | 'list' | 'reschedule'

export function PortalPage() {
  const { slug = '' } = useParams<{ slug: string }>()
  const navigate = useNavigate()

  const [step, setStep] = useState<PortalStep>('phone')
  const [phone, setPhone] = useState('')
  const [result, setResult] = useState<PortalResult | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [notice, setNotice] = useState('')
  const [showCancelConfirm, setShowCancelConfirm] = useState(false)
  const [cancellingId, setCancellingId] = useState('')

  const [rescheduling, setRescheduling] = useState<PortalAppointment | null>(null)
  const [dateKeys] = useState(() => nextDateKeys(14))
  const [selectedDateKey, setSelectedDateKey] = useState('')
  const [slots, setSlots] = useState<Slot[]>([])
  const [loadingSlots, setLoadingSlots] = useState(false)
  const [selectedSlot, setSelectedSlot] = useState<Slot | null>(null)
  const [saving, setSaving] = useState(false)

  async function handleLookup(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const trimmed = phone.trim()
    if (!trimmed) {
      setError('Informe seu telefone para consultar os agendamentos.')
      return
    }
    setLoading(true)
    setError('')
    setNotice('')
    try {
      const data = await portalLookup(slug, trimmed)
      setResult(data)
      setStep('list')
    } catch (lookupError) {
      setError((lookupError as Error).message)
    } finally {
      setLoading(false)
    }
  }

  async function refresh() {
    const data = await portalLookup(slug, result?.customer.phone ?? phone)
    setResult(data)
  }

  function startReschedule(appointment: PortalAppointment) {
    setRescheduling(appointment)
    setSelectedDateKey('')
    setSlots([])
    setSelectedSlot(null)
    setError('')
    setNotice('')
    setStep('reschedule')
  }

  useEffect(() => {
    if (step !== 'reschedule' || !rescheduling || !selectedDateKey) {
      return
    }
    let cancelled = false
    setLoadingSlots(true)
    setSlots([])
    setSelectedSlot(null)
    fetchAvailability(slug, selectedDateKey, rescheduling.serviceId, rescheduling.professionalId)
      .then((data) => {
        if (cancelled) {
          return
        }
        if ('professionals' in data) {
          const group = data.professionals.find((item) => item.professionalId === rescheduling.professionalId)
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
  }, [step, rescheduling, selectedDateKey, slug])

  async function handleConfirmReschedule() {
    if (!rescheduling || !selectedSlot) {
      return
    }
    setSaving(true)
    setError('')
    try {
      await portalRescheduleAppointment(slug, result?.customer.phone ?? phone, rescheduling.id, selectedSlot.startsAt)
      await refresh()
      setNotice('Seu horário foi remarcado com sucesso.')
      setRescheduling(null)
      setStep('list')
    } catch (rescheduleError) {
      if (rescheduleError instanceof ApiError && rescheduleError.code === 'time_conflict') {
        setSelectedSlot(null)
        setError('Este horário acabou de ser ocupado. Escolha outro horário disponível.')
        return
      }
      setError((rescheduleError as Error).message)
    } finally {
      setSaving(false)
    }
  }

  async function handleCancel() {
    if (!result || !cancellingId) {
      return
    }
    setSaving(true)
    setError('')
    try {
      await portalCancelAppointment(slug, result.customer.phone, cancellingId, 'Cancelado pelo cliente')
      await refresh()
      setNotice('Seu agendamento foi cancelado. Até a próxima!')
      setShowCancelConfirm(false)
      setCancellingId('')
    } catch (cancelError) {
      setError((cancelError as Error).message)
    } finally {
      setSaving(false)
    }
  }

  const upcoming = result?.appointments.filter((item) => ACTIONABLE_STATUSES.includes(item.status)) ?? []
  const history = result?.appointments.filter((item) => !ACTIONABLE_STATUSES.includes(item.status)) ?? []

  return (
    <main className="flex min-h-dvh justify-center px-4 py-8">
      <div className="w-full max-w-2xl">
        <header className="mb-6 text-center">
          <span className="mx-auto flex h-12 w-12 items-center justify-center rounded-2xl bg-[var(--color-primary)] text-white">
            <UserRound className="h-6 w-6" />
          </span>
          <h1 className="mt-3 font-headline text-2xl font-extrabold text-[var(--color-ink)]">Meu portal</h1>
          <p className="mt-1 text-sm font-semibold text-[var(--color-ink-soft)]">
            Consulte, cancele ou reagende seus horários pelo telefone cadastrado.
          </p>
        </header>

        <section className="rounded-3xl border border-[var(--color-line)] bg-white p-6">
          {error ? (
            <p className="mb-4 rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p>
          ) : null}
          {notice ? (
            <p className="mb-4 rounded-xl border border-emerald-200 bg-emerald-50 p-3 text-sm font-bold text-emerald-700">{notice}</p>
          ) : null}

          {step === 'phone' ? (
            <Form onSubmit={handleLookup}>
              <Field hint="Use o mesmo telefone informado no momento do agendamento." label="Seu telefone">
                <Input
                  autoFocus
                  className="mt-2"
                  inputMode="tel"
                  onChange={(event) => {
                    setPhone(event.target.value)
                    setError('')
                    setNotice('')
                  }}
                  placeholder="(11) 99999-9999"
                  type="tel"
                  value={phone}
                />
              </Field>
              <Button icon={<Phone className="h-4 w-4" />} loading={loading} size="lg" type="submit">
                Consultar meus agendamentos
              </Button>
              <Button icon={<ArrowLeft className="h-4 w-4" />} onClick={() => navigate('/')} variant="ghost">
                Voltar ao início
              </Button>
            </Form>
          ) : step === 'reschedule' && rescheduling ? (
            <ReschedulePanel
              appointment={rescheduling}
              dateKeys={dateKeys}
              loadingSlots={loadingSlots}
              onBack={() => {
                setStep('list')
                setRescheduling(null)
                setError('')
              }}
              onConfirm={() => void handleConfirmReschedule()}
              onSelectDate={(dateKey) => {
                setSelectedDateKey(dateKey)
                setError('')
              }}
              onSelectSlot={(slot) => setSelectedSlot(slot)}
              saving={saving}
              selectedDateKey={selectedDateKey}
              selectedSlot={selectedSlot}
              slots={slots}
            />
          ) : result ? (
            <PortalList
              history={history}
              name={result.customer.name}
              onCancel={(id) => {
                setCancellingId(id)
                setShowCancelConfirm(true)
              }}
              onReschedule={startReschedule}
              onReset={() => {
                setResult(null)
                setStep('phone')
                setPhone('')
                setError('')
                setNotice('')
              }}
              upcoming={upcoming}
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
              <Button disabled={saving} onClick={() => setShowCancelConfirm(false)} variant="secondary">
                Voltar
              </Button>
              <Button loading={saving} onClick={() => void handleCancel()} variant="danger">
                Cancelar agendamento
              </Button>
            </div>
          </div>
        </div>
      ) : null}
    </main>
  )
}

function PortalList({
  name,
  upcoming,
  history,
  onReschedule,
  onCancel,
  onReset,
}: {
  name: string
  upcoming: PortalAppointment[]
  history: PortalAppointment[]
  onReschedule: (appointment: PortalAppointment) => void
  onCancel: (appointmentId: string) => void
  onReset: () => void
}) {
  return (
    <div>
      <div className="flex items-center justify-between gap-3">
        <div>
          <h2 className="font-headline text-lg font-extrabold text-[var(--color-ink)]">Olá, {name.split(' ')[0]}!</h2>
          <p className="mt-1 text-sm font-semibold text-[var(--color-ink-soft)]">Escolha um agendamento para gerenciar.</p>
        </div>
        <Button icon={<ArrowLeft className="h-4 w-4" />} onClick={onReset} variant="ghost">
          Trocar telefone
        </Button>
      </div>

      {upcoming.length === 0 && history.length === 0 ? (
        <div className="mt-6 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-8 text-center">
          <p className="text-sm font-bold text-[var(--color-ink-soft)]">
            Nenhum agendamento encontrado para este telefone.
          </p>
          <p className="mt-1 text-sm font-semibold text-[var(--color-ink-soft)]">
            Verifique se o telefone cadastrado é o mesmo informado no momento do agendamento.
          </p>
        </div>
      ) : null}

      {upcoming.length > 0 ? (
        <>
          <h3 className="mt-6 text-sm font-extrabold uppercase tracking-wide text-[var(--color-ink-soft)]">
            Próximos agendamentos
          </h3>
          <div className="mt-3 grid gap-3">
            {upcoming.map((item) => (
              <AppointmentCard
                key={item.id}
                item={item}
                onCancel={() => onCancel(item.id)}
                onReschedule={() => onReschedule(item)}
              />
            ))}
          </div>
        </>
      ) : null}

      {history.length > 0 ? (
        <>
          <h3 className="mt-8 text-sm font-extrabold uppercase tracking-wide text-[var(--color-ink-soft)]">Histórico</h3>
          <div className="mt-3 grid gap-3">
            {history.map((item) => (
              <AppointmentCard key={item.id} item={item} />
            ))}
          </div>
        </>
      ) : null}
    </div>
  )
}

function AppointmentCard({
  item,
  onReschedule,
  onCancel,
}: {
  item: PortalAppointment
  onReschedule?: () => void
  onCancel?: () => void
}) {
  const canManage = onReschedule && onCancel && item.status !== 'in_progress'
  return (
    <div className="rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div className="min-w-0">
          <p className="font-headline font-extrabold text-[var(--color-ink)]">{item.serviceName || 'Serviço'}</p>
          <p className="mt-1 text-sm font-semibold text-[var(--color-ink-soft)]">
            {item.professionalName || 'Profissional'}
            {item.resourceName ? ` · ${item.resourceName}` : ''}
          </p>
        </div>
        <span
          className={`shrink-0 rounded-full px-3 py-1 text-xs font-bold ${
            item.status === 'cancelled' || item.status === 'no_show'
              ? 'bg-red-100 text-red-700'
              : item.status === 'completed'
                ? 'bg-[var(--color-bg)] text-[var(--color-ink-soft)]'
                : 'bg-emerald-100 text-emerald-700'
          }`}
        >
          {STATUS_LABELS[item.status] ?? item.status}
        </span>
      </div>

      <dl className="mt-3 grid gap-2 text-sm sm:grid-cols-2">
        <div className="flex items-center gap-2">
          <CalendarDays className="h-4 w-4 shrink-0 text-[var(--color-ink-soft)]" />
          <dd className="font-bold text-[var(--color-ink)]">{formatDateLong(item.startsAt.slice(0, 10))}</dd>
        </div>
        <div className="flex items-center gap-2">
          <Clock className="h-4 w-4 shrink-0 text-[var(--color-ink-soft)]" />
          <dd className="font-bold text-[var(--color-ink)]">
            {formatTime(item.startsAt)} às {formatTime(item.endsAt)}
          </dd>
        </div>
      </dl>

      {item.priceCents > 0 ? (
        <p className="mt-2 text-sm font-extrabold text-[var(--color-primary)]">{formatMoney(item.priceCents)}</p>
      ) : null}

      {item.status === 'cancelled' && item.cancellationReason ? (
        <p className="mt-2 rounded-lg bg-red-50 p-2 text-xs font-semibold text-red-700">
          Motivo do cancelamento: {item.cancellationReason}
        </p>
      ) : null}

      {canManage ? (
        <div className="mt-4 grid gap-2 sm:grid-cols-2">
          <Button icon={<CalendarDays className="h-4 w-4" />} onClick={onReschedule} size="sm" variant="secondary">
            Reagendar horário
          </Button>
          <Button icon={<Trash2 className="h-4 w-4" />} onClick={onCancel} size="sm" variant="danger">
            Cancelar agendamento
          </Button>
        </div>
      ) : null}
    </div>
  )
}

function ReschedulePanel({
  appointment,
  dateKeys,
  selectedDateKey,
  slots,
  loadingSlots,
  selectedSlot,
  saving,
  onSelectDate,
  onSelectSlot,
  onConfirm,
  onBack,
}: {
  appointment: PortalAppointment
  dateKeys: string[]
  selectedDateKey: string
  slots: Slot[]
  loadingSlots: boolean
  selectedSlot: Slot | null
  saving: boolean
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
        <span>{appointment.serviceName || 'Serviço'}</span>
        <span className="h-1 w-1 rounded-full bg-[var(--color-ink-soft)]" />
        <span>{appointment.professionalName || 'Profissional'}</span>
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
