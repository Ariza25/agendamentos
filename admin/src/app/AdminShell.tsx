import {
  AlertTriangle,
  Building2,
  CalendarDays,
  Check,
  ChevronDown,
  Copy,
  CreditCard,
  ExternalLink,
  LogOut,
  Menu,
  QrCode,
  ShieldCheck,
  Sparkles,
  X,
  Zap,
} from 'lucide-react'
import { useEffect, useState } from 'react'
import { NavLink, useLocation } from 'react-router-dom'

import type { AdminAccount } from '../lib/adminAuth'
import { fetchCurrentTenant } from '../lib/adminApi'
import { canAccessAdminArea, isMasterAdmin } from '../lib/adminPermissions'
import { roleLabel } from '../lib/adminAuth'
import { Button } from '../shared/components/Button'
import { Modal } from '../shared/components/Modal'
import { useToast } from '../shared/components/toast/ToastProvider'
import { EstablishmentModal, getPublicBookingUrl } from '../components/EstablishmentModal'
import { AdminRoutes } from './AdminRoutes'
import { navigationGroups, navigationItems, type AdminNavigationItem } from './navigation'
import type { Tenant } from '../types'

type AdminShellProps = {
  user: AdminAccount
  onLogout: () => void
}

export function AdminShell({ user, onLogout }: AdminShellProps) {
  const { showToast } = useToast()
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const [currentTenant, setCurrentTenant] = useState<Tenant | null>(null)
  const [isUpgradeModalOpen, setIsUpgradeModalOpen] = useState(false)
  const [isEstablishmentModalOpen, setIsEstablishmentModalOpen] = useState(false)
  const [copiedSlug, setCopiedSlug] = useState(false)
  const location = useLocation()

  async function handleCopyBookingUrl() {
    if (!currentTenant?.slug) return
    const url = getPublicBookingUrl(currentTenant.slug)
    try {
      await navigator.clipboard.writeText(url)
      setCopiedSlug(true)
      showToast({ message: 'Link da agenda pública copiado com sucesso!', tone: 'success' })
      setTimeout(() => setCopiedSlug(false), 2000)
    } catch {
      showToast({ message: 'Falha ao copiar link.', tone: 'error' })
    }
  }

  useEffect(() => {
    setIsMobileMenuOpen(false)
  }, [location.pathname])

  useEffect(() => {
    fetchCurrentTenant()
      .then(setCurrentTenant)
      .catch(() => {
        // Fallback for non-master without tenant record
      })
  }, [user.tenantId])

  const canAccessItem = (item: AdminNavigationItem) => {
    if (item.roles && !item.roles.includes(user.role ?? '')) {
      return false
    }
    if (item.masterOnly && !isMasterAdmin(user.tenantId)) {
      return false
    }
    // Module permissions enforcement for tenant
    if (
      !isMasterAdmin(user.tenantId) &&
      currentTenant &&
      currentTenant.activeModules &&
      currentTenant.activeModules.length > 0
    ) {
      if (!currentTenant.activeModules.includes(item.area)) {
        return false
      }
    }
    return canAccessAdminArea(item.area, user.role)
  }

  const visibleNavigationGroups = navigationGroups
    .map((group) => ({
      ...group,
      items: group.paths
        .map((path) => navigationItems.find((item) => item.to === path))
        .filter((item): item is AdminNavigationItem => Boolean(item))
        .filter(canAccessItem),
    }))
    .filter((group) => group.items.length > 0)

  const defaultRoute = visibleNavigationGroups.flatMap((group) => group.items)[0]?.to ?? '/agenda'

  // Trial calculation
  let trialDaysRemaining = 0
  if (currentTenant?.trialEndsAt) {
    trialDaysRemaining = Math.ceil(
      (new Date(currentTenant.trialEndsAt).getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24)
    )
  }

  const renderNavigation = (mode: 'desktop' | 'mobile') => (
    <nav
      className={
        mode === 'desktop'
          ? 'hidden flex-1 space-y-3 overflow-y-auto pb-2 xl:block'
          : 'flex-1 space-y-3 overflow-y-auto pb-2'
      }
    >
      {visibleNavigationGroups.map((group) => {
        const isGroupActive = group.items.some((item) => item.to === location.pathname)
        return (
          <details className="group" key={group.label} open={mode === 'mobile' || isGroupActive}>
            <summary className="flex cursor-pointer list-none items-center justify-between gap-3 rounded-2xl px-3 py-3 font-headline text-xs font-black uppercase tracking-[0.16em] text-stone-500 transition hover:bg-stone-200 hover:text-stone-900">
              <span className="truncate">{group.label}</span>
              <ChevronDown className="h-4 w-4 shrink-0 transition group-open:rotate-180" />
            </summary>
            <div className="mt-1 space-y-1">
              {group.items.map(({ to, label, icon: Icon }) => (
                <NavLink
                  className={({ isActive }) =>
                    `flex items-center gap-3 rounded-2xl px-3 py-3 font-headline text-sm font-bold tracking-tight transition ${
                      isActive
                        ? 'bg-[var(--color-primary)] text-white shadow-lg shadow-[var(--color-primary)]/20'
                        : 'text-stone-500 hover:bg-stone-200 hover:text-stone-900'
                    }`
                  }
                  end
                  key={to}
                  title={label}
                  to={to}
                >
                  <Icon className="h-4 w-4 shrink-0" />
                  <span className="min-w-0 truncate">{label}</span>
                </NavLink>
              ))}
            </div>
          </details>
        )
      })}
    </nav>
  )

  const shellSidebar = (
    <>
      <div className="mb-2 flex h-20 items-center gap-3">
        <span className="flex h-11 w-11 items-center justify-center rounded-2xl bg-[var(--color-primary)] text-white">
          <CalendarDays className="h-6 w-6" />
        </span>
        <div className="min-w-0">
          <p className="truncate font-headline text-lg font-black tracking-tight text-stone-900">
            {currentTenant?.name || 'ERP Agendamentos'}
          </p>
          <p className="truncate text-xs text-stone-500">
            {isMasterAdmin(user.tenantId)
              ? 'Painel Master Admin'
              : `Plano ${currentTenant?.plan?.toUpperCase() || 'PRO'}`}
          </p>
        </div>
      </div>

      <div className="mb-4 rounded-2xl border border-stone-200 bg-white px-3 py-3">
        <p className="truncate text-sm font-bold text-stone-900">{user.name}</p>
        <p className="truncate text-xs text-stone-500">{user.email}</p>
        <p className="mt-2 text-[11px] font-semibold uppercase tracking-wider text-[var(--color-primary)]">
          {roleLabel(user.role)}
        </p>
      </div>

      {/* Dados do Estabelecimento */}
      <div className="mb-4 rounded-2xl border border-stone-200 bg-white p-3.5 shadow-sm">
        <div className="flex items-center justify-between gap-2">
          <span className="flex items-center gap-1.5 font-headline text-[11px] font-black uppercase tracking-wider text-[var(--color-primary)]">
            <Building2 className="h-3.5 w-3.5" />
            Dados do Estabelecimento
          </span>
          <button
            type="button"
            onClick={() => setIsEstablishmentModalOpen(true)}
            className="rounded-lg bg-stone-100 px-2.5 py-1 text-[11px] font-bold text-stone-700 hover:bg-[var(--color-primary)] hover:text-white transition"
            title="Editar Dados do Estabelecimento"
          >
            Editar
          </button>
        </div>
        <div className="mt-2 min-w-0">
          <p className="truncate font-headline text-sm font-black text-stone-900">
            {currentTenant?.name || 'Meu Estabelecimento'}
          </p>
          <div className="mt-1 flex items-center justify-between text-xs">
            <span className="truncate text-stone-500 font-medium">
              Link: <strong className="font-semibold text-emerald-700">/{currentTenant?.slug || 'agenda'}</strong>
            </span>
            <button
              type="button"
              onClick={handleCopyBookingUrl}
              className="ml-2 font-bold text-[var(--color-primary)] hover:underline shrink-0 text-[11px]"
              title="Copiar link da agenda pública"
            >
              {copiedSlug ? 'Copiado!' : 'Copiar link'}
            </button>
          </div>
        </div>
      </div>

      {renderNavigation('desktop')}

      <div className="mt-auto border-t border-stone-200 pt-4">
        {!isMasterAdmin(user.tenantId) && (
          <button
            type="button"
            onClick={() => setIsUpgradeModalOpen(true)}
            className="mb-2 flex w-full items-center justify-between gap-2 rounded-xl bg-gradient-to-r from-amber-500 to-amber-600 px-3 py-2 text-xs font-bold text-white shadow-sm hover:from-amber-600 hover:to-amber-700 transition"
          >
            <span className="flex items-center gap-1.5">
              <Zap className="h-3.5 w-3.5 fill-current" /> Planos & Upgrade
            </span>
            <span className="rounded bg-black/20 px-1.5 py-0.5 text-[10px]">Ver</span>
          </button>
        )}

        <Button
          className="flex w-full items-center gap-3 rounded-2xl border border-stone-200 bg-white px-3 py-3 font-headline text-sm font-bold tracking-tight text-stone-700 transition hover:border-stone-300 hover:text-stone-900"
          onClick={onLogout}
          title="Sair"
          type="button"
          variant="unstyled"
        >
          <LogOut className="h-4 w-4 shrink-0 text-stone-400" />
          <span className="truncate">Sair</span>
        </Button>
      </div>
    </>
  )

  return (
    <div className="app-shell min-h-screen text-stone-900">
      <div className="flex min-h-screen">
        {/* Mobile Header */}
        <header className="sticky top-0 z-30 flex h-20 items-center justify-between border-b border-stone-200 bg-stone-50 px-4 xl:hidden">
          <div className="flex items-center gap-3">
            <span className="flex h-10 w-10 items-center justify-center rounded-2xl bg-[var(--color-primary)] text-white">
              <CalendarDays className="h-5 w-5" />
            </span>
            <p className="truncate font-headline text-base font-black tracking-tight">
              {currentTenant?.name || 'ERP Agendamentos'}
            </p>
          </div>
          <Button
            aria-expanded={isMobileMenuOpen}
            aria-label="Abrir menu"
            className="flex h-11 w-11 items-center justify-center rounded-2xl border border-stone-200 bg-white p-0 text-stone-700 shadow-sm"
            onClick={() => setIsMobileMenuOpen(true)}
            type="button"
            variant="unstyled"
          >
            <Menu className="h-5 w-5" />
          </Button>
        </header>

        <aside className="sticky top-0 z-30 hidden h-screen min-w-0 flex-col overflow-hidden border-r border-stone-200 bg-stone-50 px-4 py-6 text-stone-700 xl:flex">
          {shellSidebar}
        </aside>

        {isMobileMenuOpen ? (
          <div className="fixed inset-0 z-50 xl:hidden">
            <Button
              aria-label="Fechar menu"
              className="absolute inset-0 cursor-default bg-black/35"
              onClick={() => setIsMobileMenuOpen(false)}
              type="button"
              variant="unstyled"
            />
            <aside className="relative flex h-full w-[min(21rem,calc(100vw-2rem))] flex-col overflow-hidden bg-stone-50 px-4 py-5 text-stone-700 shadow-2xl">
              <div className="mb-4 flex items-center justify-between gap-3">
                <div className="min-w-0">
                  <p className="truncate text-sm font-bold text-stone-900">{user.name}</p>
                  <p className="truncate text-xs text-stone-500">{user.email}</p>
                </div>
                <Button
                  aria-label="Fechar menu"
                  className="flex h-10 w-10 items-center justify-center rounded-2xl border border-stone-200 bg-white p-0 text-stone-700"
                  onClick={() => setIsMobileMenuOpen(false)}
                  type="button"
                  variant="unstyled"
                >
                  <X className="h-5 w-5" />
                </Button>
              </div>

              {/* Dados do Estabelecimento (Mobile) */}
              <div className="mb-4 rounded-2xl border border-stone-200 bg-white p-3 shadow-sm">
                <div className="flex items-center justify-between gap-2">
                  <span className="flex items-center gap-1.5 font-headline text-[11px] font-black uppercase tracking-wider text-[var(--color-primary)]">
                    <Building2 className="h-3.5 w-3.5" />
                    Dados do Estabelecimento
                  </span>
                  <button
                    type="button"
                    onClick={() => {
                      setIsMobileMenuOpen(false)
                      setIsEstablishmentModalOpen(true)
                    }}
                    className="rounded-lg bg-stone-100 px-2.5 py-1 text-[11px] font-bold text-stone-700 hover:bg-[var(--color-primary)] hover:text-white transition"
                  >
                    Editar
                  </button>
                </div>
                <div className="mt-2 min-w-0">
                  <p className="truncate font-headline text-sm font-black text-stone-900">
                    {currentTenant?.name || 'Meu Estabelecimento'}
                  </p>
                  <p className="truncate text-xs font-semibold text-emerald-700 mt-0.5">
                    /{currentTenant?.slug || 'agenda'}
                  </p>
                </div>
              </div>

              {renderNavigation('mobile')}
              <div className="mt-4 border-t border-stone-200 pt-4">
                <Button
                  className="flex w-full items-center gap-3 rounded-2xl bg-[var(--color-primary)] px-3 py-3 font-headline text-sm font-bold tracking-tight text-white transition hover:bg-[var(--color-primary-strong)]"
                  onClick={onLogout}
                  title="Sair"
                  type="button"
                >
                  <LogOut className="h-4 w-4 shrink-0" />
                  <span className="truncate">Sair</span>
                </Button>
              </div>
            </aside>
          </div>
        ) : null}

        <div className="min-w-0 flex-1 bg-[var(--color-panel-alt)]">
          <main className="min-w-0 p-4 sm:p-5 md:p-6 xl:p-8">
            {/* Trial Active Banner */}
            {!isMasterAdmin(user.tenantId) && currentTenant?.subscriptionStatus === 'trialing' && (
              <div className="mb-6 flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-amber-300 bg-amber-50 px-4 py-3 text-amber-900 shadow-sm">
                <div className="flex items-center gap-2.5">
                  <span className="flex h-7 w-7 items-center justify-center rounded-lg bg-amber-500 text-white">
                    <Sparkles className="h-4 w-4" />
                  </span>
                  <div>
                    <span className="text-xs font-black uppercase tracking-wider block">
                      Período de Testes Grátis de 14 Dias Ativo
                    </span>
                    <span className="text-xs text-amber-800">
                      Restam{' '}
                      <strong>{trialDaysRemaining > 0 ? `${trialDaysRemaining} dias` : 'poucas horas'}</strong>.
                      Aproveite todos os módulos PRO e garanta a continuidade da sua agenda!
                    </span>
                  </div>
                </div>

                <button
                  type="button"
                  onClick={() => setIsUpgradeModalOpen(true)}
                  className="rounded-xl bg-amber-600 px-3.5 py-1.5 text-xs font-bold text-white shadow hover:bg-amber-700 transition"
                >
                  Assinar Plano
                </button>
              </div>
            )}

            {/* Past Due / Expirado Alert Banner */}
            {!isMasterAdmin(user.tenantId) && currentTenant?.subscriptionStatus === 'past_due' && (
              <div className="mb-6 flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-rose-300 bg-rose-50 px-4 py-3 text-rose-900 shadow-sm">
                <div className="flex items-center gap-2.5">
                  <span className="flex h-7 w-7 items-center justify-center rounded-lg bg-rose-600 text-white">
                    <AlertTriangle className="h-4 w-4" />
                  </span>
                  <div>
                    <span className="text-xs font-black uppercase tracking-wider block">
                      Assinatura ou Teste Expirado
                    </span>
                    <span className="text-xs text-rose-800">
                      Regularize seu plano para desbloquear todas as ferramentas administrativas.
                    </span>
                  </div>
                </div>

                <button
                  type="button"
                  onClick={() => setIsUpgradeModalOpen(true)}
                  className="rounded-xl bg-rose-600 px-3.5 py-1.5 text-xs font-bold text-white shadow hover:bg-rose-700 transition"
                >
                  Regularizar Agora
                </button>
              </div>
            )}

            <AdminRoutes
              canAccessRoute={(area) => {
                if (area === 'tenants' && !isMasterAdmin(user.tenantId)) {
                  return false
                }
                if (
                  !isMasterAdmin(user.tenantId) &&
                  currentTenant &&
                  currentTenant.activeModules &&
                  currentTenant.activeModules.length > 0
                ) {
                  if (!currentTenant.activeModules.includes(area)) {
                    return false
                  }
                }
                return canAccessAdminArea(area, user.role)
              }}
              defaultRoute={defaultRoute}
              role={user.role}
            />
          </main>
        </div>
      </div>

      {/* Upgrade / SaaS Plans Modal */}
      {isUpgradeModalOpen && (
        <SaaSPlansUpgradeModal
          currentPlan={currentTenant?.plan || 'trial'}
          onClose={() => setIsUpgradeModalOpen(false)}
        />
      )}

      {/* Establishment Settings Modal */}
      <EstablishmentModal
        open={isEstablishmentModalOpen}
        onClose={() => setIsEstablishmentModalOpen(false)}
        tenant={currentTenant}
        onTenantUpdated={(updated) => setCurrentTenant(updated)}
      />
    </div>
  )
}

function SaaSPlansUpgradeModal({
  currentPlan,
  onClose,
}: {
  currentPlan: string
  onClose: () => void
}) {
  const [selectedPlan, setSelectedPlan] = useState<'starter' | 'pro' | 'enterprise'>('pro')
  const [showPix, setShowPix] = useState(false)

  const plans = [
    {
      id: 'starter',
      name: 'Starter',
      price: 'R$ 79',
      period: '/mês',
      desc: 'Ideal para profissionais autônomos iniciando com agendamento online.',
      limit: 'Até 2 profissionais',
      features: [
        'Agenda online com link público',
        'Cadastro e CRM de clientes',
        'Lembretes por e-mail e SMS',
        'Página com marca própria',
      ],
    },
    {
      id: 'pro',
      name: 'PRO (Recomendado)',
      price: 'R$ 139',
      period: '/mês',
      popular: true,
      desc: 'O ERP completo para salões, clínicas e barbearias em expansão.',
      limit: 'Até 6 profissionais',
      features: [
        'Tudo do Starter incluso',
        'Caixa Diário (Gaveta, Sangria e Suprimento)',
        'Frente de Caixa unificada ("Fechar Comanda")',
        'Controle de Estoque & Baixa de Insumos',
        'Cálculo e Fechamento de Comissões da Equipe',
        'Relatórios e DRE Financeiro',
        'Avaliações de Clientes',
      ],
    },
    {
      id: 'enterprise',
      name: 'Enterprise',
      price: 'R$ 229',
      period: '/mês',
      desc: 'Para grandes estabelecimentos que demandam fidelização máxima.',
      limit: 'Profissionais ilimitados',
      features: [
        'Tudo do PRO incluso',
        'Clube de Assinaturas e Mensalidades',
        'Pacotes de Procedimentos & Gift Cards',
        'Sincronização com Google Calendar',
        'Widget para site próprio',
        'Suporte Prioritário VIP',
      ],
    },
  ]

  return (
    <Modal
      open
      onClose={onClose}
      size="lg"
      title="Planos & Assinatura do ERP Agendamentos"
      subtitle="Escolha o plano ideal para a operação do seu estabelecimento e desbloqueie todos os módulos."
    >
      {!showPix ? (
        <div className="mt-4">
          <div className="grid gap-4 md:grid-cols-3">
            {plans.map((p) => {
              const isSelected = selectedPlan === p.id
              return (
                <div
                  key={p.id}
                  onClick={() => setSelectedPlan(p.id as any)}
                  className={`relative flex flex-col justify-between rounded-2xl border-2 p-5 cursor-pointer transition-all ${
                    isSelected
                      ? 'border-[var(--color-primary)] bg-[var(--color-primary)]/5 shadow-md scale-[1.02]'
                      : 'border-border bg-card hover:border-slate-300'
                  }`}
                >
                  {p.popular && (
                    <span className="absolute -top-3 right-4 rounded-full bg-[var(--color-primary)] px-2.5 py-0.5 text-[10px] font-black uppercase tracking-wider text-white shadow">
                      Mais Escolhido
                    </span>
                  )}

                  <div>
                    <h4 className="text-base font-black text-foreground">{p.name}</h4>
                    <p className="mt-1 text-xs text-muted-foreground line-clamp-2">{p.desc}</p>
                    <div className="mt-3 flex items-baseline gap-1">
                      <span className="text-2xl font-black text-foreground">{p.price}</span>
                      <span className="text-xs text-muted-foreground">{p.period}</span>
                    </div>
                    <span className="mt-1 block text-[11px] font-semibold text-primary">{p.limit}</span>

                    <ul className="mt-4 space-y-2 border-t pt-3 text-xs text-muted-foreground">
                      {p.features.map((f, i) => (
                        <li key={i} className="flex items-start gap-1.5">
                          <Check className="h-3.5 w-3.5 shrink-0 text-emerald-600 mt-0.5" />
                          <span>{f}</span>
                        </li>
                      ))}
                    </ul>
                  </div>

                  <div className="mt-5 pt-3">
                    <button
                      type="button"
                      className={`w-full rounded-xl py-2 text-xs font-bold transition ${
                        isSelected
                          ? 'bg-[var(--color-primary)] text-white shadow-sm'
                          : 'bg-muted text-muted-foreground'
                      }`}
                    >
                      {isSelected ? 'Selecionado' : 'Selecionar'}
                    </button>
                  </div>
                </div>
              )
            })}
          </div>

          <div className="mt-6 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between border-t pt-4">
            <span className="text-xs text-muted-foreground">
              Pagamento seguro processado via <strong>Mercado Pago</strong> (Pix ou Cartão em até 12x).
            </span>

            <div className="flex items-center gap-2">
              <Button variant="secondary" onClick={onClose}>
                Fechar
              </Button>
              <Button
                className="bg-emerald-600 text-white hover:bg-emerald-700"
                icon={<QrCode className="h-4 w-4" />}
                onClick={() => setShowPix(true)}
              >
                Pagar via Pix / Mercado Pago
              </Button>
            </div>
          </div>
        </div>
      ) : (
        /* Pix Payment Simulated Checkout */
        <div className="mt-4 rounded-2xl border border-border bg-card p-6 text-center">
          <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-emerald-100 text-emerald-700">
            <QrCode className="h-8 w-8" />
          </div>
          <h4 className="mt-3 text-lg font-black text-foreground">
            Assinatura do Plano {selectedPlan.toUpperCase()}
          </h4>
          <p className="mt-1 text-xs text-muted-foreground">
            Escaneie o QR Code abaixo ou copie a chave Pix no seu aplicativo de banco.
          </p>

          <div className="mx-auto my-4 h-48 w-48 rounded-xl border bg-slate-50 p-3 flex items-center justify-center">
            <div className="text-center">
              <QrCode className="h-28 w-28 mx-auto text-slate-800" />
              <span className="text-[10px] font-mono text-slate-500 block mt-1">PIX MERCADO PAGO</span>
            </div>
          </div>

          <div className="rounded-xl border bg-muted p-2.5 text-xs font-mono text-muted-foreground truncate">
            00020126580014BR.GOV.BCB.PIX0136erp-agendamentos-subs-mercadopago520400005303986
          </div>

          <div className="mt-5 flex justify-center gap-3">
            <Button variant="secondary" onClick={() => setShowPix(false)}>
              Voltar
            </Button>
            <Button
              className="bg-emerald-600 text-white hover:bg-emerald-700"
              onClick={() => {
                alert('Aguardando confirmação do webhook do Mercado Pago. O status será atualizado automaticamente!')
                onClose()
              }}
            >
              Já realizei o pagamento
            </Button>
          </div>
        </div>
      )}
    </Modal>
  )
}
