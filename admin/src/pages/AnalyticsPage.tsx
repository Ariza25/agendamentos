import { CalendarCheck2, DollarSign, UserX, Users } from 'lucide-react'
import { useEffect, useState } from 'react'

import { fetchAppointmentAnalytics } from '../lib/adminApi'
import { fromDateKey, formatDate, formatMoney, todayKey } from '../lib/format'
import { Badge } from '../shared/components/Badge'
import { Input } from '../shared/components/FormControls'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { AppointmentAnalytics } from '../types'

function startOfMonthKey(): string {
  const date = new Date()
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-01`
}

export function AnalyticsPage() {
  const { showToast } = useToast()
  const [from, setFrom] = useState(startOfMonthKey())
  const [to, setTo] = useState(todayKey())
  const [analytics, setAnalytics] = useState<AppointmentAnalytics | null>(null)
  const [loading, setLoading] = useState(true)

  function loadAnalytics() {
    setLoading(true)
    fetchAppointmentAnalytics(`${from}T00:00:00`, `${to}T23:59:59`)
      .then(setAnalytics)
      .catch(() => {
        showToast({ message: 'Não foi possível carregar os indicadores.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  useEffect(loadAnalytics, [from, to])

  const maxServiceCount = Math.max(1, ...(analytics?.byService.map((item) => item.count) ?? [1]))
  const maxDayCount = Math.max(1, ...(analytics?.byDay.map((item) => item.count) ?? [1]))

  return (
    <PageContainer>
      <Section
        actions={
          <div className="flex flex-wrap items-end gap-3">
            <label className="block">
              <span className="text-xs font-bold uppercase tracking-wide text-[var(--color-ink-soft)]">De</span>
              <Input className="mt-1 w-auto" onChange={(event) => setFrom(event.target.value)} type="date" value={from} />
            </label>
            <label className="block">
              <span className="text-xs font-bold uppercase tracking-wide text-[var(--color-ink-soft)]">Até</span>
              <Input className="mt-1 w-auto" onChange={(event) => setTo(event.target.value)} type="date" value={to} />
            </label>
          </div>
        }
        subtitle="Acompanhe os indicadores de atendimento e receita do período."
        title="Dashboard"
      >
        <div className="mt-2" />
      </Section>

      {loading || !analytics ? (
        <Section>
          <Spinner />
        </Section>
      ) : (
        <>
          <Section title="Resumo do período">
            <div className="mt-4 grid gap-3 md:grid-cols-2 xl:grid-cols-5">
              <StatCard
                icon={<CalendarCheck2 className="h-5 w-5" />}
                label="Agendamentos"
                tone="accent"
                value={String(analytics.totalAppointments)}
              />
              <StatCard
                icon={<Users className="h-5 w-5" />}
                label="Confirmados"
                tone="success"
                value={String(analytics.confirmed)}
              />
              <StatCard
                icon={<CalendarCheck2 className="h-5 w-5" />}
                label="Concluídos"
                tone="success"
                value={String(analytics.completed)}
              />
              <StatCard
                icon={<UserX className="h-5 w-5" />}
                label="Não compareceram"
                tone="warning"
                value={`${analytics.noShow} (${Math.round(analytics.noShowRate * 100)}%)`}
              />
              <StatCard
                icon={<DollarSign className="h-5 w-5" />}
                label="Receita"
                tone="accent"
                value={formatMoney(analytics.revenueCents)}
              />
            </div>
          </Section>

          <div className="grid gap-6 xl:grid-cols-2">
            <Section title="Por serviço">
              {analytics.byService.length === 0 ? (
                <p className="mt-3 text-sm font-semibold text-[var(--color-ink-soft)]">
                  Nenhum atendimento no período selecionado.
                </p>
              ) : (
                <ul className="mt-4 grid gap-3">
                  {analytics.byService.map((item) => (
                    <li key={item.serviceId}>
                      <div className="flex items-center justify-between gap-3 text-sm">
                        <span className="truncate font-bold text-[var(--color-ink)]">{item.serviceName}</span>
                        <span className="shrink-0 font-semibold text-[var(--color-ink-soft)]">
                          {item.count} · {formatMoney(item.revenueCents)}
                        </span>
                      </div>
                      <div className="mt-1.5 h-2 overflow-hidden rounded-full bg-[var(--color-bg-muted)]">
                        <div
                          className="h-full rounded-full bg-[var(--color-primary)]"
                          style={{ width: `${(item.count / maxServiceCount) * 100}%` }}
                        />
                      </div>
                    </li>
                  ))}
                </ul>
              )}
            </Section>

            <Section title="Por profissional">
              {analytics.byProfessional.length === 0 ? (
                <p className="mt-3 text-sm font-semibold text-[var(--color-ink-soft)]">
                  Nenhum atendimento no período selecionado.
                </p>
              ) : (
                <ul className="mt-4 grid gap-2">
                  {analytics.byProfessional.map((item) => (
                    <li
                      className="flex items-center justify-between gap-3 rounded-xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] px-4 py-3 text-sm"
                      key={item.professionalId}
                    >
                      <span className="truncate font-bold text-[var(--color-ink)]">{item.professionalName}</span>
                      <Badge tone="neutral">{item.count} atendimentos</Badge>
                    </li>
                  ))}
                </ul>
              )}
            </Section>
          </div>

          <Section title="Atendimentos por dia">
            {analytics.byDay.length === 0 ? (
              <p className="mt-3 text-sm font-semibold text-[var(--color-ink-soft)]">
                Nenhum atendimento no período selecionado.
              </p>
            ) : (
              <div className="mt-4 flex h-40 items-end gap-1.5">
                {analytics.byDay.map((item) => (
                  <div
                    className="flex min-w-0 flex-1 flex-col items-center gap-1"
                    key={item.date}
                    title={`${formatDate(item.date)}: ${item.count}`}
                  >
                    <div
                      className="w-full rounded-t-lg bg-[var(--color-primary)] transition"
                      style={{ height: `${Math.max(6, (item.count / maxDayCount) * 100)}%` }}
                    />
                    <span className="text-[10px] font-bold text-[var(--color-ink-soft)]">
                      {fromDateKey(item.date).getDate()}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </Section>
        </>
      )}

    </PageContainer>
  )
}

function StatCard({ label, value, icon, tone }: { label: string; value: string; icon: React.ReactNode; tone: 'accent' | 'success' | 'warning' }) {
  const toneClasses: Record<string, string> = {
    accent: 'bg-orange-50 text-orange-900',
    success: 'bg-emerald-50 text-emerald-700',
    warning: 'bg-amber-50 text-amber-800',
  }
  return (
    <div className="rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4">
      <div className="flex items-center gap-2">
        <span className={`flex h-8 w-8 items-center justify-center rounded-xl ${toneClasses[tone]}`}>{icon}</span>
        <span className="truncate text-xs font-bold uppercase tracking-wide text-[var(--color-ink-soft)]">{label}</span>
      </div>
      <p className="mt-3 truncate font-headline text-2xl font-black text-[var(--color-ink)]">{value}</p>
    </div>
  )
}
