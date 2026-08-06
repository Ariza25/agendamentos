import { CalendarDays, ListChecks, Settings2, Sparkles, UserRound, UsersRound } from 'lucide-react'

const FEATURES = [
  {
    icon: CalendarDays,
    title: 'Agenda de atendimentos',
    description: 'Visualize, crie, remaneje e acompanhe todos os agendamentos em um só lugar.',
  },
  {
    icon: Sparkles,
    title: 'Agendamento online',
    description: 'Página pública por estabelecimento: o cliente escolhe serviço, profissional e horário.',
  },
  {
    icon: UserRound,
    title: 'Cadastro de clientes',
    description: 'Histórico e informações de contato de cada cliente para um atendimento pessoal.',
  },
  {
    icon: ListChecks,
    title: 'Serviços e duração',
    description: 'Configure o catálogo com preço, duração e disponibilidade de cada serviço.',
  },
  {
    icon: UsersRound,
    title: 'Equipe e profissionais',
    description: 'Ative ou desative profissionais e defina quem atende cada serviço.',
  },
  {
    icon: Settings2,
    title: 'Configuração completa',
    description: 'Horário de funcionamento, antecedência e outras regras da agenda do negócio.',
  },
]

export function FeaturesSection() {
  return (
    <section id="funcionalidades" className="border-t border-line bg-bg-elevated/60 py-16 sm:py-20">
      <div className="mx-auto max-w-6xl px-6">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="font-headline text-3xl font-extrabold text-ink sm:text-4xl">
            Tudo o que sua agenda precisa
          </h2>
          <p className="mt-4 text-lg text-ink-soft">
            Do catálogo de serviços à confirmação do horário, o ERP Agendamentos centraliza a
            operação do seu negócio.
          </p>
        </div>

        <div className="mt-12 grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
          {FEATURES.map((feature) => (
            <div
              key={feature.title}
              className="rounded-xl border border-line bg-bg-elevated p-6 transition-shadow hover:shadow-md"
            >
              <span className="flex size-10 items-center justify-center rounded-lg bg-primary/10 text-primary">
                <feature.icon className="size-5" />
              </span>
              <h3 className="mt-4 font-headline text-lg font-bold text-ink">{feature.title}</h3>
              <p className="mt-2 text-sm leading-relaxed text-ink-soft">{feature.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
