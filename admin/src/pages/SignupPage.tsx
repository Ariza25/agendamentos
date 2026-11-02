import { CalendarDays, Eye, EyeOff } from 'lucide-react'
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
    .slice(0, 45)
}

export function SignupPage({ onSignup }: SignupPageProps) {
  const { showToast } = useToast()

  const [ownerName, setOwnerName] = useState('')
  const [email, setEmail] = useState('')
  const [phone, setPhone] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [acceptTerms, setAcceptTerms] = useState(false)

  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({})

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError('')
    setFieldErrors({})

    if (!acceptTerms) {
      setError('É obrigatório aceitar os Termos de Uso e a Política de Privacidade.')
      return
    }

    const trimmedName = ownerName.trim()
    if (!trimmedName) {
      setError('Informe seu nome completo.')
      return
    }

    setIsSubmitting(true)

    // Compute automatic business name and slug so backend receives valid establishment defaults
    const uniqueSuffix = Date.now().toString(36).slice(-4)
    const baseSlug = slugify(trimmedName) || 'estabelecimento'
    const fallbackSlug = `${baseSlug}-${uniqueSuffix}`

    try {
      const user = await signupAdmin({
        businessName: `Estabelecimento de ${trimmedName}`,
        slug: fallbackSlug,
        segment: 'barbearia',
        ownerName: trimmedName,
        email: email.trim(),
        phone: phone.trim(),
        password,
      })
      showToast({ message: 'Conta criada com sucesso! Bem-vindo.', tone: 'success' })
      onSignup(user)
    } catch (err) {
      const authErr = err as AdminAuthError
      setError(authErr.message || 'Falha ao cadastrar conta.')
      if (authErr.fieldErrors) {
        setFieldErrors(authErr.fieldErrors)
      }
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-[var(--color-panel)] px-4 py-12 sm:px-6 sm:py-16">
      <div className="w-full max-w-2xl">
        <div className="flex items-center gap-3.5 mb-2">
          <span className="flex h-14 w-14 items-center justify-center rounded-2xl bg-[var(--color-primary)] text-white shadow-md">
            <CalendarDays className="h-8 w-8" />
          </span>
          <div>
            <p className="font-headline text-2xl font-black tracking-tight text-[var(--color-ink)]">ERP Agendamentos</p>
            <p className="text-xs font-bold uppercase tracking-wider text-[var(--color-primary)]">Cadastro de Administrador</p>
          </div>
        </div>

        <h1 className="mt-6 font-headline text-3xl font-black tracking-tight text-[var(--color-ink)] sm:text-4xl">
          Crie sua conta
        </h1>
        <p className="mt-2 text-sm leading-6 text-[var(--color-ink-soft)] sm:text-base">
          Cadastre seu acesso para gerenciar sua agenda, equipe e clientes.
        </p>

        <form className="mt-8 grid gap-5" onSubmit={handleSubmit}>
          <div className="rounded-3xl border border-stone-200 bg-white p-6 sm:p-8 shadow-sm space-y-5">
            <Field error={fieldErrors.ownerName} label="Seu Nome Completo">
              <Input
                autoFocus
                className="mt-1.5 h-12 w-full px-4 text-sm sm:text-base font-medium"
                onChange={(e) => setOwnerName(e.target.value)}
                placeholder="Ex: Carlos Oliveira"
                required
                value={ownerName}
              />
            </Field>

            <div className="grid gap-5 sm:grid-cols-2">
              <Field error={fieldErrors.email} label="E-mail de Acesso">
                <Input
                  autoComplete="email"
                  className="mt-1.5 h-12 w-full px-4 text-sm sm:text-base font-medium"
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
                  className="mt-1.5 h-12 w-full px-4 text-sm sm:text-base font-medium"
                  onChange={(e) => setPhone(e.target.value)}
                  placeholder="(11) 99999-9999"
                  required
                  type="tel"
                  value={phone}
                />
              </Field>
            </div>

            <Field error={fieldErrors.password} hint="Mínimo 8 caracteres." label="Senha">
              <div className="relative mt-1.5">
                <Input
                  autoComplete="new-password"
                  className="h-12 w-full px-4 pr-12 text-sm sm:text-base font-medium"
                  minLength={8}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="Crie uma senha segura"
                  required
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                />
                <button
                  aria-label={showPassword ? 'Ocultar senha' : 'Mostrar senha'}
                  className="absolute right-2 top-1/2 flex h-9 w-9 -translate-y-1/2 items-center justify-center rounded-xl text-stone-500 transition-colors hover:text-[var(--color-primary)]"
                  onClick={() => setShowPassword((current) => !current)}
                  type="button"
                >
                  {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
            </Field>
          </div>

          <div className="rounded-2xl border border-stone-200 bg-white p-4 sm:p-5 shadow-sm">
            <label className="flex items-start gap-3.5 cursor-pointer">
              <input
                checked={acceptTerms}
                className="mt-1 h-4 w-4 rounded border-stone-300 text-[var(--color-primary)] focus:ring-[var(--color-primary)] shrink-0"
                onChange={(e) => setAcceptTerms(e.target.checked)}
                required
                type="checkbox"
              />
              <span className="text-xs sm:text-sm text-[var(--color-ink)] leading-relaxed">
                Declaro que li e concordo com os <strong>Termos de Uso</strong> e com a{' '}
                <strong>Política de Privacidade (LGPD)</strong> do ERP Agendamentos.
              </span>
            </label>
          </div>

          {error ? (
            <p className="rounded-2xl border border-red-200 bg-red-50 p-4 text-sm font-bold text-red-700">
              {error}
            </p>
          ) : null}

          <Button className="h-13 w-full rounded-2xl px-6 text-base font-black shadow-md" loading={isSubmitting} type="submit">
            Cadastrar e entrar
          </Button>

          <p className="text-center text-sm font-semibold text-[var(--color-ink-soft)] pt-1">
            Já possui uma conta?{' '}
            <Link className="font-bold text-[var(--color-primary)] hover:underline" to="/">
              Entrar no painel
            </Link>
          </p>
        </form>
      </div>
    </main>
  )
}
