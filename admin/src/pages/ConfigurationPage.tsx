import { Save } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import { fetchConfiguration, updateConfiguration } from '../lib/adminApi'
import { DAY_LABELS } from '../lib/format'
import { Button } from '../shared/components/Button'
import { Checkbox, Field, Form, Input } from '../shared/components/FormControls'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { Configuration, WorkingHours } from '../types'

const EMPTY_DAY: Omit<WorkingHours, 'dayOfWeek'> = { openTime: '09:00', closeTime: '18:00', enabled: false }

export function ConfigurationPage() {
  const { showToast } = useToast()
  const [configuration, setConfiguration] = useState<Configuration | null>(null)
  const [establishmentName, setEstablishmentName] = useState('')
  const [workingHours, setWorkingHours] = useState<WorkingHours[]>([])
  const [appointmentIntervalMinutes, setAppointmentIntervalMinutes] = useState(30)
  const [defaultDurationMinutes, setDefaultDurationMinutes] = useState(30)
  const [allowOnlineBooking, setAllowOnlineBooking] = useState(true)
  const [advanceNoticeHours, setAdvanceNoticeHours] = useState(0)
  const [loading, setLoading] = useState(true)
  const [isSaving, setIsSaving] = useState(false)
  const [error, setError] = useState('')

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

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (!establishmentName.trim()) {
      setError('Informe o nome do estabelecimento.')
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

        {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
      </Form>
    </PageContainer>
  )
}
