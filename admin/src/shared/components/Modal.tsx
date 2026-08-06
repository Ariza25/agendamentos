import type { ReactNode } from 'react'
import { X } from 'lucide-react'

type ModalProps = {
  open: boolean
  onClose: () => void
  title?: string
  subtitle?: string
  children: ReactNode
  footer?: ReactNode
  size?: 'sm' | 'md' | 'lg'
}

const sizeClasses = {
  sm: 'max-w-md',
  md: 'max-w-xl',
  lg: 'max-w-3xl',
}

export function Modal({ open, onClose, title, subtitle, children, footer, size = 'md' }: ModalProps) {
  if (!open) {
    return null
  }

  return (
    <div
      aria-modal="true"
      className="app-modal-backdrop fixed inset-0 z-50 flex items-center justify-center overflow-y-auto p-4"
      role="dialog"
    >
      <div
        className="absolute inset-0"
        onClick={onClose}
        onKeyDown={(event) => {
          if (event.key === 'Escape') {
            onClose()
          }
        }}
        role="presentation"
      />
      <div
        className={[
          'app-modal-surface relative my-8 w-full rounded-[1.75rem] p-6 shadow-2xl',
          sizeClasses[size],
        ].join(' ')}
        style={{ animation: 'modal-in 180ms ease-out' }}
      >
        <div className="flex items-start justify-between gap-4">
          <div>
            {title ? <h3 className="font-headline text-2xl font-extrabold text-[var(--color-ink)]">{title}</h3> : null}
            {subtitle ? <p className="mt-1 text-sm text-[var(--color-ink-soft)]">{subtitle}</p> : null}
          </div>
          <button
            aria-label="Fechar"
            className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl text-[var(--color-ink-soft)] transition-colors hover:bg-[var(--color-bg-muted)]"
            onClick={onClose}
            type="button"
          >
            <X className="h-5 w-5" />
          </button>
        </div>
        <div className="mt-5">{children}</div>
        {footer ? <div className="mt-6 flex flex-col-reverse gap-3 sm:flex-row sm:justify-end">{footer}</div> : null}
      </div>
    </div>
  )
}
