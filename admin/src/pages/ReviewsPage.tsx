import { Copy, Download, Star, MessageSquare } from 'lucide-react'
import { useEffect, useState } from 'react'

import { downloadAdminCsv, fetchReviews } from '../lib/adminApi'
import { formatDateTime } from '../lib/format'
import { Badge } from '../shared/components/Badge'
import { Button } from '../shared/components/Button'
import { EmptyState } from '../shared/components/EmptyState'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Select } from '../shared/components/FormControls'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { Review } from '../types'

const REVIEW_STATUS_LABELS: Record<string, string> = {
  invited: 'Convite enviado',
  submitted: 'Respondida',
}

const bookingOrigin = import.meta.env.VITE_BOOKING_ORIGIN?.trim() || (typeof window !== 'undefined' ? `${window.location.origin}/booking` : '')

function reviewStatusTone(status: string): 'neutral' | 'success' | 'warning' | 'danger' {
  if (status === 'submitted') {
    return 'success'
  }
  return 'warning'
}

function reviewLink(review: Review): string {
  const slug = review.tenantSlug || review.tenantId
  return `${bookingOrigin}/#/review/${encodeURIComponent(slug)}/${encodeURIComponent(review.id)}`
}

export function ReviewsPage() {
  const { showToast } = useToast()
  const [reviews, setReviews] = useState<Review[]>([])
  const [status, setStatus] = useState('')
  const [loading, setLoading] = useState(true)
  const [exporting, setExporting] = useState(false)

  function loadReviews() {
    setLoading(true)
    fetchReviews(status)
      .then(setReviews)
      .catch(() => {
        showToast({ message: 'Não foi possível carregar as avaliações.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  useEffect(loadReviews, [status])

  async function handleExport() {
    setExporting(true)
    try {
      await downloadAdminCsv('/api/v1/admin/export/reviews.csv', 'avaliacoes.csv')
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    } finally {
      setExporting(false)
    }
  }

  async function handleCopyLink(review: Review) {
    try {
      await navigator.clipboard.writeText(reviewLink(review))
      showToast({ message: 'Link da avaliação copiado.', tone: 'success' })
    } catch {
      showToast({ message: 'Não foi possível copiar o link.', tone: 'error' })
    }
  }

  return (
    <PageContainer>
      <Section
        actions={
          <>
            <Select
              aria-label="Filtrar por status"
              className="w-44"
              onChange={(event) => setStatus(event.target.value)}
              value={status}
            >
              <option value="">Todos os status</option>
              <option value="invited">Convites enviados</option>
              <option value="submitted">Respondidas</option>
            </Select>
            <Button disabled={exporting} icon={<Download className="h-4 w-4" />} loading={exporting} onClick={() => void handleExport()}>
              Exportar CSV
            </Button>
          </>
        }
        subtitle="Avaliações enviadas automaticamente após a conclusão dos agendamentos."
        title="Avaliações"
      >
        <div className="mt-2" />
      </Section>

      {loading ? (
        <Section>
          <Spinner />
        </Section>
      ) : reviews.length === 0 ? (
        <Section>
          <EmptyState
            description="Quando um agendamento for concluído, um convite de avaliação é enviado ao cliente."
            icon={<Star className="h-10 w-10" />}
            title="Nenhuma avaliação encontrada"
          />
        </Section>
      ) : (
        <Section>
          <ul className="mt-2 grid gap-3">
            {reviews.map((review) => (
              <li
                className="flex flex-col gap-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] px-4 py-3 md:flex-row md:items-center"
                key={review.id}
              >
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <span className="font-headline text-lg font-black text-[var(--color-ink)]">
                      {review.customerName || review.customerId || 'Cliente'}
                    </span>
                    <Badge tone={reviewStatusTone(review.status)}>{REVIEW_STATUS_LABELS[review.status] ?? review.status}</Badge>
                    <span className="text-xs font-semibold text-[var(--color-ink-soft)]">
                      {formatDateTime(review.invitedAt)}
                    </span>
                  </div>
                  {review.status === 'submitted' ? (
                    <div className="mt-1 flex flex-wrap items-center gap-3">
                      <span className="flex items-center gap-1">
                        {[1, 2, 3, 4, 5].map((value) => (
                          <Star
                            className={`h-4 w-4 ${value <= review.rating ? 'fill-amber-400 text-amber-400' : 'text-stone-300'}`}
                            key={value}
                          />
                        ))}
                      </span>
                      {review.comment ? (
                        <span className="flex items-center gap-1 truncate text-sm text-[var(--color-ink-soft)]">
                          <MessageSquare className="h-4 w-4 shrink-0" />
                          {review.comment}
                        </span>
                      ) : null}
                    </div>
                  ) : (
                    <p className="mt-1 text-sm text-[var(--color-ink-soft)]">Aguardando resposta do cliente.</p>
                  )}
                </div>
                {review.status === 'invited' ? (
                  <div className="flex shrink-0 items-center gap-2">
                    <Button
                      icon={<Copy className="h-4 w-4" />}
                      onClick={() => void handleCopyLink(review)}
                      size="sm"
                      variant="secondary"
                    >
                      Copiar link
                    </Button>
                  </div>
                ) : null}
              </li>
            ))}
          </ul>
        </Section>
      )}
    </PageContainer>
  )
}
