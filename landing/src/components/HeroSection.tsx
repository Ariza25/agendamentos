import { CalendarCheck2, Clock, Users } from 'lucide-react'

import { APP_LINKS } from '../lib/links'

export function HeroSection() {
  return (
    <section id="topo" className="mx-auto max-w-6xl px-6 pb-16 pt-16 text-center sm:pt-24">
      <p className="mx-auto mb-4 inline-flex items-center gap-2 rounded-full border border-line bg-bg-elevated px-4 py-1.5 text-sm font-medium text-ink-soft">
        Agendamentos online para o seu negócio
      </p>
      <h1 className="mx-auto max-w-3xl font-headline text-4xl font-extrabold leading-tight text-ink sm:text-5xl lg:text-6xl">
        Sua agenda organizada, seus clientes{' '}
        <span className="text-primary">agendando sozinhos</span>
      </h1>
      <p className="mx-auto mt-6 max-w-2xl text-lg leading-relaxed text-ink-soft">
        Gerencie serviços, profissionais, horários e clientes em um só lugar. Deixe seus
        clientes reservarem online a qualquer hora, sem ligações e sem retrabalho.
      </p>
      <div className="mt-10 flex flex-col items-center justify-center gap-3 sm:flex-row">
        <a
          href={APP_LINKS.booking}
          className="rounded-lg bg-primary px-6 py-3 text-base font-semibold text-white shadow-sm transition-colors hover:bg-primary-strong"
        >
          Experimentar agendamento
        </a>
        <a
          href={APP_LINKS.admin}
          className="rounded-lg border border-line bg-bg-elevated px-6 py-3 text-base font-semibold text-ink transition-colors hover:bg-panel-alt"
        >
          Acessar o painel
        </a>
      </div>

      <div className="mx-auto mt-16 grid max-w-3xl gap-4 sm:grid-cols-3">
        <div className="rounded-xl border border-line bg-bg-elevated p-5 text-left">
          <Clock className="size-6 text-primary" />
          <p className="mt-3 font-headline text-2xl font-bold text-ink">Agenda unificada</p>
          <p className="mt-1 text-sm text-ink-soft">Todos os horários em uma visão clara.</p>
        </div>
        <div className="rounded-xl border border-line bg-bg-elevated p-5 text-left">
          <Users className="size-6 text-primary" />
          <p className="mt-3 font-headline text-2xl font-bold text-ink">Equipe e serviços</p>
          <p className="mt-1 text-sm text-ink-soft">Configure quem atende e o que oferece.</p>
        </div>
        <div className="rounded-xl border border-line bg-bg-elevated p-5 text-left">
          <CalendarCheck2 className="size-6 text-primary" />
          <p className="mt-3 font-headline text-2xl font-bold text-ink">Confirmação automática</p>
          <p className="mt-1 text-sm text-ink-soft">Clientes confirmam sem contato manual.</p>
        </div>
      </div>
    </section>
  )
}
