import { CalendarDays, ChevronDown, LogOut, Menu, X } from 'lucide-react'
import { useEffect, useState } from 'react'
import { NavLink, useLocation } from 'react-router-dom'

import type { AdminAccount } from '../lib/adminAuth'
import { canAccessAdminArea, isMasterAdmin } from '../lib/adminPermissions'
import { roleLabel } from '../lib/adminAuth'
import { Button } from '../shared/components/Button'
import { AdminRoutes } from './AdminRoutes'
import { navigationGroups, navigationItems, type AdminNavigationItem } from './navigation'

type AdminShellProps = {
  user: AdminAccount
  onLogout: () => void
}

export function AdminShell({ user, onLogout }: AdminShellProps) {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const location = useLocation()

  useEffect(() => {
    setIsMobileMenuOpen(false)
  }, [location.pathname])

  const canAccessItem = (item: AdminNavigationItem) => {
    if (item.roles && !item.roles.includes(user.role ?? '')) {
      return false
    }
    if (item.masterOnly && !isMasterAdmin(user.tenantId)) {
      return false
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
          <p className="truncate font-headline text-lg font-black tracking-tight text-stone-900">ERP Agendamentos</p>
          <p className="truncate text-xs text-stone-500">Painel administrativo</p>
        </div>
      </div>
      <div className="mb-6 rounded-2xl border border-stone-200 bg-white px-3 py-3">
        <p className="truncate text-sm font-bold text-stone-900">{user.name}</p>
        <p className="truncate text-xs text-stone-500">
          {roleLabel(user.role)} · {user.email}
        </p>
      </div>

      {renderNavigation('desktop')}

      <div className="mt-6 border-t border-stone-200 pt-6">
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
    </>
  )

  return (
    <div className="min-h-screen bg-[var(--color-panel)]">
      <div className="mx-auto flex min-h-screen w-full max-w-[1920px] flex-col bg-[var(--color-panel)] xl:grid xl:grid-cols-[20rem_minmax(0,1fr)]">
        <header className="sticky top-0 z-40 flex h-16 items-center justify-between border-b border-stone-200 bg-stone-50 px-4 text-stone-700 shadow-sm xl:hidden">
          <div className="flex min-w-0 items-center gap-3">
            <span className="flex h-9 w-9 items-center justify-center rounded-xl bg-[var(--color-primary)] text-white">
              <CalendarDays className="h-5 w-5" />
            </span>
            <p className="truncate font-headline text-base font-black tracking-tight">ERP Agendamentos</p>
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
            <AdminRoutes
              canAccessRoute={(area) => {
                if (area === 'tenants' && !isMasterAdmin(user.tenantId)) {
                  return false
                }
                return canAccessAdminArea(area, user.role)
              }}
              defaultRoute={defaultRoute}
              role={user.role}
            />
          </main>
        </div>
      </div>
    </div>
  )
}
