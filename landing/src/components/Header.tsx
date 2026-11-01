import { CalendarClock } from 'lucide-react'

import { APP_LINKS, SECTION_LINKS } from '../lib/links'

const NAV_ITEMS = [
  { label: 'Funcionalidades', href: SECTION_LINKS.features },
  { label: 'Como funciona', href: SECTION_LINKS.howItWorks },
  { label: 'Planos & Preços', href: SECTION_LINKS.pricing },
]

export function Header() {
  return (
    <header className="sticky top-0 z-50 border-b border-line bg-bg-base/85 backdrop-blur">
      <div className="mx-auto flex max-w-6xl items-center justify-between gap-4 px-6 py-4">
        <a href="#topo" className="flex items-center gap-2 font-headline text-lg font-extrabold text-ink">
          <span className="flex size-8 items-center justify-center rounded-lg bg-primary text-white">
            <CalendarClock className="size-5" />
          </span>
          ERP Agendamentos
        </a>

        <nav className="hidden items-center gap-6 sm:flex">
          {NAV_ITEMS.map((item) => (
            <a
              key={item.href}
              href={item.href}
              className="text-sm font-medium text-ink-soft transition-colors hover:text-ink"
            >
              {item.label}
            </a>
          ))}
        </nav>

        <div className="flex items-center gap-3">
          <a
            href={APP_LINKS.admin}
            className="rounded-lg px-4 py-2 text-sm font-semibold text-ink transition-colors hover:bg-panel-alt"
          >
            Entrar
          </a>
          <a
            href={APP_LINKS.signup}
            className="rounded-lg bg-primary px-4 py-2 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-primary-strong"
          >
            Cadastre-se grátis
          </a>
        </div>
      </div>
    </header>
  )
}
