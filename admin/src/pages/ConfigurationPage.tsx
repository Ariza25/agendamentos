import { Check, Copy, Save, Trash2 } from 'lucide-react'
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

export function ConfigurationPage() {
  const { showToast } = useToast()
  const [configuration, setConfiguration] = useState<Configuration | null>(null)
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
      })
      setConfiguration(updated)
      showToast({ message: 'Configurações salvas.', tone: 'success' })
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

  return (
    <PageContainer>
      <Section
        actions={
          <Button icon={<Save className="h-4 w-4" />} loading={isSaving} type="submit" form="configuration-form">
            Salvar configurações
          </Button>
        }
        subtitle="Defina os períodos de atendimento e as regras do agendamento online."
        title="Configuração"
      >
        <div className="mt-2" />
      </Section>

      <Form className="mt-0 grid gap-6" id="configuration-form" onSubmit={handleSubmit}>
        <Section subtitle="Identificação exibida na tela pública de agendamento." title="Estabelecimento">
          <Field label="Nome do estabelecimento">
            <Input className="mt-2" onChange={(event) => setEstablishmentName(event.target.value)} value={establishmentName} />
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

        {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
      </Form>
    </PageContainer>
  )
}
