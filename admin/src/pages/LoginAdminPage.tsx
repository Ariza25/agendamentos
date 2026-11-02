import { CalendarDays, Eye, EyeOff, KeyRound } from 'lucide-react'
import { useState } from 'react'
import type { FormEvent } from 'react'
import { Link } from 'react-router-dom'

import { Button } from '../shared/components/Button'
import { Field, Input } from '../shared/components/FormControls'
import { useToast } from '../shared/components/toast/ToastProvider'
import { loginAdmin, requestAdminPasswordReset, type AdminAuthError, type AdminAccount } from '../lib/adminAuth'

type LoginAdminPageProps = {
  onLogin: (user: AdminAccount) => void
}

export function AdminLoginPage({ onLogin }: LoginAdminPageProps) {
  const { showToast } = useToast()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [mfaCode, setMfaCode] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [requiresMfa, setRequiresMfa] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({})
  const [isRecoverySent, setIsRecoverySent] = useState(false)
  const [isRequestingRecovery, setIsRequestingRecovery] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    setIsSubmitting(true)
    setError('')
    setFieldErrors({})

    try {
      const user = await loginAdmin({ email, password, ...(mfaCode ? { mfaCode } : {}) })
      onLogin(user)
    } catch (authError) {
      const err = authError as AdminAuthError
      if (err.code === 'mfa_required') {
        setRequiresMfa(true)
        setError('Informe o codigo do aplicativo autenticador para concluir o acesso.')
      } else {
        setRequiresMfa(false)
        setFieldErrors(err.fieldErrors ?? {})
        setError(err.message)
      }
    } finally {
      setIsSubmitting(false)
    }
  }

  async function handleRecovery() {
    if (!email.trim()) {
      showToast({ message: 'Informe seu email para recuperar a senha.', tone: 'error' })
      return
    }
    setIsRequestingRecovery(true)
    try {
      await requestAdminPasswordReset(email.trim())
      setIsRecoverySent(true)
    } catch (recoveryError) {
      showToast({ message: (recoveryError as Error).message, tone: 'error' })
    } finally {
      setIsRequestingRecovery(false)
    }
  }

  return (
    <main className="grid min-h-screen bg-[var(--color-panel)] lg:grid-cols-[1.05fr_0.95fr]">
      <section className="flex min-h-screen justify-center px-6 pb-10 pt-12 md:px-12 md:pt-16 lg:pt-40">
        <div className="mx-auto w-full max-w-md">
          <div className="flex items-center gap-3">
            <span className="flex h-12 w-12 items-center justify-center rounded-2xl bg-[var(--color-primary)] text-white">
              <CalendarDays className="h-7 w-7" />
            </span>
            <div>
              <p className="font-headline text-xl font-black tracking-tight text-[var(--color-ink)]">ERP Agendamentos</p>
              <p className="text-xs font-semibold text-[var(--color-ink-soft)]">Painel administrativo</p>
            </div>
          </div>

          <h1 className="mt-8 text-center font-headline text-4xl font-black tracking-tight text-[var(--color-ink)] lg:text-left">
            Acesso administrativo
          </h1>
          <p className="mt-3 text-center text-sm leading-6 text-[var(--color-ink-soft)] lg:text-left">
            Entre com uma conta autorizada para gerenciar agenda, serviços, clientes e configurações.
          </p>

          <form className="mt-8 grid gap-4" onSubmit={handleSubmit}>
            <Field error={fieldErrors.email} label="Email">
              <Input
                autoComplete="email"
                className="mt-2 h-12 w-full px-4"
                onChange={(event) => setEmail(event.target.value)}
                placeholder="voce@empresa.com"
                type="email"
                value={email}
              />
            </Field>

            <Field error={fieldErrors.password} label="Senha">
              <div className="relative mt-2">
                <Input
                  autoComplete="current-password"
                  className="h-12 w-full px-4 pr-12"
                  onChange={(event) => setPassword(event.target.value)}
                  placeholder="Sua senha"
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                />
                <button
                  aria-label={showPassword ? 'Ocultar senha' : 'Mostrar senha'}
                  className="absolute right-2 top-1/2 flex h-9 w-9 -translate-y-1/2 items-center justify-center rounded-lg text-stone-500 transition-colors hover:text-[var(--color-primary)]"
                  onClick={() => setShowPassword((current) => !current)}
                  type="button"
                >
                  {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
            </Field>

            {requiresMfa ? (
              <Field error={fieldErrors.mfaCode} hint="Digite o codigo de 6 digitos gerado no aplicativo autenticador." label="Codigo autenticador">
                <div className="relative mt-2">
                  <Input
                    autoComplete="one-time-code"
                    className="h-12 w-full px-4 pl-11"
                    inputMode="numeric"
                    maxLength={6}
                    onChange={(event) => setMfaCode(event.target.value.replace(/\D/g, ''))}
                    placeholder="000000"
                    value={mfaCode}
                  />
                  <KeyRound className="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-[var(--color-ink-soft)]" />
                </div>
              </Field>
            ) : null}

            {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
            {isRecoverySent ? (
              <p className="rounded-xl border border-emerald-200 bg-emerald-50 p-3 text-sm font-bold text-emerald-700">
                Se o email estiver cadastrado, enviamos um link de recuperação de senha.
              </p>
            ) : null}

            <Button className="h-12 w-full rounded-xl px-5 text-sm font-black" loading={isSubmitting} type="submit">
              Entrar
            </Button>

            <button
              className="text-center text-sm font-bold text-[var(--color-ink-soft)] transition-colors hover:text-[var(--color-primary)]"
              disabled={isRequestingRecovery}
              onClick={() => void handleRecovery()}
              type="button"
            >
              {isRequestingRecovery ? 'Enviando...' : 'Esqueci minha senha'}
            </button>

            <div className="mt-2 border-t border-stone-200 pt-4 text-center text-sm font-semibold text-[var(--color-ink-soft)]">
              Não tem uma conta?{' '}
              <Link to="/signup" className="font-bold text-[var(--color-primary)] hover:underline">
                Criar conta
              </Link>
            </div>
          </form>
        </div>
      </section>

      <section className="hidden min-h-screen flex-col justify-center bg-[var(--color-primary)] px-12 py-16 text-white lg:flex">
        <p className="max-w-md font-headline text-4xl font-black leading-tight tracking-tight">
          Organize sua agenda e deixe seus clientes agendarem sozinhos.
        </p>
        <p className="mt-5 max-w-md text-base leading-7 text-white/85">
          Configure os períodos de atendimento, cadastre serviços e profissionais, e disponibilize uma tela pública
          de agendamento para o cliente final.
        </p>
      </section>
    </main>
  )
}
