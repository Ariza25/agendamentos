import {
  Building2,
  Check,
  Copy,
  ExternalLink,
  Globe,
  Instagram,
  MapPin,
  Palette,
  Phone,
  Save,
  Sliders,
  Trash2,
} from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import { fetchConfiguration, updateConfiguration } from '../lib/adminApi'
import { buildApiUrl } from '../lib/apiBase'
import { DAY_LABELS } from '../lib/format'
import { Button } from '../shared/components/Button'
import { Checkbox, Field, Form, Input, Select } from '../shared/components/FormControls'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { ClosedDay, Configuration, WorkingHours } from '../types'

const EMPTY_DAY: Omit<WorkingHours, 'dayOfWeek'> = { openTime: '09:00', closeTime: '18:00', enabled: false }

const CHANNEL_OPTIONS = [
  { value: 'sms', label: 'SMS' },
  { value: 'email', label: 'Email' },
  { value: 'whatsapp', label: 'WhatsApp' },
]

const PRESET_COLORS = [
  { name: 'Azul Marinho (Padrão)', value: '#1e3a8a' },
  { name: 'Azul Cobalto', value: '#2563eb' },
  { name: 'Azul Gelo / Sky', value: '#0284c7' },
  { name: 'Esmeralda', value: '#059669' },
  { name: 'Âmbar Ouro', value: '#d97706' },
  { name: 'Grafite Elegante', value: '#1e293b' },
  { name: 'Carmim Nobre', value: '#b42318' },
]

type ActiveTab = 'rules' | 'branding'

export function ConfigurationPage() {
  const { showToast } = useToast()
  const [activeTab, setActiveTab] = useState<ActiveTab>('rules')
  const [configuration, setConfiguration] = useState<Configuration | null>(null)

  // Rules state
  const [establishmentName, setEstablishmentName] = useState('')
  const [workingHours, setWorkingHours] = useState<WorkingHours[]>([])
  const [appointmentIntervalMinutes, setAppointmentIntervalMinutes] = useState(30)
  const [defaultDurationMinutes, setDefaultDurationMinutes] = useState(30)
  const [allowOnlineBooking, setAllowOnlineBooking] = useState(true)
  const [advanceNoticeHours, setAdvanceNoticeHours] = useState(0)
  const [bufferMinutes, setBufferMinutes] = useState(0)
  const [closedDays, setClosedDays] = useState<ClosedDay[]>([])
  const [notificationsEnabled, setNotificationsEnabled] = useState(true)
  const [notificationReminderHours, setNotificationReminderHours] = useState(24)
  const [notificationChannel, setNotificationChannel] = useState('email')
  const [depositEnabled, setDepositEnabled] = useState(false)
  const [depositType, setDepositType] = useState('percentage')
  const [depositAmountCents, setDepositAmountCents] = useState(0)
  const [depositPercentage, setDepositPercentage] = useState(20)
  const [noShowFeeCents, setNoShowFeeCents] = useState(0)
  const [cancellationPolicyHours, setCancellationPolicyHours] = useState(24)
  const [calendarSyncEnabled, setCalendarSyncEnabled] = useState(false)

  // Branding & White-label state
  const [logoUrl, setLogoUrl] = useState('')
  const [primaryColor, setPrimaryColor] = useState('#1e3a8a')
  const [address, setAddress] = useState('')
  const [addressNumber, setAddressNumber] = useState('')
  const [neighborhood, setNeighborhood] = useState('')
  const [city, setCity] = useState('')
  const [state, setState] = useState('')
  const [postalCode, setPostalCode] = useState('')
  const [contactPhone, setContactPhone] = useState('')
  const [contactWhatsApp, setContactWhatsApp] = useState('')
  const [instagram, setInstagram] = useState('')
  const [facebook, setFacebook] = useState('')
  const [website, setWebsite] = useState('')
  const [description, setDescription] = useState('')

  const [loading, setLoading] = useState(true)
  const [isSaving, setIsSaving] = useState(false)
  const [error, setError] = useState('')
  const [copied, setCopied] = useState(false)

  const icsUrl = buildApiUrl('/api/v1/admin/calendar/export.ics')

  useEffect(() => {
    let active = true
    fetchConfiguration()
      .then((data) => {
        if (!active) {
          return
        }
        setConfiguration(data)
        setEstablishmentName(data.establishmentName)
        setWorkingHours(data.workingHours)
        setAppointmentIntervalMinutes(data.appointmentIntervalMinutes)
        setDefaultDurationMinutes(data.defaultDurationMinutes)
        setAllowOnlineBooking(data.allowOnlineBooking)
        setAdvanceNoticeHours(data.advanceNoticeHours)
        setBufferMinutes(data.bufferMinutes ?? 0)
        setClosedDays(data.closedDays ?? [])
        setNotificationsEnabled(data.notificationsEnabled ?? true)
        setNotificationReminderHours(data.notificationReminderHours ?? 24)
        setNotificationChannel(data.notificationChannel || 'email')
        setDepositEnabled(data.depositEnabled ?? false)
        setDepositType(data.depositType || 'percentage')
        setDepositAmountCents(data.depositAmountCents ?? 0)
        setDepositPercentage(data.depositPercentage ?? 20)
        setNoShowFeeCents(data.noShowFeeCents ?? 0)
        setCancellationPolicyHours(data.cancellationPolicyHours ?? 24)
        setCalendarSyncEnabled(data.calendarSyncEnabled ?? false)

        // Branding
        setLogoUrl(data.logoUrl || '')
        setPrimaryColor(data.primaryColor || '#1e3a8a')
        setAddress(data.address || '')
        setAddressNumber(data.addressNumber || '')
        setNeighborhood(data.neighborhood || '')
        setCity(data.city || '')
        setState(data.state || '')
        setPostalCode(data.postalCode || '')
        setContactPhone(data.contactPhone || '')
        setContactWhatsApp(data.contactWhatsApp || '')
        setInstagram(data.instagram || '')
        setFacebook(data.facebook || '')
        setWebsite(data.website || '')
        setDescription(data.description || '')
      })
      .catch(() => {
        if (active) {
          showToast({ message: 'Não foi possível carregar as configurações.', tone: 'error' })
        }
      })
      .finally(() => {
        if (active) {
          setLoading(false)
        }
      })

    return () => {
      active = false
    }
  }, [])

  function updateDay(dayOfWeek: number, patch: Partial<WorkingHours>) {
    setWorkingHours((current) => current.map((day) => (day.dayOfWeek === dayOfWeek ? { ...day, ...patch } : day)))
  }

  function addClosedDay() {
    setClosedDays((current) => [
      ...current,
      { date: '', name: '', allDay: true },
    ])
  }

  function updateClosedDay(index: number, patch: Partial<ClosedDay>) {
    setClosedDays((current) => current.map((item, itemIndex) => (itemIndex === index ? { ...item, ...patch } : item)))
  }

  function removeClosedDay(index: number) {
    setClosedDays((current) => current.filter((_, itemIndex) => itemIndex !== index))
  }

  async function copyIcsUrl() {
    try {
      await navigator.clipboard.writeText(icsUrl)
      setCopied(true)
      window.setTimeout(() => setCopied(false), 2000)
    } catch {
      showToast({ message: 'Não foi possível copiar o link.', tone: 'error' })
    }
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (!establishmentName.trim()) {
      setError('Informe o nome do estabelecimento.')
      return
    }
    if (closedDays.some((day) => !day.date || !day.name.trim())) {
      setError('Preencha a data e o nome de cada dia de fechamento.')
      return
    }

    setIsSaving(true)
    setError('')

    try {
      const updated = await updateConfiguration({
        establishmentName: establishmentName.trim(),
        workingHours,
        appointmentIntervalMinutes,
        defaultDurationMinutes,
        allowOnlineBooking,
        advanceNoticeHours,
        bufferMinutes,
        closedDays,
        notificationsEnabled,
        notificationReminderHours,
        notificationChannel,
        depositEnabled,
        depositType,
        depositAmountCents,
        depositPercentage,
        noShowFeeCents,
        cancellationPolicyHours,
        calendarSyncEnabled,
        // Branding & Contacts
        logoUrl: logoUrl.trim(),
        primaryColor: primaryColor.trim() || '#1e3a8a',
        address: address.trim(),
        addressNumber: addressNumber.trim(),
        neighborhood: neighborhood.trim(),
        city: city.trim(),
        state: state.trim(),
        postalCode: postalCode.trim(),
        contactPhone: contactPhone.trim(),
        contactWhatsApp: contactWhatsApp.trim(),
        instagram: instagram.trim(),
        facebook: facebook.trim(),
        website: website.trim(),
        description: description.trim(),
      })
      setConfiguration(updated)
      showToast({ message: 'Configurações salvas com sucesso.', tone: 'success' })
    } catch (submitError) {
      setError((submitError as Error).message)
    } finally {
      setIsSaving(false)
    }
  }

  if (loading) {
    return (
      <PageContainer>
        <Section>
          <Spinner />
        </Section>
      </PageContainer>
    )
  }

  const fullAddressPreview = [
    address,
    addressNumber,
    neighborhood,
    city && state ? `${city} - ${state}` : city || state,
  ]
    .filter(Boolean)
    .join(', ')

  return (
    <PageContainer>
      <Section
        actions={
          <Button icon={<Save className="h-4 w-4" />} loading={isSaving} type="submit" form="configuration-form">
            Salvar configurações
          </Button>
        }
        subtitle="Defina os horários de atendimento, regras da agenda e personalize a identidade visual da sua página pública."
        title="Configurações & Marca"
      >
        {/* Navigation Tabs */}
        <div className="mt-4 flex border-b border-[var(--color-line)]">
          <button
            type="button"
            onClick={() => setActiveTab('rules')}
            className={`flex items-center gap-2 border-b-2 px-4 py-2.5 text-sm font-semibold transition-colors ${
              activeTab === 'rules'
                ? 'border-[var(--color-primary)] text-[var(--color-primary)]'
                : 'border-transparent text-[var(--color-ink-soft)] hover:text-[var(--color-ink)]'
            }`}
          >
            <Sliders className="h-4 w-4" />
            Agendamento & Regras
          </button>
          <button
            type="button"
            onClick={() => setActiveTab('branding')}
            className={`flex items-center gap-2 border-b-2 px-4 py-2.5 text-sm font-semibold transition-colors ${
              activeTab === 'branding'
                ? 'border-[var(--color-primary)] text-[var(--color-primary)]'
                : 'border-transparent text-[var(--color-ink-soft)] hover:text-[var(--color-ink)]'
            }`}
          >
            <Palette className="h-4 w-4" />
            Página Pública & Marca
          </button>
        </div>
      </Section>

      <Form className="mt-6 grid gap-6" id="configuration-form" onSubmit={handleSubmit}>
        {activeTab === 'rules' ? (
          <>
            <Section subtitle="Identificação básica exibida aos clientes." title="Estabelecimento">
              <Field label="Nome do estabelecimento">
                <Input
                  className="mt-2"
                  onChange={(event) => setEstablishmentName(event.target.value)}
                  value={establishmentName}
                  placeholder="Ex.: Barbearia Vintage / Clínica Bella"
                />
              </Field>
            </Section>

            <Section
              subtitle="Configure cada dia da semana. Horários no formato HH:MM (24h)."
              title="Períodos de atendimento (working hours)"
            >
              <div className="mt-4 grid gap-3 md:grid-cols-2">
                {Array.from({ length: 7 }, (_, dayOfWeek) => dayOfWeek).map((dayOfWeek) => {
                  const day = workingHours.find((item) => item.dayOfWeek === dayOfWeek) ?? { ...EMPTY_DAY, dayOfWeek }
                  return (
                    <div
                      className="grid grid-cols-[auto_1fr_1fr] items-center gap-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] px-4 py-3"
                      key={dayOfWeek}
                    >
                      <div className="flex min-w-0 items-center gap-2">
                        <Checkbox
                          checked={day.enabled}
                          onChange={(event) => updateDay(dayOfWeek, { enabled: event.target.checked })}
                        />
                        <span className="truncate text-sm font-bold text-[var(--color-ink)]">{DAY_LABELS[dayOfWeek]}</span>
                      </div>
                      <Input
                        aria-label={`Abertura ${DAY_LABELS[dayOfWeek]}`}
                        disabled={!day.enabled}
                        onChange={(event) => updateDay(dayOfWeek, { openTime: event.target.value })}
                        type="time"
                        value={day.openTime}
                      />
                      <Input
                        aria-label={`Fechamento ${DAY_LABELS[dayOfWeek]}`}
                        disabled={!day.enabled}
                        onChange={(event) => updateDay(dayOfWeek, { closeTime: event.target.value })}
                        type="time"
                        value={day.closeTime}
                      />
                    </div>
                  )
                })}
              </div>
              <p className="mt-3 text-xs font-semibold text-[var(--color-ink-soft)]">
                Dias desabilitados não geram horários de agendamento online.
              </p>
            </Section>

            <Section subtitle="Regras que controlam a geração de horários disponíveis." title="Agendamento online">
              <div className="mt-4 grid gap-4 md:grid-cols-3">
                <Field label="Intervalo entre horários (minutos)">
                  <Input
                    className="mt-2"
                    min={5}
                    onChange={(event) => setAppointmentIntervalMinutes(Number(event.target.value))}
                    type="number"
                    value={appointmentIntervalMinutes}
                  />
                </Field>
                <Field label="Duração padrão (minutos)">
                  <Input
                    className="mt-2"
                    min={5}
                    onChange={(event) => setDefaultDurationMinutes(Number(event.target.value))}
                    type="number"
                    value={defaultDurationMinutes}
                  />
                </Field>
                <Field hint="0 = sem restrição. Usado para não aceitar agendamentos em cima da hora." label="Antecedência mínima (horas)">
                  <Input
                    className="mt-2"
                    min={0}
                    onChange={(event) => setAdvanceNoticeHours(Number(event.target.value))}
                    type="number"
                    value={advanceNoticeHours}
                  />
                </Field>
              </div>
              <div className="mt-4 grid gap-4 md:grid-cols-2">
                <Field hint="Tempo reservado entre um agendamento e outro de um mesmo profissional." label="Intervalo de segurança / buffer (minutos)">
                  <Input
                    className="mt-2"
                    min={0}
                    onChange={(event) => setBufferMinutes(Number(event.target.value))}
                    type="number"
                    value={bufferMinutes}
                  />
                </Field>
              </div>
              <div className="mt-4 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] px-4 py-3">
                <Checkbox checked={allowOnlineBooking} onChange={(event) => setAllowOnlineBooking(event.target.checked)} />
                <span className="ml-2 text-sm font-bold text-[var(--color-ink)]">Permitir agendamento online pelo cliente final</span>
              </div>
              <p className="mt-3 text-xs font-semibold text-[var(--color-ink-soft)]">
                {configuration?.tenantId
                  ? `Link público de agendamento habilitado para o estabelecimento "${establishmentName.trim() || configuration.tenantId}".`
                  : ''}
              </p>
            </Section>

            <Section
              actions={
                <Button onClick={addClosedDay} size="sm" variant="secondary">
                  + Adicionar
                </Button>
              }
              subtitle="Dias em que o estabelecimento não abre, como feriados. Podem ter horário reduzido."
              title="Dias de fechamento"
            >
              {closedDays.length === 0 ? (
                <p className="mt-3 text-sm font-semibold text-[var(--color-ink-soft)]">
                  Nenhum dia de fechamento cadastrado.
                </p>
              ) : (
                <ul className="mt-4 grid gap-3">
                  {closedDays.map((day, index) => (
                    <li
                      className="grid gap-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] px-4 py-3 md:grid-cols-[auto_1fr_1fr_1fr_1fr_auto] md:items-center"
                      key={`${day.date}-${index}`}
                    >
                      <div className="flex items-center gap-2">
                        <Checkbox
                          checked={day.allDay}
                          onChange={(event) => updateClosedDay(index, { allDay: event.target.checked })}
                        />
                        <span className="text-xs font-bold uppercase tracking-wide text-[var(--color-ink-soft)]">Dia todo</span>
                      </div>
                      <Input
                        aria-label="Data"
                        onChange={(event) => updateClosedDay(index, { date: event.target.value })}
                        type="date"
                        value={day.date}
                      />
                      <Input
                        aria-label="Nome do dia"
                        onChange={(event) => updateClosedDay(index, { name: event.target.value })}
                        placeholder="Ex.: Natal"
                        value={day.name}
                      />
                      <Input
                        aria-label="Abertura"
                        disabled={day.allDay}
                        onChange={(event) => updateClosedDay(index, { openTime: event.target.value })}
                        type="time"
                        value={day.openTime ?? ''}
                      />
                      <Input
                        aria-label="Fechamento"
                        disabled={day.allDay}
                        onChange={(event) => updateClosedDay(index, { closeTime: event.target.value })}
                        type="time"
                        value={day.closeTime ?? ''}
                      />
                      <Button
                        aria-label="Remover dia"
                        icon={<Trash2 className="h-4 w-4" />}
                        onClick={() => removeClosedDay(index)}
                        size="sm"
                        variant="danger"
                      />
                    </li>
                  ))}
                </ul>
              )}
            </Section>

            <Section subtitle="Lembretes automáticos enviados aos clientes antes do atendimento." title="Notificações e lembretes">
              <div className="mt-4 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] px-4 py-3">
                <Checkbox checked={notificationsEnabled} onChange={(event) => setNotificationsEnabled(event.target.checked)} />
                <span className="ml-2 text-sm font-bold text-[var(--color-ink)]">Enviar lembretes automáticos</span>
              </div>
              <div className="mt-4 grid gap-4 md:grid-cols-2">
                <Field label="Antecedência do lembrete (horas)">
                  <Input
                    className="mt-2"
                    min={1}
                    onChange={(event) => setNotificationReminderHours(Number(event.target.value))}
                    type="number"
                    value={notificationReminderHours}
                  />
                </Field>
                <Field label="Canal preferencial">
                  <Select className="mt-2" onChange={(event) => setNotificationChannel(event.target.value)} value={notificationChannel}>
                    {CHANNEL_OPTIONS.map((option) => (
                      <option key={option.value} value={option.value}>
                        {option.label}
                      </option>
                    ))}
                  </Select>
                </Field>
              </div>
            </Section>

            <Section subtitle="Depósito e política de cancelamento definidos para o agendamento online." title="Depósito, taxas e cancelamento">
              <div className="mt-4 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] px-4 py-3">
                <Checkbox checked={depositEnabled} onChange={(event) => setDepositEnabled(event.target.checked)} />
                <span className="ml-2 text-sm font-bold text-[var(--color-ink)]">Exigir depósito para agendamento online</span>
              </div>
              <div className="mt-4 grid gap-4 md:grid-cols-3">
                <Field label="Tipo de depósito">
                  <Select
                    className="mt-2"
                    disabled={!depositEnabled}
                    onChange={(event) => setDepositType(event.target.value)}
                    value={depositType}
                  >
                    <option value="percentage">Percentual</option>
                    <option value="fixed">Valor fixo</option>
                  </Select>
                </Field>
                {depositType === 'percentage' ? (
                  <Field label="Percentual (%)">
                    <Input
                      className="mt-2"
                      disabled={!depositEnabled}
                      min={1}
                      max={100}
                      onChange={(event) => setDepositPercentage(Number(event.target.value))}
                      type="number"
                      value={depositPercentage}
                    />
                  </Field>
                ) : (
                  <Field label="Valor do depósito (R$)">
                    <Input
                      className="mt-2"
                      disabled={!depositEnabled}
                      min={0}
                      onChange={(event) => setDepositAmountCents(Math.round(Number(event.target.value) * 100))}
                      step="0.01"
                      type="number"
                      value={depositAmountCents / 100}
                    />
                  </Field>
                )}
                <Field hint="Cobrado quando o cliente não comparece ao agendamento." label="Taxa de não comparecimento (R$)">
                  <Input
                    className="mt-2"
                    min={0}
                    onChange={(event) => setNoShowFeeCents(Math.round(Number(event.target.value) * 100))}
                    step="0.01"
                    type="number"
                    value={noShowFeeCents / 100}
                  />
                </Field>
              </div>
              <div className="mt-4 grid gap-4 md:grid-cols-2">
                <Field hint="Janela para o cliente cancelar ou reagendar sem custo." label="Política de cancelamento (horas)">
                  <Input
                    className="mt-2"
                    min={0}
                    onChange={(event) => setCancellationPolicyHours(Number(event.target.value))}
                    type="number"
                    value={cancellationPolicyHours}
                  />
                </Field>
              </div>
            </Section>

            <Section subtitle="Exporte seus agendamentos para o Google Calendar ou outro calendário via link .ics." title="Sincronização de calendário">
              <div className="mt-4 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] px-4 py-3">
                <Checkbox checked={calendarSyncEnabled} onChange={(event) => setCalendarSyncEnabled(event.target.checked)} />
                <span className="ml-2 text-sm font-bold text-[var(--color-ink)]">Habilitar exportação de agenda (.ics)</span>
              </div>
              {calendarSyncEnabled ? (
                <div className="mt-3 grid gap-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] px-4 py-3">
                  <p className="text-sm font-bold text-[var(--color-ink)]">Link de assinatura (.ics)</p>
                  <div className="flex flex-wrap items-center gap-2">
                    <code className="min-w-0 flex-1 truncate rounded-lg bg-[var(--color-bg)] px-3 py-2 text-xs font-semibold text-[var(--color-ink)]">
                      {icsUrl}
                    </code>
                    <Button
                      icon={copied ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
                      onClick={() => void copyIcsUrl()}
                      size="sm"
                      variant="secondary"
                    >
                      {copied ? 'Copiado!' : 'Copiar'}
                    </Button>
                  </div>
                  <p className="text-xs font-semibold text-[var(--color-ink-soft)]">
                    Cole o link em "Adicionar por URL" do Google Calendar ou abra-o no navegador para baixar o arquivo .ics.
                  </p>
                </div>
              ) : null}
            </Section>
          </>
        ) : (
          /* Branding & White-Label Tab */
          <>
            <Section
              subtitle="Personalize a aparência, logotipo e paleta de cores exclusiva da página pública de agendamento do seu estabelecimento."
              title="Identidade Visual & Cores"
            >
              <div className="mt-4 grid gap-6 md:grid-cols-2">
                <Field hint="URL de uma imagem pública quadrada ou horizontal (PNG, JPG, SVG)." label="URL do Logotipo">
                  <Input
                    className="mt-2"
                    placeholder="https://exemplo.com/sua-logo.png"
                    value={logoUrl}
                    onChange={(e) => setLogoUrl(e.target.value)}
                  />
                  {logoUrl ? (
                    <div className="mt-3 flex items-center gap-3 rounded-xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-2">
                      <img
                        src={logoUrl}
                        alt="Prévia da Logo"
                        className="h-10 w-10 rounded-lg object-contain bg-white p-1 border border-border"
                        onError={(e) => {
                          ;(e.target as HTMLElement).style.display = 'none'
                        }}
                      />
                      <span className="text-xs font-semibold text-[var(--color-ink-soft)]">Prévia da imagem carregada</span>
                    </div>
                  ) : null}
                </Field>

                <div>
                  <label className="text-sm font-bold text-[var(--color-ink)]">Cor Primária da Marca</label>
                  <p className="text-xs text-[var(--color-ink-soft)] mt-0.5">
                    Utilizada nos botões, datas selecionadas e destaques da página de agendamento.
                  </p>
                  <div className="mt-3 flex items-center gap-3">
                    <input
                      type="color"
                      aria-label="Seletor de cor primária"
                      value={primaryColor || '#1e3a8a'}
                      onChange={(e) => setPrimaryColor(e.target.value)}
                      className="h-10 w-14 cursor-pointer rounded-lg border border-[var(--color-line)] p-1 bg-transparent"
                    />
                    <Input
                      className="w-32 font-mono text-sm"
                      placeholder="#1e3a8a"
                      value={primaryColor}
                      onChange={(e) => setPrimaryColor(e.target.value)}
                    />
                  </div>

                  {/* Preset swatches */}
                  <div className="mt-3 flex flex-wrap gap-2">
                    {PRESET_COLORS.map((preset) => (
                      <button
                        key={preset.value}
                        type="button"
                        onClick={() => setPrimaryColor(preset.value)}
                        className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium border transition-transform ${
                          primaryColor.toLowerCase() === preset.value.toLowerCase()
                            ? 'border-foreground shadow-sm scale-105'
                            : 'border-border opacity-80 hover:opacity-100'
                        }`}
                      >
                        <span
                          className="h-3 w-3 rounded-full"
                          style={{ backgroundColor: preset.value }}
                        />
                        {preset.name}
                      </button>
                    ))}
                  </div>
                </div>
              </div>

              <div className="mt-6">
                <Field
                  hint="Uma breve mensagem ou slogan exibido no topo da página de agendamento do seu cliente."
                  label="Descrição / Apresentação do Estabelecimento"
                >
                  <textarea
                    rows={3}
                    className="mt-2 w-full rounded-xl border border-[var(--color-line)] bg-[var(--color-bg)] px-3 py-2 text-sm text-[var(--color-ink)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
                    placeholder="Ex.: Especialistas em barbearia clássica e cuidados masculinos. Atendimento com hora marcada e estacionamento no local."
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                  />
                </Field>
              </div>
            </Section>

            <Section
              subtitle="Endereço físico exibido no rodapé do agendamento com atalho direto para navegação via GPS / Google Maps."
              title="Endereço & Localização"
            >
              <div className="mt-4 grid gap-4 md:grid-cols-3">
                <div className="md:col-span-2">
                  <Field label="Rua / Logradouro">
                    <Input
                      className="mt-2"
                      placeholder="Ex.: Av. Paulista"
                      value={address}
                      onChange={(e) => setAddress(e.target.value)}
                    />
                  </Field>
                </div>
                <div>
                  <Field label="Número / Compl.">
                    <Input
                      className="mt-2"
                      placeholder="Ex.: 1000, Sala 42"
                      value={addressNumber}
                      onChange={(e) => setAddressNumber(e.target.value)}
                    />
                  </Field>
                </div>
              </div>

              <div className="mt-4 grid gap-4 md:grid-cols-4">
                <div>
                  <Field label="Bairro">
                    <Input
                      className="mt-2"
                      placeholder="Ex.: Bela Vista"
                      value={neighborhood}
                      onChange={(e) => setNeighborhood(e.target.value)}
                    />
                  </Field>
                </div>
                <div>
                  <Field label="Cidade">
                    <Input
                      className="mt-2"
                      placeholder="Ex.: São Paulo"
                      value={city}
                      onChange={(e) => setCity(e.target.value)}
                    />
                  </Field>
                </div>
                <div>
                  <Field label="Estado (UF)">
                    <Input
                      className="mt-2"
                      placeholder="Ex.: SP"
                      maxLength={2}
                      value={state}
                      onChange={(e) => setState(e.target.value.toUpperCase())}
                    />
                  </Field>
                </div>
                <div>
                  <Field label="CEP">
                    <Input
                      className="mt-2"
                      placeholder="Ex.: 01310-100"
                      value={postalCode}
                      onChange={(e) => setPostalCode(e.target.value)}
                    />
                  </Field>
                </div>
              </div>
            </Section>

            <Section
              subtitle="Canais diretos para os clientes tirarem dúvidas antes ou depois de agendarem."
              title="Contatos & Redes Sociais"
            >
              <div className="mt-4 grid gap-4 md:grid-cols-2">
                <Field
                  hint="Número com DDD (ex: 11999998888). Cria botão direto para iniciar conversa no WhatsApp."
                  label="WhatsApp de Contato"
                >
                  <div className="relative mt-2">
                    <Input
                      placeholder="Ex.: 11999998888"
                      value={contactWhatsApp}
                      onChange={(e) => setContactWhatsApp(e.target.value)}
                    />
                  </div>
                </Field>

                <Field label="Telefone Fixo / Comercial">
                  <Input
                    className="mt-2"
                    placeholder="Ex.: (11) 3222-1111"
                    value={contactPhone}
                    onChange={(e) => setContactPhone(e.target.value)}
                  />
                </Field>
              </div>

              <div className="mt-4 grid gap-4 md:grid-cols-3">
                <Field hint="Ex: @seunegocio ou link completo" label="Instagram">
                  <Input
                    className="mt-2"
                    placeholder="@seunegocio"
                    value={instagram}
                    onChange={(e) => setInstagram(e.target.value)}
                  />
                </Field>

                <Field label="Facebook">
                  <Input
                    className="mt-2"
                    placeholder="fb.com/seunegocio"
                    value={facebook}
                    onChange={(e) => setFacebook(e.target.value)}
                  />
                </Field>

                <Field label="Website Oficial">
                  <Input
                    className="mt-2"
                    placeholder="https://seusite.com.br"
                    value={website}
                    onChange={(e) => setWebsite(e.target.value)}
                  />
                </Field>
              </div>
            </Section>

            {/* Live Interactive Preview */}
            <Section
              subtitle="Veja como o cabeçalho e as informações do seu estabelecimento serão exibidos aos clientes."
              title="Pré-visualização da Sua Página Pública"
            >
              <div
                className="mt-4 overflow-hidden rounded-2xl border border-[var(--color-line)] bg-white p-6 shadow-sm transition-all"
                style={{
                  borderTop: `6px solid ${primaryColor || '#1e3a8a'}`,
                }}
              >
                <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                  <div className="flex items-center gap-4">
                    {logoUrl ? (
                      <img
                        src={logoUrl}
                        alt="Logo"
                        className="h-16 w-16 rounded-2xl object-cover shadow-sm border border-slate-100"
                        onError={(e) => {
                          ;(e.target as HTMLElement).style.display = 'none'
                        }}
                      />
                    ) : (
                      <div
                        className="flex h-16 w-16 items-center justify-center rounded-2xl text-white shadow-sm"
                        style={{ backgroundColor: primaryColor || '#1e3a8a' }}
                      >
                        <Building2 className="h-8 w-8" />
                      </div>
                    )}
                    <div>
                      <h3 className="text-xl font-bold text-slate-900">
                        {establishmentName.trim() || 'Nome do Seu Estabelecimento'}
                      </h3>
                      <p className="mt-0.5 max-w-xl text-sm text-slate-600 line-clamp-2">
                        {description || 'Sua bio ou descrição atrativa aparecerá aqui para os clientes.'}
                      </p>
                    </div>
                  </div>

                  <button
                    type="button"
                    style={{ backgroundColor: primaryColor || '#1e3a8a' }}
                    className="inline-flex items-center justify-center rounded-xl px-4 py-2 text-sm font-semibold text-white shadow-sm transition-opacity hover:opacity-90"
                  >
                    Agendar Horário
                  </button>
                </div>

                {/* Footer preview */}
                <div className="mt-6 flex flex-wrap items-center justify-between gap-4 border-t border-slate-100 pt-4 text-xs text-slate-500">
                  <div className="flex items-center gap-1.5 font-medium">
                    <MapPin className="h-3.5 w-3.5 text-slate-400" />
                    <span>
                      {fullAddressPreview || 'Endereço do estabelecimento, Cidade - UF'}
                    </span>
                  </div>

                  <div className="flex items-center gap-3">
                    {contactWhatsApp ? (
                      <span className="flex items-center gap-1 font-semibold text-emerald-600">
                        <Phone className="h-3.5 w-3.5" />
                        WhatsApp: {contactWhatsApp}
                      </span>
                    ) : null}
                    {instagram ? (
                      <span className="flex items-center gap-1 font-semibold text-pink-600">
                        <Instagram className="h-3.5 w-3.5" />
                        {instagram.startsWith('@') ? instagram : `@${instagram}`}
                      </span>
                    ) : null}
                    {website ? (
                      <span className="flex items-center gap-1 font-medium text-slate-600">
                        <Globe className="h-3.5 w-3.5" />
                        Website
                      </span>
                    ) : null}
                  </div>
                </div>
              </div>
            </Section>
          </>
        )}

        {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
      </Form>
    </PageContainer>
  )
}
