import { Plus, RefreshCcw } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import {
  createMemberPlan,
  createMembership,
  deleteMembership,
  fetchMemberPlans,
  fetchMemberships,
  setMembershipActive,
  updateMembership,
  updateMemberPlanStatus,
} from '../lib/adminApi'
import { formatDateTime, formatMoney } from '../lib/format'
import { useReferenceData } from '../shared/hooks/useReferenceData'
import { Badge } from '../shared/components/Badge'
import { Button } from '../shared/components/Button'
import { ConfirmModal } from '../shared/components/ConfirmModal'
import { EmptyState } from '../shared/components/EmptyState'
import { Field, Form, Input, Select, Textarea } from '../shared/components/FormControls'
import { Modal } from '../shared/components/Modal'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { MemberPlan, Membership } from '../types'

const PLAN_STATUS_LABELS: Record<string, string> = {
  active: 'Ativa',
  paused: 'Pausada',
  cancelled: 'Cancelada',
}

function planStatusTone(status: string): 'neutral' | 'success' | 'warning' | 'danger' {
  if (status === 'active') {
    return 'success'
  }
  if (status === 'paused') {
    return 'warning'
  }
  if (status === 'cancelled') {
    return 'danger'
  }
  return 'neutral'
}

export function MembershipsPage() {
  const { showToast } = useToast()
  const referenceData = useReferenceData()
  const [memberships, setMemberships] = useState<Membership[]>([])
  const [memberPlans, setMemberPlans] = useState<MemberPlan[]>([])
  const [loading, setLoading] = useState(true)
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [isPlanOpen, setIsPlanOpen] = useState(false)
  const [editingMembership, setEditingMembership] = useState<Membership | null>(null)
  const [membershipToDelete, setMembershipToDelete] = useState<Membership | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  function loadMemberships() {
    setLoading(true)
    Promise.allSettled([fetchMemberships(), fetchMemberPlans()])
      .then(([membershipsResult, plansResult]) => {
        setMemberships(membershipsResult.status === 'fulfilled' ? membershipsResult.value : [])
        setMemberPlans(plansResult.status === 'fulfilled' ? plansResult.value : [])
      })
      .catch(() => {
        showToast({ message: 'Não foi possível carregar as assinaturas.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  useEffect(loadMemberships, [])

  async function handleToggleActive(item: Membership) {
    try {
      await setMembershipActive(item.id, !item.active)
      showToast({ message: item.active ? 'Assinatura desativada.' : 'Assinatura ativada.', tone: 'success' })
      loadMemberships()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    }
  }

  async function handlePlanStatusChange(plan: MemberPlan, status: string) {
    try {
      await updateMemberPlanStatus(plan.id, status)
      showToast({ message: 'Plano atualizado.', tone: 'success' })
      loadMemberships()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    }
  }

  async function handleDelete() {
    if (!membershipToDelete) {
      return
    }
    setIsDeleting(true)
    try {
      await deleteMembership(membershipToDelete.id)
      showToast({ message: 'Assinatura removida.', tone: 'success' })
      setMembershipToDelete(null)
      loadMemberships()
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
          <div className="flex flex-wrap gap-2">
            <Button icon={<RefreshCcw className="h-4 w-4" />} onClick={() => setIsPlanOpen(true)} variant="secondary">
              Assinar cliente
            </Button>
            <Button
              icon={<Plus className="h-4 w-4" />}
              onClick={() => {
                setEditingMembership(null)
                setIsFormOpen(true)
              }}
            >
              Nova assinatura
            </Button>
          </div>
        }
        subtitle="Assinaturas recorrentes e planos ativos dos clientes."
        title="Assinaturas"
      >
        <div className="mt-2" />
      </Section>

      {loading ? (
        <Section>
          <Spinner />
        </Section>
      ) : memberships.length === 0 ? (
        <Section>
          <EmptyState
            description="Crie assinaturas para oferecer planos recorrentes aos clientes."
            icon={<RefreshCcw className="h-10 w-10" />}
            title="Nenhuma assinatura cadastrada"
          />
        </Section>
      ) : (
        <Section>
          <ul className="mt-2 grid gap-3 md:grid-cols-2 xl:grid-cols-3">
            {memberships.map((item) => (
              <li
                className="flex flex-col gap-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4"
                key={item.id}
              >
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <p className="truncate font-bold text-[var(--color-ink)]">{item.name}</p>
                    <p className="mt-0.5 text-sm font-semibold text-[var(--color-ink-soft)]">
                      Cobrança a cada {item.billingPeriodDays} dias
                    </p>
                  </div>
                  <Badge tone={item.active ? 'success' : 'danger'}>{item.active ? 'Ativo' : 'Inativo'}</Badge>
                </div>
                {item.description ? (
                  <p className="line-clamp-2 text-sm text-[var(--color-ink-soft)]">{item.description}</p>
                ) : null}
                <p className="font-headline text-lg font-black text-[var(--color-ink)]">
                  {formatMoney(item.priceCents)}/período
                </p>
                <div className="mt-auto flex flex-wrap items-center gap-2">
                  <Button
                    onClick={() => {
                      setEditingMembership(item)
                      setIsFormOpen(true)
                    }}
                    size="sm"
                    variant="secondary"
                  >
                    Editar
                  </Button>
                  <Button onClick={() => void handleToggleActive(item)} size="sm" variant="secondary">
                    {item.active ? 'Desativar' : 'Ativar'}
                  </Button>
                  <Button onClick={() => setMembershipToDelete(item)} size="sm" variant="danger">
                    Remover
                  </Button>
                </div>
              </li>
            ))}
          </ul>
        </Section>
      )}

      <Section subtitle="Planos de assinatura ativos por cliente." title="Planos de clientes">
        {memberPlans.length === 0 ? (
          <p className="mt-3 text-sm font-semibold text-[var(--color-ink-soft)]">
            Nenhum cliente assinando. Use "Assinar cliente" para iniciar um plano.
          </p>
        ) : (
          <ul className="mt-4 grid gap-2">
            {memberPlans.map((plan) => (
              <li
                className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] px-4 py-3"
                key={plan.id}
              >
                <div className="min-w-0">
                  <p className="truncate font-bold text-[var(--color-ink)]">{plan.customerName ?? plan.customerId}</p>
                  <p className="truncate text-sm text-[var(--color-ink-soft)]">
                    {plan.membershipName ?? plan.membershipId}
                    {plan.startedAt ? ` · desde ${formatDateTime(plan.startedAt)}` : ''}
                  </p>
                </div>
                <div className="flex shrink-0 items-center gap-2">
                  <Badge tone={planStatusTone(plan.status)}>{PLAN_STATUS_LABELS[plan.status] ?? plan.status}</Badge>
                  <Select
                    className="h-9 w-auto min-w-32 text-xs"
                    onChange={(event) => void handlePlanStatusChange(plan, event.target.value)}
                    value=""
                  >
                    <option value="">Alterar...</option>
                    {Object.entries(PLAN_STATUS_LABELS).map(([value, label]) => (
                      <option key={value} value={value}>
                        {label}
                      </option>
                    ))}
                  </Select>
                </div>
              </li>
            ))}
          </ul>
        )}
      </Section>

      {isFormOpen ? (
        <MembershipFormModal
          item={editingMembership}
          onClose={() => setIsFormOpen(false)}
          onSaved={() => {
            setIsFormOpen(false)
            loadMemberships()
          }}
        />
      ) : null}

      {isPlanOpen ? (
        <CreatePlanModal
          customers={referenceData.customers}
          memberships={memberships}
          onClose={() => setIsPlanOpen(false)}
          onSaved={() => {
            setIsPlanOpen(false)
            loadMemberships()
          }}
        />
      ) : null}

      <ConfirmModal
        confirmLabel="Remover"
        loading={isDeleting}
        message={`A assinatura "${membershipToDelete?.name ?? ''}" será removida.`}
        onCancel={() => setMembershipToDelete(null)}
        onConfirm={() => void handleDelete()}
        open={Boolean(membershipToDelete)}
        title="Remover assinatura?"
      />
    </PageContainer>
  )
}

function MembershipFormModal({
  item,
  onClose,
  onSaved,
}: {
  item: Membership | null
  onClose: () => void
  onSaved: () => void
}) {
  const { showToast } = useToast()
  const [name, setName] = useState(item?.name ?? '')
  const [description, setDescription] = useState(item?.description ?? '')
  const [price, setPrice] = useState(item ? String(item.priceCents / 100) : '')
  const [billingPeriodDays, setBillingPeriodDays] = useState(item?.billingPeriodDays ?? 30)
  const [active, setActive] = useState(item?.active ?? true)
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (!name.trim()) {
      setError('Informe o nome da assinatura.')
      return
    }
    const priceCents = Math.round(Number(price) * 100)
    if (!Number.isFinite(priceCents) || priceCents <= 0) {
      setError('Informe um preço válido.')
      return
    }
    if (billingPeriodDays < 1) {
      setError('Informe o período de cobrança.')
      return
    }

    setIsSubmitting(true)
    setError('')

    try {
      const input = {
        name: name.trim(),
        description: description.trim(),
        priceCents,
        billingPeriodDays,
        active,
      }
      if (item) {
        await updateMembership(item.id, input)
        showToast({ message: 'Assinatura atualizada.', tone: 'success' })
      } else {
        await createMembership(input)
        showToast({ message: 'Assinatura criada.', tone: 'success' })
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
          <Button loading={isSubmitting} type="submit" form="membership-form">
            Salvar
          </Button>
        </>
      }
      onClose={onClose}
      open
      size="md"
      subtitle={item ? 'Atualize os dados da assinatura.' : 'Cadastre uma nova assinatura.'}
      title={item ? 'Editar assinatura' : 'Nova assinatura'}
    >
      <Form id="membership-form" onSubmit={handleSubmit}>
        <Field label="Nome">
          <Input className="mt-2" onChange={(event) => setName(event.target.value)} value={name} />
        </Field>
        <Field label="Descrição">
          <Textarea className="mt-2" onChange={(event) => setDescription(event.target.value)} value={description} />
        </Field>
        <div className="grid gap-4 md:grid-cols-2">
          <Field label="Preço (R$)">
            <Input
              className="mt-2"
              min="0"
              onChange={(event) => setPrice(event.target.value)}
              step="0.01"
              type="number"
              value={price}
            />
          </Field>
          <Field label="Período de cobrança (dias)">
            <Input
              className="mt-2"
              min={1}
              onChange={(event) => setBillingPeriodDays(Number(event.target.value))}
              type="number"
              value={billingPeriodDays}
            />
          </Field>
        </div>
        <Field label=" ">
          <input
            checked={active}
            className="h-4 w-4 rounded border-[var(--color-line)] text-[var(--color-primary)]"
            onChange={(event) => setActive(event.target.checked)}
            type="checkbox"
          />
          <span className="ml-2 text-sm font-bold text-[var(--color-ink)]">Assinatura ativa (pode ser contratada)</span>
        </Field>
        {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
      </Form>
    </Modal>
  )
}

function CreatePlanModal({
  memberships,
  customers,
  onClose,
  onSaved,
}: {
  memberships: Membership[]
  customers: { id: string; name: string }[]
  onClose: () => void
  onSaved: () => void
}) {
  const { showToast } = useToast()
  const [membershipId, setMembershipId] = useState(memberships[0]?.id ?? '')
  const [customerId, setCustomerId] = useState(customers[0]?.id ?? '')
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (!membershipId || !customerId) {
      setError('Selecione a assinatura e o cliente.')
      return
    }

    setIsSubmitting(true)
    setError('')

    try {
      await createMemberPlan({ membershipId, customerId })
      showToast({ message: 'Plano iniciado para o cliente.', tone: 'success' })
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
          <Button loading={isSubmitting} type="submit" form="plan-form">
            Assinar
          </Button>
        </>
      }
      onClose={onClose}
      open
      size="md"
      subtitle="Inicie um plano de assinatura para um cliente."
      title="Assinar cliente"
    >
      <Form id="plan-form" onSubmit={handleSubmit}>
        <Field label="Assinatura">
          <Select className="mt-2" onChange={(event) => setMembershipId(event.target.value)} value={membershipId}>
            {memberships.map((item) => (
              <option key={item.id} value={item.id}>
                {item.name} · {formatMoney(item.priceCents)}/{item.billingPeriodDays} dias
              </option>
            ))}
          </Select>
        </Field>
        <Field label="Cliente">
          <Select className="mt-2" onChange={(event) => setCustomerId(event.target.value)} value={customerId}>
            {customers.map((customer) => (
              <option key={customer.id} value={customer.id}>
                {customer.name}
              </option>
            ))}
          </Select>
        </Field>
        {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
      </Form>
    </Modal>
  )
}
