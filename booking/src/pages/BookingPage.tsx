import { ArrowLeft, CalendarDays, Check, Clock, Copy, User, Users } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent, ReactNode } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'

import {
  createPublicBooking,
  fetchAvailability,
  fetchPublicInfo,
  fetchPublicProfessionals,
  fetchPublicServices,
} from '../lib/bookingApi'
import type { BookingResult, Professional, PublicInfo, Service, Slot } from '../lib/bookingApi'
import { ApiError } from '../lib/apiError'
import { formatDateLong, formatMoney, formatTime, nextDateKeys, toDateKey } from '../lib/format'
import { Button } from '../shared/components/Button'
import { Field, Form, Input, Textarea } from '../shared/components/FormControls'
import { Spinner } from '../shared/components/Spinner'

type Step = 'service' | 'professional' | 'datetime' | 'details' | 'done'

const STEP_LABELS: Record<Step, string> = {
  service: 'Serviço',
  professional: 'Profissional',
  datetime: 'Data e horário',
  details: 'Seus dados',
  done: 'Confirmado',
}

const STEP_ORDER: Step[] = ['service', 'professional', 'datetime', 'details', 'done']

export function BookingPage() {
  const { slug = '' } = useParams<{ slug: string }>()
  const navigate = useNavigate()

  const [info, setInfo] = useState<PublicInfo | null>(null)
  const [services, setServices] = useState<Service[]>([])
  const [professionals, setProfessionals] = useState<Professional[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [step, setStep] = useState<Step>('service')
  const [service, setService] = useState<Service | null>(null)
  const [professional, setProfessional] = useState<Professional | null>(null)
  const [dateKeys] = useState(() => nextDateKeys(14))
  const [selectedDateKey, setSelectedDateKey] = useState('')
  const [slots, setSlots] = useState<Slot[]>([])
  const [loadingSlots, setLoadingSlots] = useState(false)
  const [selectedSlot, setSelectedSlot] = useState<Slot | null>(null)
  const [booking, setBooking] = useState<BookingResult | null>(null)

  function loadAll() {
    setLoading(true)
    setError('')
    Promise.all([
      fetchPublicInfo(slug),
      fetchPublicServices(slug),
      fetchPublicProfessionals(slug),
    ])
      .then(([infoData, servicesData, professionalsData]) => {
        setInfo(infoData)
        setServices(servicesData)
        setProfessionals(professionalsData)
      })
      .catch((loadError: unknown) => {
        setError((loadError as Error).message)
      })
      .finally(() => setLoading(false))
  }

  useEffect(loadAll, [slug])

  useEffect(() => {
    if (!service || !professional || !selectedDateKey) {
      return
    }
    let cancelled = false
    setLoadingSlots(true)
    setSlots([])
    setSelectedSlot(null)
    fetchAvailability(slug, selectedDateKey, service.id, professional.id)
      .then((data) => {
        if (cancelled) {
          return
        }
        if ('professionals' in data) {
          const group = data.professionals.find((item) => item.professionalId === professional.id)
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
  }, [slug, service, professional, selectedDateKey])

  const stepIndex = STEP_ORDER.indexOf(step)

  function goBack() {
    if (stepIndex <= 0) {
      return
    }
    setStep(STEP_ORDER[stepIndex - 1])
    setError('')
  }

  function selectService(nextService: Service) {
    setService(nextService)
    setProfessional(null)
    setSelectedDateKey('')
    setSlots([])
    setSelectedSlot(null)
    setStep('professional')
    setError('')
  }

  function selectProfessional(nextProfessional: Professional) {
    setProfessional(nextProfessional)
    setSelectedDateKey('')
    setSlots([])
    setSelectedSlot(null)
    setStep('datetime')
    setError('')
  }

  function selectDate(dateKey: string) {
    setSelectedDateKey(dateKey)
    setSelectedSlot(null)
    setError('')
  }

  function selectSlot(slot: Slot) {
    setSelectedSlot(slot)
    setStep('details')
    setError('')
  }

  async function handleBookingSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!service || !professional || !selectedSlot) {
      return
    }
    const formData = new FormData(event.currentTarget)
    const name = String(formData.get('name') ?? '').trim()
    const phone = String(formData.get('phone') ?? '').trim()
    const email = String(formData.get('email') ?? '').trim()
    const notes = String(formData.get('notes') ?? '').trim()

    if (!name) {
      setError('Informe seu nome.')
      return
    }
    if (!phone) {
      setError('Informe seu telefone para contato.')
      return
    }

    setError('')
    try {
      const result = await createPublicBooking(slug, {
        serviceId: service.id,
        professionalId: professional.id,
        startsAt: selectedSlot.startsAt,
        name,
        phone,
        email,
        notes,
      })
      setBooking(result)
      setStep('done')
    } catch (submitError) {
      if (submitError instanceof ApiError && submitError.code === 'time_conflict') {
        setStep('datetime')
        setSelectedSlot(null)
        setError('Este horário acabou de ser ocupado. Escolha outro horário disponível.')
        return
      }
      setError((submitError as Error).message)
    }
  }

  if (loading) {
    return (
      <CenteredLayout>
        <Spinner label="Carregando estabelecimento..." />
      </CenteredLayout>
    )
  }

  if (error && !info) {
    return (
      <CenteredLayout>
        <ErrorCard message={error} onRetry={loadAll} onBack={() => navigate('/')} />
      </CenteredLayout>
    )
  }

  if (info && !info.configuration.allowOnlineBooking) {
    return (
      <CenteredLayout>
        <ErrorCard
          message="No momento este estabelecimento não está aceitando agendamentos online. Tente novamente mais tarde."
          onBack={() => navigate('/')}
        />
      </CenteredLayout>
    )
  }

  if (!info) {
    return (
      <CenteredLayout>
        <ErrorCard message="Não foi possível carregar o estabelecimento." onRetry={loadAll} onBack={() => navigate('/')} />
      </CenteredLayout>
    )
  }

  if (step === 'done' && booking) {
    return <Confirmation info={info} booking={booking} service={service} professional={professional} />
  }

  return (
    <main className="flex min-h-dvh justify-center px-4 py-8">
      <div className="w-full max-w-2xl">
        <header className="mb-6 text-center">
          <span className="mx-auto flex h-12 w-12 items-center justify-center rounded-2xl bg-[var(--color-primary)] text-white">
            <CalendarDays className="h-7 w-7" />
          </span>
          <h1 className="mt-3 font-headline text-2xl font-extrabold text-[var(--color-ink)]">{info.tenant.name}</h1>
          <p className="mt-1 text-sm font-semibold text-[var(--color-ink-soft)]">
            Agende seu horário em poucos passos.
          </p>
        </header>

        <StepProgress labels={STEP_LABELS} steps={STEP_ORDER} currentStep={step} />

        <section className="mt-6 rounded-3xl border border-[var(--color-line)] bg-white p-6">
          {error ? (
            <p className="mb-4 rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p>
          ) : null}

          {stepIndex > 0 ? (
            <Button className="mb-4" icon={<ArrowLeft className="h-4 w-4" />} onClick={goBack} variant="ghost">
              Voltar
            </Button>
          ) : null}

          {step === 'service' ? (
            <ServiceStep onSelect={selectService} services={services} />
          ) : null}

          {step === 'professional' ? (
            <ProfessionalStep onSelect={selectProfessional} professionals={professionals} />
          ) : null}

          {step === 'datetime' && service && professional ? (
            <DateTimeStep
              dateKeys={dateKeys}
              loadingSlots={loadingSlots}
              onSelectDate={selectDate}
              onSelectSlot={selectSlot}
              professionalName={professional.name}
              selectedDateKey={selectedDateKey}
              selectedSlot={selectedSlot}
              serviceName={service.name}
              slots={slots}
            />
          ) : null}

          {step === 'details' && service && professional && selectedSlot ? (
            <DetailsStep
              onSubmit={handleBookingSubmit}
              professionalName={professional.name}
              serviceName={service.name}
              slot={selectedSlot}
            />
          ) : null}
        </section>
      </div>
    </main>
  )
}

function StepProgress({
  labels,
  steps,
  currentStep,
}: {
  labels: Record<Step, string>
  steps: Step[]
  currentStep: Step
}) {
  const currentIndex = steps.indexOf(currentStep)
  return (
    <ol className="flex items-center justify-center gap-2">
      {steps.slice(0, 4).map((stepName, index) => (
        <li className="flex items-center gap-2" key={stepName}>
          <span
            className={`flex items-center gap-2 rounded-full px-3 py-1.5 text-xs font-bold ${
              index === currentIndex
                ? 'bg-[var(--color-primary)] text-white'
                : index < currentIndex
                  ? 'bg-[var(--color-bg-muted)] text-[var(--color-primary)]'
                  : 'bg-[var(--color-bg-muted)] text-[var(--color-ink-soft)]'
            }`}
          >
            {index < currentIndex ? <Check className="h-3.5 w-3.5" /> : <span>{index + 1}</span>}
            <span className={index === currentIndex ? '' : 'hidden sm:inline'}>{labels[stepName]}</span>
          </span>
          {index < steps.slice(0, 4).length - 1 ? (
            <span className="h-px w-4 bg-[var(--color-line)] sm:w-8" />
          ) : null}
        </li>
      ))}
    </ol>
  )
}

function ServiceStep({ services, onSelect }: { services: Service[]; onSelect: (service: Service) => void }) {
  if (services.length === 0) {
    return <EmptyNote message="Este estabelecimento ainda não possui serviços disponíveis para agendamento." />
  }
  return (
    <div>
      <h2 className="font-headline text-lg font-extrabold text-[var(--color-ink)]">Escolha o serviço</h2>
      <ul className="mt-4 grid gap-3">
        {services.map((item) => (
          <li key={item.id}>
            <button
              className="flex w-full items-center gap-4 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4 text-left transition hover:border-[var(--color-primary)] hover:bg-[var(--color-bg-elevated)]"
              onClick={() => onSelect(item)}
              type="button"
            >
              <span className="h-4 w-4 shrink-0 rounded-full" style={{ backgroundColor: item.color || '#b42318' }} />
              <span className="min-w-0 flex-1">
                <span className="block truncate font-bold text-[var(--color-ink)]">{item.name}</span>
                {item.description ? (
                  <span className="mt-0.5 line-clamp-2 block text-sm text-[var(--color-ink-soft)]">{item.description}</span>
                ) : null}
              </span>
              <span className="flex shrink-0 flex-col items-end gap-1 text-sm font-semibold text-[var(--color-ink-soft)]">
                <span>{item.durationMinutes} min</span>
                <span className="font-bold text-[var(--color-ink)]">{formatMoney(item.priceCents)}</span>
              </span>
            </button>
          </li>
        ))}
      </ul>
    </div>
  )
}

function ProfessionalStep({
  professionals,
  onSelect,
}: {
  professionals: Professional[]
  onSelect: (professional: Professional) => void
}) {
  if (professionals.length === 0) {
    return <EmptyNote message="Este estabelecimento ainda não possui profissionais disponíveis para agendamento." />
  }
  return (
    <div>
      <h2 className="font-headline text-lg font-extrabold text-[var(--color-ink)]">Escolha o profissional</h2>
      <ul className="mt-4 grid gap-3 sm:grid-cols-2">
        {professionals.map((item) => (
          <li key={item.id}>
            <button
              className="flex w-full items-center gap-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4 text-left transition hover:border-[var(--color-primary)] hover:bg-[var(--color-bg-elevated)]"
              onClick={() => onSelect(item)}
              type="button"
            >
              <span
                className="flex h-11 w-11 shrink-0 items-center justify-center rounded-full text-white"
                style={{ backgroundColor: item.color || '#b42318' }}
              >
                <User className="h-5 w-5" />
              </span>
              <span className="min-w-0">
                <span className="block truncate font-bold text-[var(--color-ink)]">{item.name}</span>
                {item.email ? <span className="mt-0.5 block truncate text-sm text-[var(--color-ink-soft)]">{item.email}</span> : null}
              </span>
            </button>
          </li>
        ))}
      </ul>
    </div>
  )
}

function DateTimeStep({
  dateKeys,
  selectedDateKey,
  slots,
  loadingSlots,
  selectedSlot,
  serviceName,
  professionalName,
  onSelectDate,
  onSelectSlot,
}: {
  dateKeys: string[]
  selectedDateKey: string
  slots: Slot[]
  loadingSlots: boolean
  selectedSlot: Slot | null
  serviceName: string
  professionalName: string
  onSelectDate: (dateKey: string) => void
  onSelectSlot: (slot: Slot) => void
}) {
  const todayKey = toDateKey(new Date())
  return (
    <div>
      <h2 className="font-headline text-lg font-extrabold text-[var(--color-ink)]">Escolha a data e o horário</h2>
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
              <EmptyNote message="Não há horários disponíveis nesta data. Escolha outro dia." />
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
    </div>
  )
}

function DetailsStep({
  serviceName,
  professionalName,
  slot,
  onSubmit,
}: {
  serviceName: string
  professionalName: string
  slot: Slot
  onSubmit: (event: FormEvent<HTMLFormElement>) => void
}) {
  const [name, setName] = useState('')
  const [phone, setPhone] = useState('')
  return (
    <div>
      <h2 className="font-headline text-lg font-extrabold text-[var(--color-ink)]">Confirme seus dados</h2>

      <div className="mt-4 flex flex-wrap items-center gap-x-4 gap-y-1 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4 text-sm">
        <span className="inline-flex items-center gap-2 font-bold text-[var(--color-ink)]">
          <CalendarDays className="h-4 w-4" />
          {formatDateLong(slot.startsAt.slice(0, 10))}
        </span>
        <span className="inline-flex items-center gap-2 font-bold text-[var(--color-ink)]">
          <Clock className="h-4 w-4" />
          {formatTime(slot.startsAt)} às {formatTime(slot.endsAt)}
        </span>
      </div>

      <p className="mt-4 flex flex-wrap items-center gap-x-3 gap-y-1 text-sm font-semibold text-[var(--color-ink-soft)]">
        <span>{serviceName}</span>
        <span className="h-1 w-1 rounded-full bg-[var(--color-ink-soft)]" />
        <span className="inline-flex items-center gap-1.5">
          <Users className="h-4 w-4" />
          {professionalName}
        </span>
      </p>

      <Form onSubmit={onSubmit}>
        <Field label="Nome completo">
          <Input
            autoFocus
            className="mt-2"
            name="name"
            onChange={(event) => setName(event.target.value)}
            placeholder="Como podemos te chamar?"
            value={name}
          />
        </Field>
        <Field label="Telefone / WhatsApp">
          <Input
            className="mt-2"
            name="phone"
            onChange={(event) => setPhone(event.target.value)}
            placeholder="(00) 00000-0000"
            value={phone}
          />
        </Field>
        <Field hint="Opcional" label="E-mail">
          <Input className="mt-2" name="email" placeholder="voce@email.com" type="email" />
        </Field>
        <Field hint="Opcional" label="Observações">
          <Textarea className="mt-2" name="notes" placeholder="Conte-nos algo que gostaria de acrescentar..." />
        </Field>
        <Button size="lg" type="submit">
          Confirmar agendamento
        </Button>
      </Form>
    </div>
  )
}

function Confirmation({
  info,
  booking,
  service,
  professional,
}: {
  info: PublicInfo
  booking: BookingResult
  service: Service | null
  professional: Professional | null
}) {
  const navigate = useNavigate()
  return (
    <main className="flex min-h-dvh items-center justify-center px-4 py-10">
      <section className="w-full max-w-md rounded-3xl border border-[var(--color-line)] bg-white p-8 text-center shadow-sm">
        <span className="mx-auto flex h-14 w-14 items-center justify-center rounded-full bg-[var(--color-success)] text-white">
          <Check className="h-7 w-7" />
        </span>
        <h1 className="mt-4 font-headline text-2xl font-extrabold text-[var(--color-ink)]">Horário reservado!</h1>
        <p className="mt-1 text-sm font-semibold text-[var(--color-ink-soft)]">
          {info.tenant.name} já recebeu sua solicitação. Qualquer dúvida, entre em contato por telefone.
        </p>

        {booking.managementCode ? (
          <div className="mt-6 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4">
            <p className="text-xs font-bold uppercase tracking-wide text-[var(--color-ink-soft)]">
              Código do agendamento
            </p>
            <p className="mt-1 flex items-center justify-center gap-2">
              <span className="font-headline text-3xl font-black tracking-[0.2em] text-[var(--color-ink)]">
                {booking.managementCode}
              </span>
              <button
                aria-label="Copiar código"
                className="rounded-lg p-1.5 text-[var(--color-ink-soft)] transition hover:bg-[var(--color-bg-elevated)] hover:text-[var(--color-primary)]"
                onClick={() => {
                  void navigator.clipboard?.writeText(booking.managementCode)
                }}
                type="button"
              >
                <Copy className="h-4 w-4" />
              </button>
            </p>
            <p className="mt-2 text-xs font-semibold text-[var(--color-ink-soft)]">
              Guarde este código para cancelar ou reagendar pelo site.
            </p>
          </div>
        ) : null}

        <dl className="mt-6 space-y-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4 text-left text-sm">
          <div className="flex items-center justify-between gap-3">
            <dt className="font-semibold text-[var(--color-ink-soft)]">Serviço</dt>
            <dd className="truncate font-bold text-[var(--color-ink)]">{service?.name ?? '—'}</dd>
          </div>
          <div className="flex items-center justify-between gap-3">
            <dt className="font-semibold text-[var(--color-ink-soft)]">Profissional</dt>
            <dd className="truncate font-bold text-[var(--color-ink)]">{professional?.name ?? '—'}</dd>
          </div>
          <div className="flex items-center justify-between gap-3">
            <dt className="font-semibold text-[var(--color-ink-soft)]">Data</dt>
            <dd className="font-bold text-[var(--color-ink)]">{formatDateLong(booking.appointment.startsAt.slice(0, 10))}</dd>
          </div>
          <div className="flex items-center justify-between gap-3">
            <dt className="font-semibold text-[var(--color-ink-soft)]">Horário</dt>
            <dd className="font-bold text-[var(--color-ink)]">
              {formatTime(booking.appointment.startsAt)} às {formatTime(booking.appointment.endsAt)}
            </dd>
          </div>
          <div className="flex items-center justify-between gap-3">
            <dt className="font-semibold text-[var(--color-ink-soft)]">Status</dt>
            <dd className="font-bold text-[var(--color-success)]">Confirmado</dd>
          </div>
        </dl>

        <Button className="mt-6" onClick={() => navigate('/')} size="lg" variant="secondary">
          Fazer outro agendamento
        </Button>
        {booking.managementCode ? (
          <Link
            className="mt-3 block text-sm font-bold text-[var(--color-primary)] underline-offset-4 hover:underline"
            to={`/manage/${info.tenant.slug}`}
          >
            Cancelar ou reagendar este horário
          </Link>
        ) : null}
      </section>
    </main>
  )
}

function EmptyNote({ message }: { message: string }) {
  return (
    <p className="rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-6 text-center text-sm font-bold text-[var(--color-ink-soft)]">
      {message}
    </p>
  )
}

function ErrorCard({
  message,
  onRetry,
  onBack,
}: {
  message: string
  onRetry?: () => void
  onBack: () => void
}) {
  return (
    <section className="w-full max-w-md rounded-3xl border border-[var(--color-line)] bg-white p-8 text-center shadow-sm">
      <h1 className="font-headline text-xl font-extrabold text-[var(--color-ink)]">Não foi possível agendar</h1>
      <p className="mt-2 text-sm font-semibold text-[var(--color-ink-soft)]">{message}</p>
      <div className="mt-6 flex flex-col gap-3">
        {onRetry ? (
          <Button onClick={onRetry}>Tentar novamente</Button>
        ) : null}
        <Button onClick={onBack} variant="secondary">
          Voltar ao início
        </Button>
      </div>
    </section>
  )
}

function CenteredLayout({ children }: { children: ReactNode }) {
  return (
    <main className="flex min-h-dvh items-center justify-center px-4 py-10">
      <div className="w-full max-w-md">{children}</div>
    </main>
  )
}
