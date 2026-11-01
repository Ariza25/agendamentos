import { CalendarDays, CheckCircle2, Eye, EyeOff, Globe, Sparkles } from 'lucide-react'
import { useState } from 'react'
import type { FormEvent } from 'react'
import { Link } from 'react-router-dom'

import { Button } from '../shared/components/Button'
import { Field, Input } from '../shared/components/FormControls'
import { useToast } from '../shared/components/toast/ToastProvider'
import { signupAdmin, type AdminAccount, type AdminAuthError } from '../lib/adminAuth'

type SignupPageProps = {
  onSignup: (user: AdminAccount) => void
}

function slugify(text: string): string {
  return text
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .slice(0, 50)
}

export function SignupPage({ onSignup }: SignupPageProps) {
  const { showToast } = useToast()

  const [businessName, setBusinessName] = useState('')
  const [slug, setSlug] = useState('')
  const [isSlugManual, setIsSlugManual] = useState(false)
  const [segment, setSegment] = useState<'barbearia' | 'salao' | 'geral'>('barbearia')
  const [ownerName, setOwnerName] = useState('')
  const [email, setEmail] = useState('')
  const [phone, setPhone] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [acceptTerms, setAcceptTerms] = useState(false)

  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({})

  function handleBusinessNameChange(name: string) {
    setBusinessName(name)
    if (!isSlugManual) {
      setSlug(slugify(name))
    }
  }

  function handleSlugChange(value: string) {
    setIsSlugManual(true)
    setSlug(slugify(value))
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError('')
    setFieldErrors({})

    if (!acceptTerms) {
      setError('É obrigatório aceitar os Termos de Uso e a Política de Privacidade.')
      return
    }

    if (!slug) {
      setError('Informe um identificador (link) para o seu estabelecimento.')
      return
    }

    setIsSubmitting(true)

    try {
      const user = await signupAdmin({
        businessName: businessName.trim(),
        slug: slug.trim(),
        segment,
        ownerName: ownerName.trim(),
        email: email.trim(),
        phone: phone.trim(),
        password,
      })
      showToast({ message: 'Estabelecimento cadastrado com sucesso! Bem-vindo.', tone: 'success' })
      onSignup(user)
    } catch (err) {
      const authErr = err as AdminAuthError
      setError(authErr.message || 'Falha ao cadastrar estabelecimento.')
      if (authErr.fieldErrors) {
        setFieldErrors(authErr.fieldErrors)
      }
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <main className="grid min-h-screen bg-[var(--color-panel)] lg:grid-cols-[1.1fr_0.9fr]">
      <section className="flex min-h-screen justify-center px-6 py-10 md:px-12 md:py-16">
        <div className="mx-auto w-full max-w-lg">
          <div className="flex items-center gap-3">
            <span className="flex h-12 w-12 items-center justify-center rounded-2xl bg-[var(--color-primary)] text-white shadow-md">
              <CalendarDays className="h-7 w-7" />
            </span>
            <div>
              <p className="font-headline text-xl font-black tracking-tight text-[var(--color-ink)]">ERP Agendamentos</p>
              <p className="text-xs font-semibold text-[var(--color-ink-soft)]">Cadastro de novo estabelecimento</p>
            </div>
          </div>

          <h1 className="mt-8 font-headline text-3xl font-black tracking-tight text-[var(--color-ink)] md:text-4xl">
            Crie seu estabelecimento
          </h1>
          <p className="mt-2 text-sm leading-6 text-[var(--color-ink-soft)]">
            Configure seu negócio em menos de 1 minuto com serviços e horários prontos para uso.
          </p>

          <form className="mt-8 grid gap-4" onSubmit={handleSubmit}>
            <div className="rounded-2xl border border-stone-200 bg-white p-5 shadow-sm space-y-4">
              <h2 className="text-xs font-black uppercase tracking-wider text-[var(--color-primary)]">
                1. Dados do Estabelecimento
              </h2>

              <Field error={fieldErrors.businessName} label="Nome do Estabelecimento">
                <Input
                  autoFocus
                  className="mt-1 h-11 w-full px-4"
                  onChange={(e) => handleBusinessNameChange(e.target.value)}
                  placeholder="Ex: Barbearia Estilo, Studio Bella"
                  required
                  value={businessName}
                />
              </Field>

              <Field
                error={fieldErrors.slug}
                hint="Seus clientes acessarão sua página através deste link."
                label="Endereço da sua agenda pública"
              >
                <div className="relative mt-1">
                  <Input
                    className="h-11 w-full pl-9 pr-4 text-sm font-semibold"
                    onChange={(e) => handleSlugChange(e.target.value)}
                    placeholder="meu-estabelecimento"
                    required
                    value={slug}
                  />
                  <Globe className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-stone-400" />
                </div>
                {slug ? (
                  <p className="mt-1.5 text-xs font-medium text-emerald-700">
                    Link público: <strong className="font-bold underline">agendei.me/{slug}</strong>
                  </p>
                ) : null}
              </Field>

              <Field label="Segmento de Atuação">
                <select
                  className="mt-1 h-11 w-full rounded-xl border border-stone-200 bg-white px-3 text-sm font-semibold text-[var(--color-ink)] focus:border-[var(--color-primary)] focus:outline-none"
                  onChange={(e) => setSegment(e.target.value as 'barbearia' | 'salao' | 'geral')}
                  value={segment}
                >
                  <option value="barbearia">Barbearia (Corte, Barba, Combo)</option>
                  <option value="salao">Salão de Beleza / Estética (Manicure, Escova, etc.)</option>
                  <option value="geral">Clínica / Consultório / Serviços Gerais</option>
                </select>
                <p className="mt-1 text-xs text-[var(--color-ink-soft)]">
                  Criaremos serviços e horários sugeridos baseados no seu segmento.
                </p>
              </Field>
            </div>

            <div className="rounded-2xl border border-stone-200 bg-white p-5 shadow-sm space-y-4">
              <h2 className="text-xs font-black uppercase tracking-wider text-[var(--color-primary)]">
                2. Conta do Administrador
              </h2>

              <Field error={fieldErrors.ownerName} label="Seu Nome Completo">
                <Input
                  className="mt-1 h-11 w-full px-4"
                  onChange={(e) => setOwnerName(e.target.value)}
                  placeholder="Ex: Carlos Oliveira"
                  required
                  value={ownerName}
                />
              </Field>

              <div className="grid gap-3 sm:grid-cols-2">
                <Field error={fieldErrors.email} label="E-mail de Acesso">
                  <Input
                    autoComplete="email"
                    className="mt-1 h-11 w-full px-4"
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="voce@empresa.com"
                    required
                    type="email"
                    value={email}
                  />
                </Field>

                <Field error={fieldErrors.phone} label="WhatsApp / Celular">
                  <Input
                    autoComplete="tel"
                    className="mt-1 h-11 w-full px-4"
                    onChange={(e) => setPhone(e.target.value)}
                    placeholder="(11) 99999-9999"
                    required
                    type="tel"
                    value={phone}
                  />
                </Field>
              </div>

              <Field error={fieldErrors.password} hint="Mínimo 8 caracteres." label="Senha">
                <div className="relative mt-1">
                  <Input
                    autoComplete="new-password"
                    className="h-11 w-full px-4 pr-12"
                    minLength={8}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="Crie uma senha segura"
                    required
                    type={showPassword ? 'text' : 'password'}
                    value={password}
                  />
                  <button
                    aria-label={showPassword ? 'Ocultar senha' : 'Mostrar senha'}
                    className="absolute right-2 top-1/2 flex h-8 w-8 -translate-y-1/2 items-center justify-center rounded-lg text-stone-500 transition-colors hover:text-[var(--color-primary)]"
                    onClick={() => setShowPassword((current) => !current)}
                    type="button"
                  >
                    {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                  </button>
                </div>
              </Field>
            </div>

            <div className="rounded-xl border border-stone-200 bg-white p-4">
              <label className="flex items-start gap-3 cursor-pointer">
                <input
                  checked={acceptTerms}
                  className="mt-1 h-4 w-4 rounded border-stone-300 text-[var(--color-primary)] focus:ring-[var(--color-primary)]"
                  onChange={(e) => setAcceptTerms(e.target.checked)}
                  required
                  type="checkbox"
                />
                <span className="text-xs text-[var(--color-ink)] leading-relaxed">
                  Declaro que li e concordo com os <strong>Termos de Uso</strong> e com a{' '}
                  <strong>Política de Privacidade (LGPD)</strong> do ERP Agendamentos.
                </span>
              </label>
            </div>

            {error ? (
              <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">
                {error}
              </p>
            ) : null}

            <Button className="h-12 w-full rounded-xl px-5 text-sm font-black" loading={isSubmitting} type="submit">
              Criar estabelecimento e entrar
            </Button>

            <p className="text-center text-sm font-semibold text-[var(--color-ink-soft)]">
              Já possui uma conta?{' '}
              <Link className="font-bold text-[var(--color-primary)] hover:underline" to="/">
                Entrar no painel
              </Link>
            </p>
          </form>
        </div>
      </section>

      <section className="hidden min-h-screen flex-col justify-center bg-[var(--color-primary)] px-12 py-16 text-white lg:flex">
        <div className="max-w-md space-y-6">
          <div className="inline-flex items-center gap-2 rounded-full bg-white/15 px-3 py-1 text-xs font-bold uppercase tracking-wide">
            <Sparkles className="h-3.5 w-3.5" />
            Tudo em uma única plataforma
          </div>

          <h2 className="font-headline text-4xl font-black leading-tight tracking-tight">
            A forma moderna de gerenciar seu negócio e agendamentos.
          </h2>

          <p className="text-base leading-7 text-white/90">
            Mais de centenas de estabelecimentos utilizam nossa solução para fidelizar clientes, reduzir faltas e automatizar recebimentos.
          </p>

          <ul className="space-y-4 pt-4 text-sm font-semibold text-white/95">
            <li className="flex items-center gap-3">
              <CheckCircle2 className="h-5 w-5 shrink-0 text-emerald-300" />
              Página de agendamento online pública exclusiva para o seu negócio
            </li>
            <li className="flex items-center gap-3">
              <CheckCircle2 className="h-5 w-5 shrink-0 text-emerald-300" />
              Pagamentos integrados via Pix (Mercado Pago) com baixa automática
            </li>
            <li className="flex items-center gap-3">
              <CheckCircle2 className="h-5 w-5 shrink-0 text-emerald-300" />
              Lembretes automáticos e gestão de faltas (No-Show)
            </li>
            <li className="flex items-center gap-3">
              <CheckCircle2 className="h-5 w-5 shrink-0 text-emerald-300" />
              Fluxo de caixa, controle de estoque e comissão por profissional
            </li>
            <li className="flex items-center gap-3">
              <CheckCircle2 className="h-5 w-5 shrink-0 text-emerald-300" />
              Em total conformidade com a LGPD (Lei nº 13.709/2018)
            </li>
          </ul>
        </div>
      </section>
    </main>
  )
}

