import {
  AlertTriangle,
  Building2,
  Calendar,
  CheckCircle,
  Clock,
  Layers,
  Lock,
  Plus,
  ShieldCheck,
  SlidersHorizontal,
  Sparkles,
  Unlock,
} from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import {
  createTenant,
  deleteTenant,
  fetchTenants,
  setTenantActive,
  updateTenant,
  updateTenantSubscription,
} from '../lib/adminApi'
import { formatDateTime } from '../lib/format'
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
import type { Tenant, TenantPlan, TenantSubscriptionStatus, UpdateTenantSubscriptionInput } from '../types'

const SLUG_PATTERN = /^[a-z0-9](?:[a-z0-9-]{0,60}[a-z0-9])?$/

export const AVAILABLE_MODULES = [
  { id: 'agenda', label: 'Agenda & Visão Diária', desc: 'Visualização da grade de horários' },
  { id: 'appointments', label: 'Agendamentos', desc: 'Gestão de horários e lista de agendados' },
  { id: 'customers', label: 'Clientes & CRM', desc: 'Histórico, contatos e dados de clientes' },
  { id: 'services', label: 'Serviços & Preços', desc: 'Catálogo de procedimentos e durações' },
  { id: 'professionals', label: 'Profissionais & Equipe', desc: 'Cadastros da equipe e horários' },
  { id: 'finance', label: 'Financeiro & DRE', desc: 'Fluxo de caixa, receitas e despesas' },
  { id: 'pos_cashregister', label: 'Caixa Diário & Frente', desc: 'Gaveta de troco, sangria, suprimento e fechar comanda' },
  { id: 'commissions', label: 'Comissões & Repasses', desc: 'Cálculo automático e fechamento de repasses' },
  { id: 'inventory', label: 'Estoque & Insumos', desc: 'Produtos para venda e insumos por serviço' },
  { id: 'commerce', label: 'Planos, Pacotes & Gift Cards', desc: 'Clube de membros e pacotes de sessões' },
  { id: 'reviews', label: 'Avaliações de Clientes', desc: 'Coleta de notas e feedbacks' },
  { id: 'notifications', label: 'Lembretes & Notificações', desc: 'Disparos automáticos por e-mail/SMS/WhatsApp' },
  { id: 'widget', label: 'Widget Integrado', desc: 'Script para incorporar agendamento no site do cliente' },
  { id: 'calendar_sync', label: 'Sincronização Google Calendar', desc: 'Integração bidirecional de calendários' },
]

export const PLAN_PRESETS: Record<TenantPlan, string[]> = {
  starter: ['agenda', 'appointments', 'customers', 'services', 'professionals', 'notifications'],
  pro: [
    'agenda',
    'appointments',
    'customers',
    'services',
    'professionals',
    'finance',
    'pos_cashregister',
    'commissions',
    'inventory',
    'reviews',
    'notifications',
  ],
  enterprise: AVAILABLE_MODULES.map((m) => m.id),
  trial: AVAILABLE_MODULES.map((m) => m.id),
  custom: [],
}

export function TenantsPage() {
  const { showToast } = useToast()
  const [tenants, setTenants] = useState<Tenant[]>([])
  const [loading, setLoading] = useState(true)

  const [isFormOpen, setIsFormOpen] = useState(false)
  const [editingTenant, setEditingTenant] = useState<Tenant | null>(null)

  const [managingSubscriptionTenant, setManagingSubscriptionTenant] = useState<Tenant | null>(null)

  const [tenantToDelete, setTenantToDelete] = useState<Tenant | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  function loadTenants() {
    setLoading(true)
    fetchTenants()
      .then(setTenants)
      .catch(() => {
        showToast({ message: 'Não foi possível carregar os estabelecimentos.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  useEffect(loadTenants, [])

  async function handleToggleActive(tenant: Tenant) {
    try {
      await setTenantActive(tenant.id, !tenant.active)
      showToast({
        message: tenant.active ? 'Estabelecimento desativado.' : 'Estabelecimento ativado.',
        tone: 'success',
      })
      loadTenants()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    }
  }

  async function handleDelete() {
    if (!tenantToDelete) return
    setIsDeleting(true)
    try {
      await deleteTenant(tenantToDelete.id)
      showToast({ message: 'Estabelecimento removido.', tone: 'success' })
      setTenantToDelete(null)
      loadTenants()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    } finally {
      setIsDeleting(false)
    }
  }

  // KPIs
  const totalTenants = tenants.length
  const activeSubs = tenants.filter((t) => t.subscriptionStatus === 'active').length
  const trialingSubs = tenants.filter((t) => t.subscriptionStatus === 'trialing').length
  const pastDueSubs = tenants.filter((t) => t.subscriptionStatus === 'past_due').length

  function getPlanBadge(plan: TenantPlan) {
    switch (plan) {
      case 'enterprise':
        return <Badge tone="accent">Enterprise (Ilimitado)</Badge>
      case 'pro':
        return <Badge tone="success">Plano PRO</Badge>
      case 'starter':
        return <Badge tone="neutral">Plano Starter</Badge>
      case 'trial':
        return <Badge tone="warning">Trial (14 dias)</Badge>
      default:
        return <Badge tone="neutral">Personalizado</Badge>
    }
  }

  function getStatusBadge(status: TenantSubscriptionStatus, trialEndsAt?: string) {
    if (status === 'trialing') {
      let daysRemaining = ''
      if (trialEndsAt) {
        const diffDays = Math.ceil(
          (new Date(trialEndsAt).getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24)
        )
        daysRemaining = diffDays > 0 ? ` (${diffDays}d restantes)` : ' (Expirado)'
      }
      return <Badge tone="warning">Em Teste{daysRemaining}</Badge>
    }
    if (status === 'active') {
      return <Badge tone="success">Assinatura Ativa</Badge>
    }
    if (status === 'past_due') {
      return <Badge tone="danger">Em Atraso</Badge>
    }
    return <Badge tone="neutral">Cancelado</Badge>
  }

  return (
    <PageContainer>
      <Section
        actions={
          <Button
            icon={<Plus className="h-4 w-4" />}
            onClick={() => {
              setEditingTenant(null)
              setIsFormOpen(true)
            }}
          >
            Novo estabelecimento
          </Button>
        }
        subtitle="Painel de controle Master Admin: gerencie todos os estabelecimentos, assinaturas e permissões de módulos."
        title="Admin Master: Estabelecimentos & Assinaturas"
      >
        {/* KPI Summary Cards */}
        <div className="mt-2 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <div className="rounded-xl border border-border bg-card p-4 shadow-sm">
            <span className="text-xs font-semibold text-muted-foreground">Total de Estabelecimentos</span>
            <p className="mt-2 text-2xl font-bold text-foreground">{totalTenants}</p>
          </div>
          <div className="rounded-xl border border-border bg-card p-4 shadow-sm">
            <span className="text-xs font-semibold text-emerald-600">Assinaturas Ativas (Pagas)</span>
            <p className="mt-2 text-2xl font-bold text-emerald-600">{activeSubs}</p>
          </div>
          <div className="rounded-xl border border-border bg-card p-4 shadow-sm">
            <span className="text-xs font-semibold text-amber-600">Em Período de Testes (Trial)</span>
            <p className="mt-2 text-2xl font-bold text-amber-600">{trialingSubs}</p>
          </div>
          <div className="rounded-xl border border-border bg-card p-4 shadow-sm">
            <span className="text-xs font-semibold text-rose-600">Em Atraso / Pendentes</span>
            <p className="mt-2 text-2xl font-bold text-rose-600">{pastDueSubs}</p>
          </div>
        </div>
      </Section>

      {loading ? (
        <Section>
          <Spinner />
        </Section>
      ) : tenants.length === 0 ? (
        <Section>
          <EmptyState
            icon={<Building2 className="h-10 w-10" />}
            description="Nenhum estabelecimento cadastrado. Crie um estabelecimento para começar a operar."
            title="Sem estabelecimentos"
          />
        </Section>
      ) : (
        <Section title="Lista de Estabelecimentos (Tenants)">
          <div className="overflow-x-auto">
            <table className="w-full border-collapse text-left text-sm">
              <thead>
                <tr className="border-b border-border text-xs uppercase text-muted-foreground">
                  <th className="py-3 px-4">Estabelecimento</th>
                  <th className="py-3 px-4">Plano</th>
                  <th className="py-3 px-4">Status da Assinatura</th>
                  <th className="py-3 px-4">Módulos Habilitados</th>
                  <th className="py-3 px-4">Limite Profissionais</th>
                  <th className="py-3 px-4">Status</th>
                  <th className="py-3 px-4 text-right">Ações Master</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {tenants.map((tenant) => (
                  <tr key={tenant.id} className="hover:bg-muted/50 transition-colors">
                    <td className="py-3 px-4">
                      <div className="font-bold text-foreground">{tenant.name}</div>
                      <div className="text-xs text-muted-foreground">/{tenant.slug}</div>
                    </td>
                    <td className="py-3 px-4">{getPlanBadge(tenant.plan)}</td>
                    <td className="py-3 px-4">
                      {getStatusBadge(tenant.subscriptionStatus, tenant.trialEndsAt)}
                    </td>
                    <td className="py-3 px-4">
                      <span className="inline-flex items-center gap-1 font-semibold text-xs text-muted-foreground">
                        <Layers className="h-3.5 w-3.5 text-primary" />
                        {tenant.activeModules?.length || 0} / {AVAILABLE_MODULES.length} módulos
                      </span>
                    </td>
                    <td className="py-3 px-4 text-xs font-semibold text-foreground">
                      {tenant.maxProfessionals >= 999 ? 'Ilimitado' : `Até ${tenant.maxProfessionals}`}
                    </td>
                    <td className="py-3 px-4">
                      <Badge tone={tenant.active ? 'success' : 'danger'}>
                        {tenant.active ? 'Ativo' : 'Inativo'}
                      </Badge>
                    </td>
                    <td className="py-3 px-4 text-right">
                      <div className="flex items-center justify-end gap-1.5">
                        <Button
                          size="sm"
                          className="bg-primary/10 text-primary hover:bg-primary/20 border-0"
                          icon={<SlidersHorizontal className="h-3.5 w-3.5" />}
                          onClick={() => setManagingSubscriptionTenant(tenant)}
                        >
                          Módulos & Plano
                        </Button>
                        <Button
                          size="sm"
                          variant="secondary"
                          onClick={() => {
                            setEditingTenant(tenant)
                            setIsFormOpen(true)
                          }}
                        >
                          Editar
                        </Button>
                        <Button
                          size="sm"
                          variant="secondary"
                          onClick={() => void handleToggleActive(tenant)}
                        >
                          {tenant.active ? 'Desativar' : 'Ativar'}
                        </Button>
                        <Button
                          disabled={tenant.id === 'default'}
                          onClick={() => setTenantToDelete(tenant)}
                          size="sm"
                          title={
                            tenant.id === 'default'
                              ? 'O estabelecimento padrão não pode ser removido.'
                              : 'Remover'
                          }
                          variant="danger"
                        >
                          Remover
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Section>
      )}

      {/* Modal: Gerenciar Módulos & Assinatura (Master Admin) */}
      {managingSubscriptionTenant && (
        <TenantSubscriptionModal
          tenant={managingSubscriptionTenant}
          onClose={() => setManagingSubscriptionTenant(null)}
          onSaved={() => {
            setManagingSubscriptionTenant(null)
            loadTenants()
          }}
        />
      )}

      {/* Modal: Criar / Editar Básico */}
      {isFormOpen ? (
        <TenantFormModal
          onClose={() => setIsFormOpen(false)}
          onSaved={() => {
            setIsFormOpen(false)
            loadTenants()
          }}
          tenant={editingTenant}
        />
      ) : null}

      <ConfirmModal
        confirmLabel="Remover"
        loading={isDeleting}
        message={`O estabelecimento "${tenantToDelete?.name ?? ''}" e todos os dados vinculados a ele serão removidos permanentemente.`}
        onCancel={() => setTenantToDelete(null)}
        onConfirm={() => void handleDelete()}
        open={Boolean(tenantToDelete)}
        title="Remover estabelecimento?"
      />
    </PageContainer>
  )
}

function TenantSubscriptionModal({
  tenant,
  onClose,
  onSaved,
}: {
  tenant: Tenant
  onClose: () => void
  onSaved: () => void
}) {
  const { showToast } = useToast()
  const [plan, setPlan] = useState<TenantPlan>(tenant.plan || 'trial')
  const [status, setStatus] = useState<TenantSubscriptionStatus>(tenant.subscriptionStatus || 'trialing')
  const [maxProfessionals, setMaxProfessionals] = useState(tenant.maxProfessionals || 5)
  const [activeModules, setActiveModules] = useState<string[]>(
    tenant.activeModules?.length ? tenant.activeModules : AVAILABLE_MODULES.map((m) => m.id)
  )
  const [trialEndsDate, setTrialEndsDate] = useState(() => {
    if (tenant.trialEndsAt) {
      return new Date(tenant.trialEndsAt).toISOString().slice(0, 10)
    }
    const d = new Date()
    d.setDate(d.getDate() + 14)
    return d.toISOString().slice(0, 10)
  })
  const [isSubmitting, setIsSubmitting] = useState(false)

  function handlePlanChange(newPlan: TenantPlan) {
    setPlan(newPlan)
    if (newPlan !== 'custom') {
      const preset = PLAN_PRESETS[newPlan]
      if (preset && preset.length > 0) {
        setActiveModules(preset)
      }
      if (newPlan === 'starter') setMaxProfessionals(2)
      if (newPlan === 'pro') setMaxProfessionals(6)
      if (newPlan === 'enterprise') setMaxProfessionals(999)
    }
  }

  function toggleModule(moduleId: string) {
    setActiveModules((prev) =>
      prev.includes(moduleId) ? prev.filter((id) => id !== moduleId) : [...prev, moduleId]
    )
  }

  function addDaysToTrial(days: number) {
    const cur = new Date(trialEndsDate)
    cur.setDate(cur.getDate() + days)
    setTrialEndsDate(cur.toISOString().slice(0, 10))
    setStatus('trialing')
    showToast({ message: `+${days} dias adicionados ao trial.`, tone: 'success' })
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setIsSubmitting(true)
    try {
      await updateTenantSubscription(tenant.id, {
        plan,
        subscriptionStatus: status,
        trialEndsAt: `${trialEndsDate}T23:59:59Z`,
        activeModules,
        maxProfessionals: Number(maxProfessionals) || 5,
      })
      showToast({ message: 'Assinatura e módulos do tenant atualizados!', tone: 'success' })
      onSaved()
    } catch (err) {
      showToast({ message: (err as Error).message, tone: 'error' })
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Modal
      open
      onClose={onClose}
      size="lg"
      title={`Gestão de Assinatura & Módulos: ${tenant.name}`}
      subtitle="Defina o plano comercial, status de cobrança e quais módulos estarão liberados para este estabelecimento."
    >
      <Form className="mt-4 grid gap-5" onSubmit={handleSubmit}>
        <div className="grid gap-4 sm:grid-cols-3">
          <Field label="Plano SaaS">
            <Select
              value={plan}
              onChange={(e) => handlePlanChange(e.target.value as TenantPlan)}
            >
              <option value="trial">Trial Grátis (14 dias)</option>
              <option value="starter">Plano Starter (Básico)</option>
              <option value="pro">Plano PRO (ERP Completo)</option>
              <option value="enterprise">Plano Enterprise (Ilimitado)</option>
              <option value="custom">Personalizado (Custom)</option>
            </Select>
          </Field>

          <Field label="Status da Assinatura">
            <Select
              value={status}
              onChange={(e) => setStatus(e.target.value as TenantSubscriptionStatus)}
            >
              <option value="trialing">Em Teste (Trial)</option>
              <option value="active">Ativo (Pago)</option>
              <option value="past_due">Em Atraso / Bloqueio</option>
              <option value="canceled">Cancelado</option>
            </Select>
          </Field>

          <Field label="Limite de Profissionais">
            <Input
              type="number"
              min={1}
              value={maxProfessionals}
              onChange={(e) => setMaxProfessionals(Number(e.target.value))}
            />
          </Field>
        </div>

        {/* Período de Teste & Extensões Rápidas */}
        <div className="rounded-xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4">
          <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <span className="text-xs font-bold uppercase tracking-wider text-[var(--color-ink-soft)]">
                Término do Período de Teste / Vencimento
              </span>
              <div className="mt-1 flex items-center gap-2">
                <Calendar className="h-4 w-4 text-[var(--color-primary)]" />
                <Input
                  type="date"
                  className="w-40 text-xs font-semibold"
                  value={trialEndsDate}
                  onChange={(e) => setTrialEndsDate(e.target.value)}
                />
              </div>
            </div>

            <div className="flex flex-wrap items-center gap-1.5 pt-2 sm:pt-0">
              <span className="text-xs text-[var(--color-ink-soft)] mr-1">Extensão rápida:</span>
              <Button type="button" size="sm" variant="secondary" onClick={() => addDaysToTrial(7)}>
                +7 dias
              </Button>
              <Button type="button" size="sm" variant="secondary" onClick={() => addDaysToTrial(14)}>
                +14 dias
              </Button>
              <Button type="button" size="sm" variant="secondary" onClick={() => addDaysToTrial(30)}>
                +30 dias
              </Button>
            </div>
          </div>
        </div>

        {/* Módulos Habilitados */}
        <div>
          <div className="flex items-center justify-between mb-2">
            <div>
              <h4 className="text-sm font-bold text-[var(--color-ink)] flex items-center gap-2">
                <ShieldCheck className="h-4 w-4 text-[var(--color-primary)]" />
                Módulos Habilitados para este Estabelecimento
              </h4>
              <p className="text-xs text-[var(--color-ink-soft)]">
                Módulos desmarcados serão ocultados do menu deste tenant ou exibirão aviso de upgrade.
              </p>
            </div>
            <span className="text-xs font-bold text-[var(--color-primary)]">
              {activeModules.length} de {AVAILABLE_MODULES.length} ativos
            </span>
          </div>

          <div className="grid gap-2.5 sm:grid-cols-2 max-h-72 overflow-y-auto p-1 border rounded-xl border-border">
            {AVAILABLE_MODULES.map((mod) => {
              const isChecked = activeModules.includes(mod.id)
              return (
                <label
                  key={mod.id}
                  className={`flex items-start gap-2.5 p-2.5 rounded-lg border cursor-pointer transition-colors ${
                    isChecked
                      ? 'border-[var(--color-primary)]/40 bg-[var(--color-primary)]/5'
                      : 'border-border opacity-70 hover:opacity-100'
                  }`}
                >
                  <Checkbox
                    checked={isChecked}
                    onChange={() => toggleModule(mod.id)}
                    className="mt-0.5"
                  />
                  <div className="min-w-0">
                    <span className="text-xs font-bold text-[var(--color-ink)] block">
                      {mod.label}
                    </span>
                    <span className="text-[11px] text-[var(--color-ink-soft)] line-clamp-1">
                      {mod.desc}
                    </span>
                  </div>
                </label>
              )
            })}
          </div>
        </div>

        <div className="mt-2 flex items-center justify-end gap-2 border-t pt-4">
          <Button type="button" variant="secondary" onClick={onClose}>
            Cancelar
          </Button>
          <Button type="submit" loading={isSubmitting}>
            Salvar Módulos & Assinatura
          </Button>
        </div>
      </Form>
    </Modal>
  )
}

function TenantFormModal({
  tenant,
  onClose,
  onSaved,
}: {
  tenant: Tenant | null
  onClose: () => void
  onSaved: () => void
}) {
  const { showToast } = useToast()
  const [name, setName] = useState(tenant?.name ?? '')
  const [slug, setSlug] = useState(tenant?.slug ?? '')
  const [active, setActive] = useState(tenant?.active ?? true)
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const trimmedName = name.trim()
    const trimmedSlug = slug.trim().toLowerCase()

    if (!trimmedName) {
      setError('Informe o nome do estabelecimento.')
      return
    }
    if (!trimmedSlug) {
      setError('Informe o identificador (slug).')
      return
    }
    if (!SLUG_PATTERN.test(trimmedSlug)) {
      setError('O identificador deve conter apenas letras minúsculas, números e hífens (ex.: barbearia-norte).')
      return
    }

    setIsSubmitting(true)
    setError('')

    try {
      const input = { name: trimmedName, slug: trimmedSlug, active }
      if (tenant) {
        await updateTenant(tenant.id, input)
        showToast({ message: 'Estabelecimento atualizado.', tone: 'success' })
      } else {
        await createTenant(input)
        showToast({ message: 'Estabelecimento criado.', tone: 'success' })
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
          <Button loading={isSubmitting} type="submit" form="tenant-form">
            Salvar
          </Button>
        </>
      }
      onClose={onClose}
      open
      size="md"
      subtitle={tenant ? 'Atualize os dados cadastrais do estabelecimento.' : 'Cadastre um novo estabelecimento.'}
      title={tenant ? 'Editar estabelecimento' : 'Novo estabelecimento'}
    >
      <Form id="tenant-form" onSubmit={handleSubmit}>
        <Field label="Nome">
          <Input
            autoFocus
            className="mt-2"
            onChange={(event) => setName(event.target.value)}
            placeholder="Ex.: Barbearia do Centro"
            value={name}
          />
        </Field>
        <Field
          hint="Usado no link público de agendamento (ex.: meubarbearia). Apenas minúsculas, números e hífens."
          label="Identificador (slug)"
        >
          <Input
            className="mt-2"
            disabled={Boolean(tenant)}
            onChange={(event) => setSlug(event.target.value.toLowerCase())}
            placeholder="meubarbearia"
            value={slug}
          />
        </Field>
        <Field label=" ">
          <Checkbox checked={active} onChange={(event) => setActive(event.target.checked)} />
          <span className="ml-2 text-sm font-bold text-[var(--color-ink)]">Estabelecimento ativo</span>
        </Field>
        {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
      </Form>
    </Modal>
  )
}
