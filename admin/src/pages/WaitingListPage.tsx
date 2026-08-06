import { Clock3 } from 'lucide-react'
import { useEffect, useState } from 'react'

import { fetchWaitingList, updateWaitingListStatus } from '../lib/adminApi'
import { formatDateTime } from '../lib/format'
import { useReferenceData } from '../shared/hooks/useReferenceData'
import { Badge } from '../shared/components/Badge'
import { Button } from '../shared/components/Button'
import { EmptyState } from '../shared/components/EmptyState'
import { Select } from '../shared/components/FormControls'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { WaitingListEntry } from '../types'

const STATUS_LABELS: Record<string, string> = {
  waiting: 'Aguardando',
  notified: 'Notificado',
  converted: 'Convertido',
  cancelled: 'Cancelado',
}

const STATUS_TONES: Record<string, 'neutral' | 'success' | 'warning' | 'danger'> = {
  waiting: 'warning',
  notified: 'warning',
  converted: 'success',
  cancelled: 'danger',
}

const NEXT_STATUSES: Record<string, string[]> = {
  waiting: ['notified', 'converted', 'cancelled'],
  notified: ['converted', 'cancelled'],
  converted: [],
  cancelled: [],
}

export function WaitingListPage() {
  const { showToast } = useToast()
  const referenceData = useReferenceData()
  const [status, setStatus] = useState('')
  const [entries, setEntries] = useState<WaitingListEntry[]>([])
  const [loading, setLoading] = useState(true)

  function loadEntries() {
    setLoading(true)
    fetchWaitingList(status)
      .then(setEntries)
      .catch(() => {
        showToast({ message: 'Não foi possível carregar a lista de espera.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  useEffect(loadEntries, [status])

  async function handleStatusChange(entry: WaitingListEntry, nextStatus: string) {
    try {
      await updateWaitingListStatus(entry.id, nextStatus)
      showToast({ message: 'Situação atualizada.', tone: 'success' })
      loadEntries()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    }
  }

  return (
    <PageContainer>
      <Section
        actions={
          <label className="block">
            <span className="text-xs font-bold uppercase tracking-wide text-[var(--color-ink-soft)]">Situação</span>
            <Select className="mt-1 w-auto" onChange={(event) => setStatus(event.target.value)} value={status}>
              <option value="">Todas</option>
              {Object.entries(STATUS_LABELS).map(([value, label]) => (
                <option key={value} value={value}>
                  {label}
                </option>
              ))}
            </Select>
          </label>
        }
        subtitle="Clientes que pediram aviso quando um horário ficar livre."
        title="Lista de espera"
      >
        <div className="mt-2" />
      </Section>

      {loading ? (
        <Section>
          <Spinner />
        </Section>
      ) : entries.length === 0 ? (
        <Section>
          <EmptyState
            description="Clientes inscritos na lista de espera do agendamento online aparecerão aqui."
            icon={<Clock3 className="h-10 w-10" />}
            title="Nenhuma inscrição na lista de espera"
          />
        </Section>
      ) : (
        <Section>
          <ul className="mt-2 grid gap-3">
            {entries.map((entry) => {
              const service = entry.serviceId ? referenceData.byId.service.get(entry.serviceId) : undefined
              const professional = entry.professionalId ? referenceData.byId.professional.get(entry.professionalId) : undefined
              return (
                <li
                  className="flex flex-col gap-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] px-4 py-3 md:flex-row md:items-center"
                  key={entry.id}
                >
                  <div className="min-w-0 flex-1">
                    <div className="flex flex-wrap items-center gap-2">
                      <span className="font-bold text-[var(--color-ink)]">{entry.name}</span>
                      <Badge tone={STATUS_TONES[entry.status] ?? 'neutral'}>{STATUS_LABELS[entry.status] ?? entry.status}</Badge>
                      <span className="text-xs font-semibold text-[var(--color-ink-soft)]">
                        {entry.createdAt ? `inscrito em ${formatDateTime(entry.createdAt)}` : ''}
                      </span>
                    </div>
                    <p className="mt-1 truncate text-sm text-[var(--color-ink-soft)]">
                      {entry.phone}
                      {service ? ` · ${service.name}` : ''}
                      {professional ? ` · ${professional.name}` : ''}
                    </p>
                  </div>
                  <div className="flex shrink-0 flex-wrap items-center gap-2">
                    {(NEXT_STATUSES[entry.status] ?? []).map((nextStatus) => (
                      <Button
                        key={nextStatus}
                        onClick={() => void handleStatusChange(entry, nextStatus)}
                        size="sm"
                        variant={nextStatus === 'converted' ? 'primary' : 'secondary'}
                      >
                        {STATUS_LABELS[nextStatus] ?? nextStatus}
                      </Button>
                    ))}
                  </div>
                </li>
              )
            })}
          </ul>
        </Section>
      )}
    </PageContainer>
  )
}
