import { ArrowLeft, CalendarDays, Check, Clock, Copy, Globe, Instagram, MapPin, MessageCircle, Phone, ShieldCheck, User, Users, X } from 'lucide-react'
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
  const [selectedServices, setSelectedServices] = useState<Service[]>([])
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
    if (info?.configuration?.primaryColor) {
      document.documentElement.style.setProperty('--color-primary', info.configuration.primaryColor)
    }
  }, [info])

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

  function selectServices(items: Service[]) {
    setSelectedServices(items)
    setService(items[0] || null)
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
    const timezone = String(formData.get('timezone') ?? '').trim() || Intl.DateTimeFormat().resolvedOptions().timeZone || 'America/Sao_Paulo'
    const recurrenceOption = String(formData.get('recurrence') ?? 'none')
    const lgpdAccepted = formData.get('lgpdAccepted') === 'on' || formData.get('lgpdAccepted') === 'true'
    const marketingConsent = formData.get('marketingConsent') === 'on' || formData.get('marketingConsent') === 'true'

    let recurrenceRule: string | undefined
    let recurrenceCount: number | undefined
    if (recurrenceOption === 'weekly_2') {
      recurrenceRule = 'weekly'
      recurrenceCount = 2
    } else if (recurrenceOption === 'weekly_4') {
      recurrenceRule = 'weekly'
      recurrenceCount = 4
    } else if (recurrenceOption === 'weekly_8') {
      recurrenceRule = 'weekly'
      recurrenceCount = 8
    }

    if (!name) {
      setError('Informe seu nome.')
      return
    }
    if (!phone) {
      setError('Informe seu telefone para contato.')
      return
    }
    if (!lgpdAccepted) {
      setError('É obrigatório aceitar os Termos de Uso e a Política de Privacidade (LGPD) para prosseguir.')
      return
    }

    setError('')
    try {
      const result = await createPublicBooking(slug, {
        serviceId: service.id,
        serviceIds: selectedServices.length > 1 ? selectedServices.map((s) => s.id) : undefined,
        professionalId: professional.id,
        startsAt: selectedSlot.startsAt,
        name,
        phone,
        email,
        notes,
        recurrenceRule,
        recurrenceCount,
        timezone,
        lgpdAccepted,
        marketingConsent,
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
          {info.configuration?.logoUrl ? (
            <img
              alt={info.tenant.name}
              className="mx-auto h-16 w-16 rounded-2xl object-cover shadow-sm border border-[var(--color-line)]"
              src={info.configuration.logoUrl}
            />
          ) : (
            <span className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-[var(--color-primary)] text-white shadow-sm">
              <CalendarDays className="h-8 w-8" />
            </span>
          )}
          <h1 className="mt-3 font-headline text-2xl font-extrabold text-[var(--color-ink)]">{info.tenant.name}</h1>
          <p className="mt-1 text-sm font-semibold text-[var(--color-ink-soft)]">
            {info.configuration?.description || 'Agende seu horário em poucos passos.'}
          </p>

          {(info.configuration?.address || info.configuration?.contactWhatsApp || info.configuration?.instagram) && (
            <div className="mt-3 flex flex-wrap items-center justify-center gap-3 text-xs font-semibold text-[var(--color-ink-soft)]">
              {info.configuration.address && (
                <a
                  className="inline-flex items-center gap-1 rounded-full bg-[var(--color-bg-muted)] px-3 py-1 hover:text-[var(--color-primary)] transition"
                  href={`https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(
                    `${info.configuration.address}, ${info.configuration.addressNumber || ''} ${info.configuration.city || ''}`,
                  )}`}
                  rel="noreferrer"
                  target="_blank"
                >
                  <MapPin className="h-3.5 w-3.5 text-[var(--color-primary)]" />
                  <span>
                    {info.configuration.address}
                    {info.configuration.addressNumber ? `, ${info.configuration.addressNumber}` : ''}
                    {info.configuration.city ? ` - ${info.configuration.city}` : ''}
                  </span>
                </a>
              )}

              {info.configuration.contactWhatsApp && (
                <a
                  className="inline-flex items-center gap-1 rounded-full bg-emerald-50 text-emerald-700 px-3 py-1 hover:bg-emerald-100 transition"
                  href={`https://wa.me/55${info.configuration.contactWhatsApp.replace(/\D/g, '')}`}
                  rel="noreferrer"
                  target="_blank"
                >
                  <MessageCircle className="h-3.5 w-3.5" />
                  <span>WhatsApp</span>
                </a>
              )}

              {info.configuration.instagram && (
                <a
                  className="inline-flex items-center gap-1 rounded-full bg-pink-50 text-pink-700 px-3 py-1 hover:bg-pink-100 transition"
                  href={`https://instagram.com/${info.configuration.instagram.replace('@', '')}`}
                  rel="noreferrer"
                  target="_blank"
                >
                  <Instagram className="h-3.5 w-3.5" />
                  <span>Instagram</span>
                </a>
              )}

              {info.configuration.website && (
                <a
                  className="inline-flex items-center gap-1 rounded-full bg-[var(--color-bg-muted)] px-3 py-1 hover:text-[var(--color-primary)] transition"
                  href={info.configuration.website}
                  rel="noreferrer"
                  target="_blank"
                >
                  <Globe className="h-3.5 w-3.5" />
                  <span>Site</span>
                </a>
              )}
            </div>
          )}
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
            <ServiceStep onSelectServices={selectServices} services={services} />
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
              serviceName={
                selectedServices.length > 1
                  ? `${selectedServices.map((s) => s.name).join(' + ')} (Combo)`
                  : service.name
              }
              slots={slots}
            />
          ) : null}

          {step === 'details' && service && professional && selectedSlot ? (
            <DetailsStep
              onSubmit={handleBookingSubmit}
              professionalName={professional.name}
              serviceName={
                selectedServices.length > 1
                  ? `${selectedServices.map((s) => s.name).join(' + ')} (Combo)`
                  : service.name
              }
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

function ServiceStep({
  services,
  onSelectServices,
}: {
  services: Service[]
  onSelectServices: (selected: Service[]) => void
}) {
  const [comboMode, setComboMode] = useState(false)
  const [selectedIds, setSelectedIds] = useState<string[]>([])

  if (services.length === 0) {
    return <EmptyNote message="Este estabelecimento ainda não possui serviços disponíveis para agendamento." />
  }

  function toggleService(item: Service) {
    if (!comboMode) {
      onSelectServices([item])
      return
    }
    setSelectedIds((prev) =>
      prev.includes(item.id) ? prev.filter((id) => id !== item.id) : [...prev, item.id],
    )
  }

  const selectedList = services.filter((s) => selectedIds.includes(s.id))
  const totalDuration = selectedList.reduce((acc, s) => acc + s.durationMinutes, 0)
  const totalPrice = selectedList.reduce((acc, s) => acc + s.priceCents, 0)

  return (
    <div>
      <div className="flex flex-wrap items-center justify-between gap-2">
        <h2 className="font-headline text-lg font-extrabold text-[var(--color-ink)]">Escolha o serviço</h2>
        <button
          type="button"
          onClick={() => {
            setComboMode(!comboMode)
            setSelectedIds([])
          }}
          className={`rounded-full px-3 py-1 text-xs font-semibold transition ${
            comboMode
              ? 'bg-[var(--color-primary)] text-white'
              : 'bg-[var(--color-bg-muted)] text-[var(--color-ink)] hover:bg-[var(--color-line)]'
          }`}
        >
          {comboMode ? '✓ Modo Combo (Múltiplos Serviços)' : '+ Agendar Combo (Vários Serviços)'}
        </button>
      </div>

      <p className="mt-1 text-xs text-[var(--color-ink-soft)]">
        {comboMode
          ? 'Selecione dois ou mais serviços para realizar no mesmo dia.'
          : 'Clique no serviço desejado para continuar.'}
      </p>

      <ul className="mt-4 grid gap-3">
        {services.map((item) => {
          const isSelected = selectedIds.includes(item.id)
          return (
            <li key={item.id}>
              <button
                className={`flex w-full items-center gap-4 rounded-2xl border p-4 text-left transition ${
                  isSelected
                    ? 'border-[var(--color-primary)] bg-[var(--color-bg-elevated)] ring-2 ring-[var(--color-primary)]'
                    : 'border-[var(--color-line)] bg-[var(--color-bg-muted)] hover:border-[var(--color-primary)] hover:bg-[var(--color-bg-elevated)]'
                }`}
                onClick={() => toggleService(item)}
                type="button"
              >
                {comboMode && (
                  <div
                    className={`flex h-5 w-5 shrink-0 items-center justify-center rounded-md border ${
                      isSelected
                        ? 'border-[var(--color-primary)] bg-[var(--color-primary)] text-white'
                        : 'border-[var(--color-line)] bg-white'
                    }`}
                  >
                    {isSelected && <Check className="h-3.5 w-3.5" />}
                  </div>
                )}
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
          )
        })}
      </ul>

      {comboMode && selectedList.length > 0 && (
        <div className="mt-5 rounded-2xl border border-[var(--color-primary)] bg-[var(--color-bg-muted)] p-4 flex flex-wrap items-center justify-between gap-3">
          <div>
            <p className="text-sm font-bold text-[var(--color-ink)]">
              {selectedList.length} serviço{selectedList.length > 1 ? 's' : ''} no combo:
            </p>
            <p className="text-xs text-[var(--color-ink-soft)]">
              Duração total estimada: <span className="font-semibold text-[var(--color-ink)]">{totalDuration} min</span> • Total:{' '}
              <span className="font-semibold text-[var(--color-primary)]">{formatMoney(totalPrice)}</span>
            </p>
          </div>
          <Button onClick={() => onSelectServices(selectedList)}>
            Continuar com Combo ({selectedList.length})
          </Button>
        </div>
      )}
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

function LegalModal({
  type,
  onClose,
}: {
  type: 'terms' | 'lgpd'
  onClose: () => void
}) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm">
      <div className="relative max-h-[85vh] w-full max-w-lg overflow-y-auto rounded-3xl border border-[var(--color-line)] bg-white p-6 shadow-2xl">
        <div className="flex items-center justify-between pb-4 border-b border-[var(--color-line)]">
          <div className="flex items-center gap-2">
            <ShieldCheck className="h-5 w-5 text-[var(--color-primary)]" />
            <h3 className="font-headline text-lg font-black text-[var(--color-ink)]">
              {type === 'terms' ? 'Termos e Condições de Uso' : 'Privacidade e Proteção de Dados (LGPD)'}
            </h3>
          </div>
          <button
            onClick={onClose}
            className="rounded-lg p-1 text-stone-500 hover:bg-stone-100 hover:text-stone-800 transition"
            type="button"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        <div className="mt-4 space-y-4 text-sm text-[var(--color-ink-soft)] leading-relaxed">
          {type === 'terms' ? (
            <>
              <p>
                <strong>1. Agendamento e Pontualidade:</strong> Ao agendar um serviço, você reserva um horário exclusivo com o profissional selecionado. Recomendamos chegar com 5 a 10 minutos de antecedência. Em caso de atraso superior a 10 minutos, o atendimento poderá ser reduzido ou remarcado.
              </p>
              <p>
                <strong>2. Cancelamento e Reagendamento:</strong> Cancelamentos ou solicitações de reagendamento devem ser realizados através do link fornecido no código do agendamento, respeitando o prazo mínimo de antecedência configurado pelo estabelecimento.
              </p>
              <p>
                <strong>3. Política de Não Comparecimento (No-Show):</strong> O não comparecimento sem aviso prévio causa prejuízo à agenda dos profissionais e poderá resultar no bloqueio de novos agendamentos online para o mesmo número de contato.
              </p>
              <p>
                <strong>4. Pagamento:</strong> Caso aplicável, os pagamentos poderão ser realizados no momento do agendamento via Pix ou diretamente no local de atendimento através das formas aceitas pelo estabelecimento.
              </p>
            </>
          ) : (
            <>
              <p>
                <strong>1. Conformidade com a LGPD:</strong> Em conformidade com a Lei Geral de Proteção de Dados (Lei nº 13.709/2018), garantimos total transparência e segurança sobre o tratamento dos seus dados pessoais.
              </p>
              <p>
                <strong>2. Coleta e Finalidade:</strong> Coletamos apenas os dados essenciais (nome, telefone/WhatsApp e e-mail opcional) com a finalidade exclusiva de identificar seu agendamento, enviar confirmações, lembretes de horário e atualizações de status.
              </p>
              <p>
                <strong>3. Segurança e Retenção:</strong> Seus dados são transmitidos com criptografia SSL/TLS e armazenados em servidores seguros, com controles rigorosos de acesso restrito apenas à equipe do estabelecimento.
              </p>
              <p>
                <strong>4. Seus Direitos (Acesso e Anonimização):</strong> Você pode solicitar a qualquer momento a confirmação da existência de tratamento, a consulta, correção ou a anonimização/exclusão definitiva dos seus dados pessoais entrando em contato diretamente com o estabelecimento ou administrador do sistema.
              </p>
            </>
          )}
        </div>

        <div className="mt-6 flex justify-end">
          <Button onClick={onClose} type="button">
            Entendido
          </Button>
        </div>
      </div>
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
  const [lgpdAccepted, setLgpdAccepted] = useState(false)
  const [marketingConsent, setMarketingConsent] = useState(false)
  const [legalModal, setLegalModal] = useState<'terms' | 'lgpd' | null>(null)
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

        <Field hint="Deseja manter o mesmo horário nas próximas semanas?" label="Recorrência">
          <select
            name="recurrence"
            defaultValue="none"
            className="mt-2 w-full rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-3 text-sm font-semibold text-[var(--color-ink)] focus:border-[var(--color-primary)] focus:outline-none"
          >
            <option value="none">Somente nesta data (Agendamento Único)</option>
            <option value="weekly_2">Repetir semanalmente por 2 semanas</option>
            <option value="weekly_4">Repetir semanalmente por 4 semanas (1 mês)</option>
            <option value="weekly_8">Repetir semanalmente por 8 semanas (2 meses)</option>
          </select>
        </Field>

        <div className="flex items-center justify-between text-xs text-[var(--color-ink-soft)] px-1">
          <span>
            Fuso horário:{' '}
            <strong className="text-[var(--color-ink)]">
              {Intl.DateTimeFormat().resolvedOptions().timeZone || 'America/Sao_Paulo'}
            </strong>
          </span>
          <input
            type="hidden"
            name="timezone"
            value={Intl.DateTimeFormat().resolvedOptions().timeZone || 'America/Sao_Paulo'}
          />
        </div>

        <div className="space-y-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4">
          <label className="flex items-start gap-3 cursor-pointer">
            <input
              type="checkbox"
              name="lgpdAccepted"
              checked={lgpdAccepted}
              onChange={(e) => setLgpdAccepted(e.target.checked)}
              required
              className="mt-1 h-4 w-4 rounded border-stone-300 text-[var(--color-primary)] focus:ring-[var(--color-primary)]"
            />
            <span className="text-xs text-[var(--color-ink)] leading-normal">
              Li e concordo com os{' '}
              <button
                type="button"
                onClick={() => setLegalModal('terms')}
                className="font-bold text-[var(--color-primary)] underline hover:opacity-80"
              >
                Termos de Uso
              </button>{' '}
              e a{' '}
              <button
                type="button"
                onClick={() => setLegalModal('lgpd')}
                className="font-bold text-[var(--color-primary)] underline hover:opacity-80"
              >
                Política de Privacidade (LGPD)
              </button>
              . <span className="text-red-500 font-bold">*</span>
            </span>
          </label>

          <label className="flex items-start gap-3 cursor-pointer">
            <input
              type="checkbox"
              name="marketingConsent"
              checked={marketingConsent}
              onChange={(e) => setMarketingConsent(e.target.checked)}
              className="mt-1 h-4 w-4 rounded border-stone-300 text-[var(--color-primary)] focus:ring-[var(--color-primary)]"
            />
            <span className="text-xs text-[var(--color-ink-soft)] leading-normal">
              Aceito receber lembretes de horário, novidades e promoções exclusivas via WhatsApp ou E-mail.
            </span>
          </label>
        </div>

        {legalModal && <LegalModal type={legalModal} onClose={() => setLegalModal(null)} />}

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
