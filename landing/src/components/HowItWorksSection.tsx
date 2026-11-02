import { CalendarClock, ExternalLink, Globe, LayoutDashboard } from 'lucide-react'

const STEPS = [
  {
    icon: LayoutDashboard,
    step: '1',
    title: 'Configure o painel',
    description: 'Acesse seu painel, adicione serviços, profissionais e regras de agenda.',
  },
  {
    icon: Globe,
    step: '2',
    title: 'Compartilhe o link',
    description: 'Divulgue a página pública de agendamento do seu negócio para os clientes.',
  },
  {
    icon: CalendarClock,
    step: '3',
    title: 'Gerencie em tempo real',
    description: 'Receba agendamentos e acompanhe a agenda do dia com atualização imediata.',
  },
]

export function HowItWorksSection() {
  return (
    <section id="como-funciona" className="py-16 sm:py-20">
      <div className="mx-auto max-w-6xl px-6">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="font-headline text-3xl font-extrabold text-ink sm:text-4xl">
            Comece em três passos
          </h2>
          <p className="mt-4 text-lg text-ink-soft">
            Sem instalação e sem complicação: seu negócio agendando online em poucos minutos.
          </p>
        </div>

        <div className="mt-12 grid gap-5 sm:grid-cols-3">
          {STEPS.map((item) => (
            <div
              key={item.step}
              className="relative rounded-xl border border-line bg-bg-elevated p-6"
            >
              <span className="absolute right-5 top-4 font-headline text-5xl font-extrabold text-primary/10">
                {item.step}
              </span>
              <span className="flex size-10 items-center justify-center rounded-lg bg-primary/10 text-primary">
                <item.icon className="size-5" />
              </span>
              <h3 className="mt-4 font-headline text-lg font-bold text-ink">{item.title}</h3>
              <p className="mt-2 text-sm leading-relaxed text-ink-soft">{item.description}</p>
            </div>
          ))}
        </div>

        <p className="mt-8 flex flex-wrap items-center justify-center gap-2 text-sm text-ink-soft">
          <ExternalLink className="size-4" />
          Cada estabelecimento possui uma página própria de agendamento em{' '}
          <code className="rounded bg-panel-alt px-2 py-0.5 font-mono text-xs">/agendar</code>.
        </p>
      </div>
    </section>
  )
}
