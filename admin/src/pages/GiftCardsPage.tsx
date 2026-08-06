import { Gift, Plus } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import { createGiftCard, deleteGiftCard, fetchGiftCards, updateGiftCard } from '../lib/adminApi'
import { formatMoney } from '../lib/format'
import { useReferenceData } from '../shared/hooks/useReferenceData'
import { Badge } from '../shared/components/Badge'
import { Button } from '../shared/components/Button'
import { ConfirmModal } from '../shared/components/ConfirmModal'
import { EmptyState } from '../shared/components/EmptyState'
import { Checkbox, Field, Form, Input, Select } from '../shared/components/FormControls'
import { Modal } from '../shared/components/Modal'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { GiftCard } from '../types'

function statusTone(status: string): 'neutral' | 'success' | 'danger' {
  if (status === 'active') {
    return 'success'
  }
  if (status === 'redeemed') {
    return 'neutral'
  }
  return 'danger'
}

export function GiftCardsPage() {
  const { showToast } = useToast()
  const referenceData = useReferenceData()
  const [giftCards, setGiftCards] = useState<GiftCard[]>([])
  const [loading, setLoading] = useState(true)
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [editingGiftCard, setEditingGiftCard] = useState<GiftCard | null>(null)
  const [giftCardToDelete, setGiftCardToDelete] = useState<GiftCard | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  function loadGiftCards() {
    setLoading(true)
    fetchGiftCards()
      .then(setGiftCards)
      .catch(() => {
        showToast({ message: 'Não foi possível carregar os gift cards.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  useEffect(loadGiftCards, [])

  async function handleDelete() {
    if (!giftCardToDelete) {
      return
    }
    setIsDeleting(true)
    try {
      await deleteGiftCard(giftCardToDelete.id)
      showToast({ message: 'Gift card removido.', tone: 'success' })
      setGiftCardToDelete(null)
      loadGiftCards()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    } finally {
      setIsDeleting(false)
    }
  }

  return (
    <PageContainer>
      <Section
        actions={
          <Button
            icon={<Plus className="h-4 w-4" />}
            onClick={() => {
              setEditingGiftCard(null)
              setIsFormOpen(true)
            }}
          >
            Novo gift card
          </Button>
        }
        subtitle="Cartões-presente emitidos e seus saldos."
        title="Gift cards"
      >
        <div className="mt-2" />
      </Section>

      {loading ? (
        <Section>
          <Spinner />
        </Section>
      ) : giftCards.length === 0 ? (
        <Section>
          <EmptyState
            description="Emita gift cards para vender créditos em forma de presente."
            icon={<Gift className="h-10 w-10" />}
            title="Nenhum gift card emitido"
          />
        </Section>
      ) : (
        <Section>
          <ul className="mt-2 grid gap-3 md:grid-cols-2 xl:grid-cols-3">
            {giftCards.map((card) => (
              <li
                className="flex flex-col gap-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4"
                key={card.id}
              >
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <p className="font-headline text-lg font-black tracking-widest text-[var(--color-ink)]">{card.code}</p>
                    <p className="mt-0.5 truncate text-sm font-semibold text-[var(--color-ink-soft)]">
                      {card.customerName ?? 'Sem cliente vinculado'}
                    </p>
                  </div>
                  <Badge tone={statusTone(card.status)}>{card.status}</Badge>
                </div>
                <div className="flex items-baseline gap-3">
                  <span className="font-headline text-xl font-black text-[var(--color-ink)]">{formatMoney(card.balanceCents)}</span>
                  <span className="text-sm font-semibold text-[var(--color-ink-soft)]">saldo de {formatMoney(card.valueCents)}</span>
                </div>
                <div className="mt-auto flex flex-wrap items-center gap-2">
                  <Button
                    onClick={() => {
                      setEditingGiftCard(card)
                      setIsFormOpen(true)
                    }}
                    size="sm"
                    variant="secondary"
                  >
                    Editar
                  </Button>
                  <Button onClick={() => setGiftCardToDelete(card)} size="sm" variant="danger">
                    Remover
                  </Button>
                </div>
              </li>
            ))}
          </ul>
        </Section>
      )}

      {isFormOpen ? (
        <GiftCardFormModal
          customers={referenceData.customers}
          giftCard={editingGiftCard}
          onClose={() => setIsFormOpen(false)}
          onSaved={() => {
            setIsFormOpen(false)
            loadGiftCards()
          }}
        />
      ) : null}

      <ConfirmModal
        confirmLabel="Remover"
        loading={isDeleting}
        message={`O gift card "${giftCardToDelete?.code ?? ''}" será removido.`}
        onCancel={() => setGiftCardToDelete(null)}
        onConfirm={() => void handleDelete()}
        open={Boolean(giftCardToDelete)}
        title="Remover gift card?"
      />
    </PageContainer>
  )
}

function GiftCardFormModal({
  giftCard,
  customers,
  onClose,
  onSaved,
}: {
  giftCard: GiftCard | null
  customers: { id: string; name: string }[]
  onClose: () => void
  onSaved: () => void
}) {
  const { showToast } = useToast()
  const [code, setCode] = useState(giftCard?.code ?? '')
  const [value, setValue] = useState(giftCard ? String(giftCard.valueCents / 100) : '')
  const [balance, setBalance] = useState(giftCard ? String(giftCard.balanceCents / 100) : '')
  const [customerId, setCustomerId] = useState(giftCard?.customerId ?? '')
  const [expiresAt, setExpiresAt] = useState(giftCard?.expiresAt ? giftCard.expiresAt.slice(0, 10) : '')
  const [active, setActive] = useState(giftCard?.status === 'active')
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const valueCents = Math.round(Number(value) * 100)
    if (!code.trim()) {
      setError('Informe o código do gift card.')
      return
    }
    if (!Number.isFinite(valueCents) || valueCents <= 0) {
      setError('Informe um valor válido.')
      return
    }

    setIsSubmitting(true)
    setError('')

    try {
      const input = {
        code: code.trim(),
        valueCents,
        balanceCents: giftCard ? Math.round(Number(balance) * 100) : valueCents,
        customerId: customerId || undefined,
        expiresAt: expiresAt ? `${expiresAt}T23:59:59` : undefined,
        status: active ? 'active' : 'inactive',
      }
      if (giftCard) {
        await updateGiftCard(giftCard.id, input)
        showToast({ message: 'Gift card atualizado.', tone: 'success' })
      } else {
        await createGiftCard(input)
        showToast({ message: 'Gift card criado.', tone: 'success' })
      }
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
          <Button loading={isSubmitting} type="submit" form="giftcard-form">
            Salvar
          </Button>
        </>
      }
      onClose={onClose}
      open
      size="md"
      subtitle={giftCard ? 'Atualize os dados do gift card.' : 'Emita um novo gift card.'}
      title={giftCard ? 'Editar gift card' : 'Novo gift card'}
    >
      <Form id="giftcard-form" onSubmit={handleSubmit}>
        <Field label="Código">
          <Input className="mt-2" onChange={(event) => setCode(event.target.value)} value={code} />
        </Field>
        <div className="grid gap-4 md:grid-cols-2">
          <Field label="Valor (R$)">
            <Input
              className="mt-2"
              min="0"
              onChange={(event) => setValue(event.target.value)}
              step="0.01"
              type="number"
              value={value}
            />
          </Field>
          <Field label="Saldo (R$)">
            <Input
              className="mt-2"
              min="0"
              onChange={(event) => setBalance(event.target.value)}
              step="0.01"
              type="number"
              value={balance}
            />
          </Field>
        </div>
        <Field label="Cliente (opcional)">
          <Select className="mt-2" onChange={(event) => setCustomerId(event.target.value)} value={customerId}>
            <option value="">Sem cliente</option>
            {customers.map((customer) => (
              <option key={customer.id} value={customer.id}>
                {customer.name}
              </option>
            ))}
          </Select>
        </Field>
        <Field label="Validade (opcional)">
          <Input className="mt-2" onChange={(event) => setExpiresAt(event.target.value)} type="date" value={expiresAt} />
        </Field>
        <Field label=" ">
          <Checkbox checked={active} onChange={(event) => setActive(event.target.checked)} />
          <span className="ml-2 text-sm font-bold text-[var(--color-ink)]">Gift card ativo (pode ser usado)</span>
        </Field>
        {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
      </Form>
    </Modal>
  )
}
