import { CalendarClock } from 'lucide-react'

import { APP_LINKS, SECTION_LINKS } from '../lib/links'

export function Footer() {
  return (
    <footer className="border-t border-line bg-bg-elevated">
      <div className="mx-auto flex max-w-6xl flex-col items-center gap-6 px-6 py-10 sm:flex-row sm:justify-between">
        <a href="#topo" className="flex items-center gap-2 font-headline text-base font-extrabold text-ink">
          <span className="flex size-7 items-center justify-center rounded-md bg-primary text-white">
            <CalendarClock className="size-4" />
          </span>
          ERP Agendamentos
        </a>

        <nav className="flex flex-wrap items-center justify-center gap-5 text-sm text-ink-soft">
          <a href={SECTION_LINKS.features} className="transition-colors hover:text-ink">
            Funcionalidades
          </a>
          <a href={SECTION_LINKS.howItWorks} className="transition-colors hover:text-ink">
            Como funciona
          </a>
          <a href={APP_LINKS.admin} className="transition-colors hover:text-ink">
            Acessar painel
          </a>
          <a href={APP_LINKS.booking} className="transition-colors hover:text-ink">
            Agendar online
          </a>
        </nav>

        <p className="text-sm text-ink-soft">
          &copy; {new Date().getFullYear()} ERP Agendamentos. Todos os direitos reservados.
        </p>
      </div>
    </footer>
  )
}
