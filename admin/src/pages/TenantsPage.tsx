import {
  Activity,
  AlertTriangle,
  Building2,
  Calendar,
  CheckCircle,
  CheckCircle2,
  Clock,
  Database,
  FileText,
  Layers,
  Lock,
  Plus,
  RefreshCw,
  Search,
  Server,
  ShieldCheck,
  SlidersHorizontal,
  Sparkles,
  Unlock,
  XCircle,
} from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import {
  checkSystemReadiness,
  createTenant,
  deleteTenant,
  fetchAuditLogs,
  fetchSystemStatus,
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
import type {
  AuditLog,
  SystemStatus,
  Tenant,
  TenantPlan,
  TenantSubscriptionStatus,
  UpdateTenantSubscriptionInput,
} from '../types'

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
  type MasterTab = 'tenants' | 'audit' | 'system'
  const [activeTab, setActiveTab] = useState<MasterTab>('tenants')

  const [tenants, setTenants] = useState<Tenant[]>([])
  const [loading, setLoading] = useState(true)

  const [isFormOpen, setIsFormOpen] = useState(false)
  const [editingTenant, setEditingTenant] = useState<Tenant | null>(null)

  const [managingSubscriptionTenant, setManagingSubscriptionTenant] = useState<Tenant | null>(null)

  const [tenantToDelete, setTenantToDelete] = useState<Tenant | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  // Audit state
  const [auditLogs, setAuditLogs] = useState<AuditLog[]>([])
  const [auditTotal, setAuditTotal] = useState(0)
  const [auditLoading, setAuditLoading] = useState(false)
  const [auditTenantFilter, setAuditTenantFilter] = useState('')
  const [selectedLogDetails, setSelectedLogDetails] = useState<AuditLog | null>(null)

  // System status state
  const [systemStatus, setSystemStatus] = useState<SystemStatus | null>(null)
  const [systemLoading, setSystemLoading] = useState(false)
  const [readinessLoading, setReadinessLoading] = useState(false)
  const [readinessResult, setReadinessResult] = useState<{ status: string; database?: string; timestamp?: string } | null>(null)

  function loadTenants() {
    setLoading(true)
    fetchTenants()
      .then(setTenants)
      .catch(() => {
        showToast({ message: 'Não foi possível carregar os estabelecimentos.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  function loadAuditLogs() {
    setAuditLoading(true)
    fetchAuditLogs({ tenantId: auditTenantFilter, limit: 50, offset: 0 })
      .then((res) => {
        setAuditLogs(res.logs || [])
        setAuditTotal(res.total || 0)
      })
      .catch(() => {
        showToast({ message: 'Não foi possível carregar a trilha de auditoria.', tone: 'error' })
      })
      .finally(() => setAuditLoading(false))
  }

  function loadSystemStatus() {
    setSystemLoading(true)
    fetchSystemStatus()
      .then(setSystemStatus)
      .catch(() => {
        showToast({ message: 'Não foi possível carregar o status do sistema.', tone: 'error' })
      })
      .finally(() => setSystemLoading(false))
  }

  async function handleCheckReadiness() {
    setReadinessLoading(true)
    try {
      const res = await checkSystemReadiness()
      setReadinessResult(res)
      showToast({ message: `Sistema pronto: Banco ${res.database || 'OK'}`, tone: 'success' })
      loadSystemStatus()
    } catch (err) {
      showToast({ message: `Falha na verificação: ${(err as Error).message}`, tone: 'error' })
    } finally {
      setReadinessLoading(false)
    }
  }

  useEffect(() => {
    loadTenants()
  }, [])

  useEffect(() => {
    if (activeTab === 'audit') {
      loadAuditLogs()
    } else if (activeTab === 'system') {
      loadSystemStatus()
    }
  }, [activeTab, auditTenantFilter])

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

  function getAuditActionBadge(action: string) {
    if (action.startsWith('checkout.')) {
      return <Badge tone="success">Comanda / Venda</Badge>
    }
    if (action.startsWith('cash_register.')) {
      return <Badge tone="accent">Caixa Diário</Badge>
    }
    if (action.startsWith('tenant.')) {
      return <Badge tone="warning">Estabelecimento / Plano</Badge>
    }
    if (action.startsWith('commission.')) {
      return <Badge tone="neutral">Comissões</Badge>
    }
    if (action.startsWith('auth.')) {
      return <Badge tone="neutral">Autenticação</Badge>
    }
    return <Badge tone="neutral">{action}</Badge>
  }

  return (
    <PageContainer>
      {/* Master Hub Navigation Tabs */}
      <div className="flex flex-wrap items-center gap-2 rounded-xl border border-border bg-card p-1.5 shadow-sm">
        <button
          type="button"
          onClick={() => setActiveTab('tenants')}
          className={`flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-semibold transition-all ${
            activeTab === 'tenants'
              ? 'bg-primary text-primary-foreground shadow'
              : 'text-muted-foreground hover:bg-muted hover:text-foreground'
          }`}
        >
          <Building2 className="h-4 w-4" />
          Estabelecimentos & Planos
        </button>
        <button
          type="button"
          onClick={() => setActiveTab('audit')}
          className={`flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-semibold transition-all ${
            activeTab === 'audit'
              ? 'bg-primary text-primary-foreground shadow'
              : 'text-muted-foreground hover:bg-muted hover:text-foreground'
          }`}
        >
          <ShieldCheck className="h-4 w-4" />
          Trilha de Auditoria (LGPD)
        </button>
        <button
          type="button"
          onClick={() => setActiveTab('system')}
          className={`flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-semibold transition-all ${
            activeTab === 'system'
              ? 'bg-primary text-primary-foreground shadow'
              : 'text-muted-foreground hover:bg-muted hover:text-foreground'
          }`}
        >
          <Server className="h-4 w-4" />
          Saúde do Sistema & Filas
        </button>
      </div>

      {activeTab === 'tenants' && (
        <>
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
        </>
      )}

      {activeTab === 'audit' && (
        <Section
          actions={
            <div className="flex items-center gap-2">
              <select
                className="rounded-lg border border-border bg-background px-3 py-1.5 text-xs font-semibold text-foreground"
                value={auditTenantFilter}
                onChange={(e) => setAuditTenantFilter(e.target.value)}
              >
                <option value="">Todos os Estabelecimentos</option>
                {tenants.map((t) => (
                  <option key={t.id} value={t.id}>
                    {t.name} ({t.slug})
                  </option>
                ))}
              </select>
              <Button
                size="sm"
                variant="secondary"
                icon={<RefreshCw className="h-3.5 w-3.5" />}
                loading={auditLoading}
                onClick={loadAuditLogs}
              >
                Atualizar
              </Button>
            </div>
          }
          subtitle="Registro cronológico imutável de operações sensíveis para compliance e auditoria LGPD."
          title="Trilha de Auditoria do Sistema (LGPD)"
        >
          {auditLoading ? (
            <Spinner />
          ) : auditLogs.length === 0 ? (
            <EmptyState
              icon={<ShieldCheck className="h-10 w-10 text-emerald-500" />}
              title="Sem registros de auditoria"
              description="Nenhuma ação auditável registrada com os filtros selecionados."
            />
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full border-collapse text-left text-sm">
                <thead>
                  <tr className="border-b border-border text-xs uppercase text-muted-foreground">
                    <th className="py-3 px-4">Data / Hora</th>
                    <th className="py-3 px-4">Estabelecimento</th>
                    <th className="py-3 px-4">Ação</th>
                    <th className="py-3 px-4">Usuário / Responsável</th>
                    <th className="py-3 px-4">Recurso & ID</th>
                    <th className="py-3 px-4">IP / Origem</th>
                    <th className="py-3 px-4 text-right">Detalhes</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {auditLogs.map((log) => (
                    <tr key={log.id} className="hover:bg-muted/50 transition-colors">
                      <td className="py-3 px-4 text-xs font-medium text-foreground whitespace-nowrap">
                        {formatDateTime(log.createdAt)}
                      </td>
                      <td className="py-3 px-4 text-xs font-semibold text-foreground">
                        {log.tenantId}
                      </td>
                      <td className="py-3 px-4">{getAuditActionBadge(log.action)}</td>
                      <td className="py-3 px-4 text-xs text-foreground">
                        {log.userEmail || log.userId || 'Sistema'}
                      </td>
                      <td className="py-3 px-4 text-xs text-muted-foreground">
                        <span className="font-semibold text-foreground">{log.resource}</span>
                        {log.resourceId && (
                          <span className="ml-1 font-mono text-[11px]">
                            #{log.resourceId.slice(0, 8)}
                          </span>
                        )}
                      </td>
                      <td className="py-3 px-4 text-xs text-muted-foreground font-mono">
                        {log.ipAddress || '-'}
                      </td>
                      <td className="py-3 px-4 text-right">
                        <Button
                          size="sm"
                          variant="secondary"
                          onClick={() => setSelectedLogDetails(log)}
                        >
                          Ver Detalhes
                        </Button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </Section>
      )}

      {activeTab === 'system' && (
        <Section
          actions={
            <div className="flex items-center gap-2">
              <Button
                size="sm"
                variant="secondary"
                icon={<RefreshCw className="h-3.5 w-3.5" />}
                loading={systemLoading}
                onClick={loadSystemStatus}
              >
                Atualizar Status
              </Button>
              <Button
                size="sm"
                icon={<Activity className="h-3.5 w-3.5" />}
                loading={readinessLoading}
                onClick={handleCheckReadiness}
              >
                Testar Prontidão (/ready)
              </Button>
            </div>
          }
          subtitle="Diagnóstico em tempo real da infraestrutura: banco de dados, filas e métricas do Cloud Run."
          title="Saúde do Sistema & Filas"
        >
          {systemLoading ? (
            <Spinner />
          ) : (
            <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
              {/* Card 1: Banco de Dados */}
              <div className="rounded-xl border border-border bg-card p-5 shadow-sm space-y-3">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Database className="h-5 w-5 text-primary" />
                    <h3 className="font-bold text-foreground">Banco de Dados & Pool</h3>
                  </div>
                  <Badge tone={systemStatus?.database.connected ? 'success' : 'danger'}>
                    {systemStatus?.database.connected ? 'Conectado' : 'Desconectado'}
                  </Badge>
                </div>
                <div className="grid grid-cols-2 gap-3 pt-2 text-xs">
                  <div className="rounded-lg bg-muted/50 p-2.5">
                    <span className="text-muted-foreground">Driver Ativo:</span>
                    <p className="mt-1 font-bold text-foreground uppercase">
                      {systemStatus?.database.stats?.driver || 'PostgreSQL'}
                    </p>
                  </div>
                  <div className="rounded-lg bg-muted/50 p-2.5">
                    <span className="text-muted-foreground">Conexões Abertas:</span>
                    <p className="mt-1 font-bold text-foreground">
                      {systemStatus?.database.stats?.openConnections ?? 0}
                    </p>
                  </div>
                  <div className="rounded-lg bg-muted/50 p-2.5">
                    <span className="text-muted-foreground">Conexões em Uso:</span>
                    <p className="mt-1 font-bold text-emerald-600">
                      {systemStatus?.database.stats?.inUse ?? 0}
                    </p>
                  </div>
                  <div className="rounded-lg bg-muted/50 p-2.5">
                    <span className="text-muted-foreground">Conexões Ociosas:</span>
                    <p className="mt-1 font-bold text-foreground">
                      {systemStatus?.database.stats?.idle ?? 0}
                    </p>
                  </div>
                </div>
                {readinessResult && (
                  <div className="mt-2 rounded-lg border border-emerald-500/30 bg-emerald-500/10 p-2.5 text-xs text-emerald-700 dark:text-emerald-300">
                    ✓ Teste de Prontidão: <strong>{readinessResult.status}</strong> (Database:{' '}
                    {readinessResult.database})
                  </div>
                )}
              </div>

              {/* Card 2: Fila Outbox */}
              <div className="rounded-xl border border-border bg-card p-5 shadow-sm space-y-3">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Activity className="h-5 w-5 text-indigo-500" />
                    <h3 className="font-bold text-foreground">Fila Transacional Outbox</h3>
                  </div>
                  <Badge tone="accent">Postgres LISTEN/NOTIFY</Badge>
                </div>
                <div className="grid grid-cols-3 gap-3 pt-2 text-xs">
                  <div className="rounded-lg bg-muted/50 p-2.5">
                    <span className="text-muted-foreground">Pendentes:</span>
                    <p className="mt-1 font-bold text-amber-600">
                      {systemStatus?.outbox.pending ?? 0}
                    </p>
                  </div>
                  <div className="rounded-lg bg-muted/50 p-2.5">
                    <span className="text-muted-foreground">Publicados:</span>
                    <p className="mt-1 font-bold text-emerald-600">
                      {systemStatus?.outbox.published ?? 0}
                    </p>
                  </div>
                  <div className="rounded-lg bg-muted/50 p-2.5">
                    <span className="text-muted-foreground">Falhas:</span>
                    <p className="mt-1 font-bold text-rose-600">
                      {systemStatus?.outbox.failed ?? 0}
                    </p>
                  </div>
                </div>
                <p className="text-[11px] text-muted-foreground">
                  Garante entrega at-least-once de eventos cross-instance e notificações sem custos de filas externas.
                </p>
              </div>

              {/* Card 3: Servidor Go & Recursos */}
              <div className="rounded-xl border border-border bg-card p-5 shadow-sm space-y-3">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Server className="h-5 w-5 text-amber-500" />
                    <h3 className="font-bold text-foreground">Processo Go & Memória</h3>
                  </div>
                  <span className="text-xs text-muted-foreground">
                    UTC: {systemStatus?.serverTime || '-'}
                  </span>
                </div>
                <div className="grid grid-cols-3 gap-3 pt-2 text-xs">
                  <div className="rounded-lg bg-muted/50 p-2.5">
                    <span className="text-muted-foreground">Tempo Ativo:</span>
                    <p className="mt-1 font-bold text-foreground">
                      {Math.floor((systemStatus?.uptimeSeconds || 0) / 60)}m{' '}
                      {(systemStatus?.uptimeSeconds || 0) % 60}s
                    </p>
                  </div>
                  <div className="rounded-lg bg-muted/50 p-2.5">
                    <span className="text-muted-foreground">Goroutines:</span>
                    <p className="mt-1 font-bold text-foreground">
                      {systemStatus?.goroutines ?? 0}
                    </p>
                  </div>
                  <div className="rounded-lg bg-muted/50 p-2.5">
                    <span className="text-muted-foreground">Memória Alocada:</span>
                    <p className="mt-1 font-bold text-foreground">
                      {systemStatus?.memoryAllocMb.toFixed(1) ?? 0} MB
                    </p>
                  </div>
                </div>
              </div>

              {/* Card 4: Camadas de Segurança Ativas */}
              <div className="rounded-xl border border-border bg-card p-5 shadow-sm space-y-3">
                <div className="flex items-center gap-2">
                  <ShieldCheck className="h-5 w-5 text-emerald-500" />
                  <h3 className="font-bold text-foreground">Políticas & Segurança Ativas</h3>
                </div>
                <div className="space-y-1.5 pt-1 text-xs">
                  <div className="flex items-center justify-between rounded-lg bg-muted/50 p-2">
                    <span className="text-foreground">PostgreSQL Row-Level Security (RLS)</span>
                    <Badge tone="success">Ativo</Badge>
                  </div>
                  <div className="flex items-center justify-between rounded-lg bg-muted/50 p-2">
                    <span className="text-foreground">
                      Idempotency Middleware (In-Flight 409 + Replay)
                    </span>
                    <Badge tone="success">Ativo</Badge>
                  </div>
                  <div className="flex items-center justify-between rounded-lg bg-muted/50 p-2">
                    <span className="text-foreground">Rate Limiting Multinível (Memory + DB)</span>
                    <Badge tone="success">Ativo</Badge>
                  </div>
                  <div className="flex items-center justify-between rounded-lg bg-muted/50 p-2">
                    <span className="text-foreground">
                      Cloud Run Auto-Scale to Zero (Custo Mínimo)
                    </span>
                    <Badge tone="accent">0 a 2 instâncias</Badge>
                  </div>
                </div>
              </div>
            </div>
          )}
        </Section>
      )}

      {/* Modal: Detalhes do Log de Auditoria */}
      {selectedLogDetails && (
        <Modal
          title="Detalhes do Registro de Auditoria"
          subtitle={`Ação: ${selectedLogDetails.action} | Recurso: ${selectedLogDetails.resource}`}
          onClose={() => setSelectedLogDetails(null)}
          open
          size="lg"
        >
          <div className="space-y-3 text-xs">
            <div className="grid grid-cols-2 gap-2 rounded-lg bg-muted/50 p-3">
              <div>
                <strong>ID:</strong>{' '}
                <span className="font-mono">{selectedLogDetails.id}</span>
              </div>
              <div>
                <strong>Estabelecimento:</strong> {selectedLogDetails.tenantId}
              </div>
              <div>
                <strong>Usuário / E-mail:</strong>{' '}
                {selectedLogDetails.userEmail || selectedLogDetails.userId || 'Sistema'}
              </div>
              <div>
                <strong>IP de Origem:</strong>{' '}
                <span className="font-mono">{selectedLogDetails.ipAddress || '-'}</span>
              </div>
              <div>
                <strong>Data / Hora:</strong> {formatDateTime(selectedLogDetails.createdAt)}
              </div>
              <div>
                <strong>Recurso ID:</strong>{' '}
                <span className="font-mono">{selectedLogDetails.resourceId || '-'}</span>
              </div>
            </div>
            <div>
              <p className="font-bold text-foreground mb-1">
                Payload / Detalhes do Evento (JSON):
              </p>
              <pre className="max-h-72 overflow-y-auto rounded-lg border border-border bg-muted p-3 text-[11px] font-mono text-foreground whitespace-pre-wrap">
                {JSON.stringify(selectedLogDetails.details || {}, null, 2)}
              </pre>
            </div>
          </div>
        </Modal>
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
