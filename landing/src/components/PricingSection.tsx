import { Check, Sparkles } from 'lucide-react'
import { APP_LINKS } from '../lib/links'

interface Plan {
  id: string
  name: string
  price: string
  period: string
  description: string
  popular?: boolean
  features: string[]
  buttonText: string
}

const PLANS: Plan[] = [
  {
    id: 'starter',
    name: 'Starter',
    price: 'R$ 79',
    period: '/mês',
    description: 'Ideal para profissionais autônomos e pequenos espaços em início.',
    features: [
      'Até 2 profissionais cadastrados',
      'Agenda inteligente e agendamentos ilimitados',
      'Página pública de agendamento online',
      'Cadastro completo de clientes e histórico',
      'Lembretes automáticos via email/notificação',
    ],
    buttonText: 'Começar 14 dias Grátis',
  },
  {
    id: 'pro',
    name: 'Pro',
    price: 'R$ 139',
    period: '/mês',
    description: 'Para barbearias, clínicas e salões que buscam gestão completa.',
    popular: true,
    features: [
      'Até 5 profissionais cadastrados',
      'Tudo do plano Starter',
      'Frente de Caixa (PDV) e Comanda Única',
      'Controle de Abertura / Sangria / Fechamento de Caixa',
      'Fechamento e Repasse de Comissões com Recibo',
      'Controle de Estoque e Consumo de Insumos',
      'Módulo Financeiro com DRE Simplificado',
    ],
    buttonText: 'Testar Pro Grátis por 14 Dias',
  },
  {
    id: 'enterprise',
    name: 'Enterprise',
    price: 'R$ 229',
    period: '/mês',
    description: 'Para redes e negócios de alta demanda com múltiplos caixas e serviços.',
    features: [
      'Profissionais e colaboradores ilimitados',
      'Tudo do plano Pro',
      'Clubes de Assinatura e Pacotes de Sessões',
      'Vouchers e Gift Cards Pré-pagos',
      'Fila de Espera inteligente automatizada',
      'Avaliações (Reviews) e Coleta de Feedback',
      'Integração com Google Calendar e Mercado Pago Pix',
      'Suporte prioritário via WhatsApp',
    ],
    buttonText: 'Começar Enterprise Grátis',
  },
]

export function PricingSection() {
  return (
    <section id="precos" className="border-t border-line bg-panel-alt/50 py-20">
      <div className="mx-auto max-w-6xl px-6">
        <div className="mx-auto max-w-2xl text-center">
          <span className="inline-block rounded-full bg-primary/10 px-3 py-1 text-xs font-bold uppercase tracking-wider text-primary">
            Planos Transparentes
          </span>
          <h2 className="mt-3 font-headline text-3xl font-black text-ink sm:text-4xl">
            Comece com 14 dias grátis. Cancele quando quiser.
          </h2>
          <p className="mt-4 text-base text-ink-soft">
            Sem pegadinhas, sem necessidade de cartão de crédito para testar. Escolha o plano perfeito para o tamanho do seu negócio.
          </p>
        </div>

        <div className="mt-14 grid gap-8 md:grid-cols-3">
          {PLANS.map((plan) => (
            <div
              key={plan.id}
              className={`relative flex flex-col justify-between rounded-2xl border p-8 transition-all ${
                plan.popular
                  ? 'border-primary bg-panel shadow-xl ring-2 ring-primary/20 -translate-y-1'
                  : 'border-line bg-panel hover:border-line-strong hover:shadow-md'
              }`}
            >
              {plan.popular && (
                <div className="absolute -top-3.5 left-1/2 -translate-x-1/2 rounded-full bg-primary px-3 py-1 text-xs font-bold tracking-wide text-white shadow-sm flex items-center gap-1.5">
                  <Sparkles className="size-3.5" />
                  MAIS ESCOLHIDO
                </div>
              )}

              <div>
                <h3 className="font-headline text-xl font-bold text-ink">{plan.name}</h3>
                <p className="mt-2 text-xs text-ink-soft min-h-[36px]">{plan.description}</p>

                <div className="mt-6 flex items-baseline gap-1">
                  <span className="font-headline text-4xl font-black text-ink">{plan.price}</span>
                  <span className="text-sm font-semibold text-ink-soft">{plan.period}</span>
                </div>

                <div className="mt-8 space-y-3">
                  <p className="text-xs font-bold uppercase tracking-wider text-ink-soft">Incluso no plano:</p>
                  <ul className="space-y-2.5 text-sm">
                    {plan.features.map((feature, idx) => (
                      <li key={idx} className="flex items-start gap-2.5">
                        <Check className="size-4 shrink-0 text-emerald-600 mt-0.5" />
                        <span className="text-ink text-xs leading-relaxed">{feature}</span>
                      </li>
                    ))}
                  </ul>
                </div>
              </div>

              <div className="mt-8 pt-4 border-t border-line">
                <a
                  href={APP_LINKS.signup}
                  className={`block w-full rounded-xl py-3 text-center text-sm font-bold transition-all shadow-sm ${
                    plan.popular
                      ? 'bg-primary text-white hover:bg-primary-strong shadow-primary/20'
                      : 'bg-panel-alt text-ink hover:bg-line border border-line'
                  }`}
                >
                  {plan.buttonText}
                </a>
                <p className="mt-2 text-center text-[11px] text-ink-soft">
                  14 dias de teste grátis sem compromisso
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}

