import { forwardRef, useState } from 'react'
import type { ButtonHTMLAttributes, MouseEvent, ReactNode } from 'react'

type ButtonVariant = 'primary' | 'secondary' | 'ghost' | 'unstyled'
type ButtonSize = 'sm' | 'md' | 'lg' | 'icon'

type ButtonProps = Omit<ButtonHTMLAttributes<HTMLButtonElement>, 'onClick'> & {
  variant?: ButtonVariant
  size?: ButtonSize
  loading?: boolean
  icon?: ReactNode
  onClick?: (event: MouseEvent<HTMLButtonElement>) => void | Promise<void>
}

const variantClasses: Record<ButtonVariant, string> = {
  primary: 'bg-[var(--color-primary)] text-white hover:bg-[var(--color-primary-strong)]',
  secondary: 'border border-[var(--color-border)] bg-white text-[var(--color-text)] hover:bg-[var(--color-bg-muted)]',
  ghost: 'bg-transparent text-[var(--color-text)] hover:bg-[var(--color-bg-muted)]',
  unstyled: '',
}

const sizeClasses: Record<ButtonSize, string> = {
  sm: 'min-h-9 px-3 py-2 text-sm',
  md: 'min-h-11 px-4 py-2.5 text-sm',
  lg: 'min-h-12 px-5 py-3 text-base',
  icon: 'h-10 w-10 shrink-0 p-0',
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(function Button(
  {
    variant = 'primary',
    size = 'md',
    loading = false,
    disabled,
    onClick,
    className = '',
    children,
    icon,
    ...props
  },
  ref,
) {
  const [internalLoading, setInternalLoading] = useState(false)
  const isLoading = loading || internalLoading
  const isDisabled = disabled || isLoading
  const baseClassName =
    'inline-flex min-w-0 items-center justify-center gap-2 rounded-xl text-center font-semibold leading-tight transition-colors disabled:cursor-not-allowed disabled:opacity-60'
  const finalClassName = [
    variant === 'unstyled' ? '' : baseClassName,
    variantClasses[variant],
    variant === 'unstyled' ? '' : sizeClasses[size],
    className,
  ]
    .filter(Boolean)
    .join(' ')

  async function handleClick(event: MouseEvent<HTMLButtonElement>) {
    if (!onClick || isDisabled) {
      return
    }

    try {
      const result = onClick(event)
      if (result && typeof (result as Promise<unknown>).then === 'function') {
        setInternalLoading(true)
        await result
      }
    } catch (error) {
      console.error(error)
    } finally {
      setInternalLoading(false)
    }
  }

  return (
    <button
      {...props}
      aria-busy={isLoading ? 'true' : undefined}
      className={finalClassName}
      disabled={isDisabled}
      onClick={(event) => {
        void handleClick(event)
      }}
      ref={ref}
    >
      {isLoading ? (
        <span className="h-4 w-4 animate-spin rounded-full border-2 border-current border-r-transparent" />
      ) : null}
      {!isLoading && icon ? icon : null}
      {children}
    </button>
  )
})
