import { CreditCard, Download, Plus } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import { createPayment, downloadAdminCsv, fetchPayments, markPaymentPaid, refundPayment } from '../lib/adminApi'
import { formatDateTime, formatMoney } from '../lib/format'
import { Badge } from '../shared/components/Badge'
import { Button } from '../shared/components/Button'
import { EmptyState } from '../shared/components/EmptyState'
import { Field, Form, Input, Select } from '../shared/components/FormControls'
import { Modal } from '../shared/components/Modal'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { Payment } from '../types'

const PAYMENT_STATUS_LABELS: Record<string, string> = {
  pending: 'Pendente',
  paid: 'Pago',
  failed: 'Falhou',
  refunded: 'Estornado',
}

const PAYMENT_METHOD_LABELS: Record<string, string> = {
  pix: 'Pix',
  card: 'Cartão',
}

function paymentStatusTone(status: string): 'neutral' | 'success' | 'warning' | 'danger' {
  if (status === 'paid') {
    return 'success'
  }
  if (status === 'pending') {
    return 'warning'
  }
  if (status === 'refunded') {
    return 'neutral'
  }
  return 'danger'
}

export function PaymentsPage() {
  const { showToast } = useToast()
  const [payments, setPayments] = useState<Payment[]>([])
  const [loading, setLoading] = useState(true)
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [exporting, setExporting] = useState(false)

  function loadPayments() {
    setLoading(true)
    fetchPayments()
      .then(setPayments)
      .catch(() => {
        showToast({ message: 'Não foi possível carregar os pagamentos.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  useEffect(loadPayments, [])

  async function handleMarkPaid(payment: Payment) {
    try {
      await markPaymentPaid(payment.id)
      showToast({ message: 'Pagamento marcado como pago.', tone: 'success' })
      loadPayments()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    }
  }

  async function handleRefund(payment: Payment) {
    try {
      await refundPayment(payment.id)
      showToast({ message: 'Pagamento estornado.', tone: 'success' })
      loadPayments()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    }
  }

  async function handleExport() {
    setExporting(true)
    try {
      await downloadAdminCsv('/api/v1/admin/export/payments.csv', 'pagamentos.csv')
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    } finally {
      setExporting(false)
    }
  }

  return (
    <PageContainer>
      <Section
        actions={
          <>
            <Button disabled={exporting} icon={<Download className="h-4 w-4" />} loading={exporting} onClick={() => void handleExport()}>
              Exportar CSV
            </Button>
            <Button icon={<Plus className="h-4 w-4" />} onClick={() => setIsFormOpen(true)}>
              Registrar pagamento
            </Button>
          </>
        }
        subtitle="Pagamentos, depósitos e estornos registrados para os agendamentos."
        title="Pagamentos"
      >
        <div className="mt-2" />
      </Section>

      {loading ? (
        <Section>
          <Spinner />
        </Section>
      ) : payments.length === 0 ? (
        <Section>
          <EmptyState
            description="Os pagamentos registrados nos agendamentos aparecerão aqui."
            icon={<CreditCard className="h-10 w-10" />}
            title="Nenhum pagamento registrado"
          />
        </Section>
      ) : (
        <Section>
          <ul className="mt-2 grid gap-3">
            {payments.map((payment) => (
              <li
                className="flex flex-col gap-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] px-4 py-3 md:flex-row md:items-center"
                key={payment.id}
              >
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <span className="font-headline text-lg font-black text-[var(--color-ink)]">
                      {formatMoney(payment.amountCents)}
                    </span>
                    <Badge tone="neutral">{PAYMENT_METHOD_LABELS[payment.method] ?? payment.method}</Badge>
                    <Badge tone={paymentStatusTone(payment.status)}>{PAYMENT_STATUS_LABELS[payment.status] ?? payment.status}</Badge>
                    <span className="text-xs font-semibold text-[var(--color-ink-soft)]">
                      {payment.createdAt ? formatDateTime(payment.createdAt) : ''}
                    </span>
                  </div>
                  {payment.description ? (
                    <p className="mt-1 truncate text-sm text-[var(--color-ink-soft)]">{payment.description}</p>
                  ) : null}
                  {payment.provider ? (
                    <p className="mt-1 text-xs font-semibold text-[var(--color-ink-soft)]">Provedor: {payment.provider}</p>
                  ) : null}
                </div>
                <div className="flex shrink-0 flex-wrap items-center gap-2">
                  {payment.status === 'pending' ? (
                    <Button onClick={() => void handleMarkPaid(payment)} size="sm" variant="secondary">
                      Marcar como pago
                    </Button>
                  ) : null}
                  {payment.status === 'paid' ? (
                    <Button onClick={() => void handleRefund(payment)} size="sm" variant="danger">
                      Estornar
                    </Button>
                  ) : null}
                </div>
              </li>
            ))}
          </ul>
        </Section>
      )}

      {isFormOpen ? (
        <CreatePaymentModal
          onClose={() => setIsFormOpen(false)}
          onSaved={() => {
            setIsFormOpen(false)
            loadPayments()
          }}
        />
      ) : null}
    </PageContainer>
  )
}

function CreatePaymentModal({ onClose, onSaved }: { onClose: () => void; onSaved: () => void }) {
  const { showToast } = useToast()
  const [appointmentId, setAppointmentId] = useState('')
  const [amount, setAmount] = useState('')
  const [method, setMethod] = useState('pix')
  const [description, setDescription] = useState('')
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const amountCents = Math.round(Number(amount) * 100)
    if (!Number.isFinite(amountCents) || amountCents <= 0) {
      setError('Informe um valor válido.')
      return
    }

    setIsSubmitting(true)
    setError('')

    try {
      await createPayment({
        appointmentId: appointmentId || undefined,
        amountCents,
        method,
        description: description.trim() || undefined,
      })
      showToast({ message: 'Pagamento registrado.', tone: 'success' })
      onSaved()
    } catch (submitError) {
      setError((submitError as Error).message)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Modal
      footer={
        <>
          <Button disabled={isSubmitting} onClick={onClose} variant="secondary">
            Cancelar
          </Button>
          <Button loading={isSubmitting} type="submit" form="payment-form">
            Registrar
          </Button>
        </>
      }
      onClose={onClose}
      open
      size="md"
      subtitle="Registre um pagamento manualmente (depósito ou pagamento no local)."
      title="Registrar pagamento"
    >
      <Form id="payment-form" onSubmit={handleSubmit}>
        <Field label="Valor (R$)">
          <Input
            className="mt-2"
            min="0"
            onChange={(event) => setAmount(event.target.value)}
            step="0.01"
            type="number"
            value={amount}
          />
        </Field>
        <Field label="Método">
          <Select className="mt-2" onChange={(event) => setMethod(event.target.value)} value={method}>
            <option value="pix">Pix</option>
            <option value="card">Cartão</option>
          </Select>
        </Field>
        <Field label="ID do agendamento (opcional)">
          <Input className="mt-2" onChange={(event) => setAppointmentId(event.target.value)} value={appointmentId} />
        </Field>
        <Field label="Descrição (opcional)">
          <Input className="mt-2" onChange={(event) => setDescription(event.target.value)} value={description} />
        </Field>
        {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
      </Form>
    </Modal>
  )
}
