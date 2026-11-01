import { Eye, EyeOff } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { Navigate, useLocation } from 'react-router-dom'

import { AdminShell } from './app/AdminShell'
import { AdminLoginPage } from './pages/LoginAdminPage'
import { SignupPage } from './pages/SignupPage'
import { Button } from './shared/components/Button'
import { Field, Input } from './shared/components/FormControls'
import { useToast } from './shared/components/toast/ToastProvider'
import {
  confirmAdminPasswordReset,
  getAdminSession,
  logoutAdmin,
  refreshAdminAccount,
  subscribeToAdminAuthChanges,
  type AdminAuthError,
} from './lib/adminAuth'

function ResetPasswordPage() {
  const { showToast } = useToast()
  const { search } = useLocation()

  const token = new URLSearchParams(search).get('token') ?? ''

  const [newPassword, setNewPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({})
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isDone, setIsDone] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    setIsSubmitting(true)
    setFieldErrors({})

    try {
      await confirmAdminPasswordReset({ token, newPassword })
      logoutAdmin()
      setIsDone(true)
      showToast({ message: 'Senha definida. Entre com seu email e a nova senha.', tone: 'success' })
    } catch (error) {
      const authError = error as AdminAuthError
      setFieldErrors(authError.fieldErrors ?? {})
      showToast({ message: authError.message, tone: 'error' })
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-[var(--color-panel)] px-6 py-10">
      <form
        className="w-full max-w-md rounded-2xl border border-stone-200 bg-white p-6 shadow-sm"
        onSubmit={handleSubmit}
      >
        <h1 className="font-headline text-3xl font-black text-[var(--color-ink)]">Definir senha</h1>

        <Field error={fieldErrors.newPassword || fieldErrors.token} label="Nova senha">
          <div className="relative mt-2">
            <Input
              className="h-12 w-full px-4 pr-12"
              disabled={isDone || !token}
              onChange={(event) => setNewPassword(event.target.value)}
              type={showPassword ? 'text' : 'password'}
              value={newPassword}
            />
            <button
              aria-label={showPassword ? 'Ocultar senha' : 'Mostrar senha'}
              className="absolute right-2 top-1/2 flex h-9 w-9 -translate-y-1/2 items-center justify-center rounded-lg text-stone-500 transition-colors hover:text-[var(--color-primary)] disabled:cursor-not-allowed disabled:opacity-50"
              disabled={isDone || !token}
              onClick={() => setShowPassword((current) => !current)}
              type="button"
            >
              {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
            </button>
          </div>
        </Field>

        <Button className="mt-5 h-12 w-full rounded-xl px-5 text-sm font-black" disabled={isSubmitting || isDone || !token} type="submit">
          {isDone ? 'Senha definida' : isSubmitting ? 'Salvando...' : 'Definir senha'}
        </Button>

        {isDone ? <Navigate replace to="/" /> : null}
      </form>
    </main>
  )
}

export default function App() {
  const { pathname } = useLocation()
  const isResetPasswordPage = pathname === '/reset-password'
  const isSignupPage = pathname === '/signup'

  const [session, setSession] = useState(getAdminSession())

  useEffect(() => {
    if (isResetPasswordPage || isSignupPage) {
      return
    }

    return subscribeToAdminAuthChanges(() => {
      setSession(getAdminSession())
    })
  }, [isResetPasswordPage, isSignupPage])

  useEffect(() => {
    if (isResetPasswordPage || isSignupPage) {
      return
    }
    if (!session) {
      return
    }

    void refreshAdminAccount().then(() => {
      setSession(getAdminSession())
    })
  }, [isResetPasswordPage, isSignupPage])

  if (isResetPasswordPage) {
    return <ResetPasswordPage />
  }

  if (isSignupPage) {
    if (session) {
      return <Navigate replace to="/" />
    }
    return <SignupPage onSignup={() => setSession(getAdminSession())} />
  }

  if (!session) {
    return <AdminLoginPage onLogin={() => setSession(getAdminSession())} />
  }

  function handleLogout() {
    logoutAdmin()
  }

  return <AdminShell onLogout={handleLogout} user={session.user} />
}
