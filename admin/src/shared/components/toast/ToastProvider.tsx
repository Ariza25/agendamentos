import { createContext, useCallback, useContext, useState } from 'react'
import type { ReactNode } from 'react'
import { CheckCircle2, Info, XCircle } from 'lucide-react'

type ToastTone = 'success' | 'error' | 'info'

type ToastItem = {
  id: number
  message: string
  tone: ToastTone
}

type ToastContextValue = {
  showToast: (input: { message: string; tone?: ToastTone }) => void
}

const ToastContext = createContext<ToastContextValue | null>(null)

let toastSequence = 0

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastItem[]>([])

  const showToast = useCallback((input: { message: string; tone?: ToastTone }) => {
    const id = ++toastSequence
    setToasts((current) => [...current, { id, message: input.message, tone: input.tone ?? 'info' }])
    window.setTimeout(() => {
      setToasts((current) => current.filter((toast) => toast.id !== id))
    }, 4000)
  }, [])

  return (
    <ToastContext.Provider value={{ showToast }}>
      {children}
      <div className="pointer-events-none fixed bottom-4 right-4 z-[60] flex w-full max-w-sm flex-col gap-2">
        {toasts.map((toast) => (
          <ToastCard key={toast.id} toast={toast} />
        ))}
      </div>
    </ToastContext.Provider>
  )
}

function ToastCard({ toast }: { toast: ToastItem }) {
  const icons = {
    success: <CheckCircle2 className="h-5 w-5 text-emerald-600" />,
    error: <XCircle className="h-5 w-5 text-red-600" />,
    info: <Info className="h-5 w-5 text-[var(--color-primary)]" />,
  }

  return (
    <div
      className="app-modal-surface pointer-events-auto flex items-start gap-3 rounded-2xl px-4 py-3 shadow-xl"
      role="status"
      style={{ animation: 'modal-in 180ms ease-out' }}
    >
      <span className="mt-0.5 shrink-0">{icons[toast.tone]}</span>
      <p className="text-sm font-semibold text-[var(--color-ink)]">{toast.message}</p>
    </div>
  )
}

export function useToast() {
  const context = useContext(ToastContext)
  if (!context) {
    throw new Error('useToast deve ser usado dentro de ToastProvider')
  }
  return context
}
