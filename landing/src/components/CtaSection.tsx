import { ArrowRight } from 'lucide-react'

import { APP_LINKS } from '../lib/links'

export function CtaSection() {
  return (
    <section className="py-16 sm:py-20">
      <div className="mx-auto max-w-6xl px-6">
        <div className="relative overflow-hidden rounded-2xl bg-primary px-8 py-12 text-center sm:px-16 sm:py-16">
          <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(255,255,255,0.18),transparent_45%)]" />
          <h2 className="relative font-headline text-3xl font-extrabold text-white sm:text-4xl">
            Pronto para organizar sua agenda?
          </h2>
          <p className="relative mx-auto mt-4 max-w-xl text-lg text-white/85">
            Acesse o painel e comece a receber agendamentos online hoje mesmo.
          </p>
          <div className="relative mt-8 flex flex-col items-center justify-center gap-3 sm:flex-row">
            <a
              href={APP_LINKS.signup}
              className="inline-flex items-center gap-2 rounded-lg bg-white px-6 py-3 text-base font-semibold text-primary transition-colors hover:bg-white/90"
            >
              Criar estabelecimento grátis
              <ArrowRight className="size-4" />
            </a>
            <a
              href={APP_LINKS.admin}
              className="inline-flex items-center gap-2 rounded-lg border border-white/40 px-6 py-3 text-base font-semibold text-white transition-colors hover:bg-white/10"
            >
              Acessar painel existente
            </a>
          </div>
        </div>
      </div>
    </section>
  )
}
