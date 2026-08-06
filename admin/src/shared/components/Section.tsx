import type { ReactNode } from 'react'

type SectionProps = {
  title?: string
  subtitle?: string
  actions?: ReactNode
  children: ReactNode
  className?: string
  titleClassName?: string
}

export function Section({ title, subtitle, actions, children, className, titleClassName }: SectionProps) {
  const sectionClassName = [
    'admin-card-shadow rounded-[1.75rem] border border-[var(--color-line)] bg-white p-6',
    className,
  ]
    .filter(Boolean)
    .join(' ')

  return (
    <section className={sectionClassName}>
      {title || actions ? (
        <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          {title ? (
            <div>
              <h3 className={['font-headline text-2xl font-extrabold text-[var(--color-ink)]', titleClassName].filter(Boolean).join(' ')}>
                {title}
              </h3>
              {subtitle ? <p className="text-sm text-[var(--color-ink-soft)]">{subtitle}</p> : null}
            </div>
          ) : null}
          {actions}
        </div>
      ) : null}
      {children}
    </section>
  )
}
