import type { ReactNode } from 'react'

type EmptyStateProps = {
  icon?: ReactNode
  title: string
  description?: string
  actions?: ReactNode
}

export function EmptyState({ icon, title, description, actions }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center rounded-2xl border border-dashed border-[var(--color-line)] bg-[var(--color-bg-muted)] px-6 py-14 text-center">
      {icon ? <span className="mb-4 text-[var(--color-ink-soft)]">{icon}</span> : null}
      <h4 className="font-headline text-lg font-extrabold text-[var(--color-ink)]">{title}</h4>
      {description ? <p className="mt-1 max-w-md text-sm text-[var(--color-ink-soft)]">{description}</p> : null}
      {actions ? <div className="mt-5">{actions}</div> : null}
    </div>
  )
}
