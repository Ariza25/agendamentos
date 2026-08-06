import type { ReactNode } from 'react'

type PageContainerProps = {
  children: ReactNode
  className?: string
  spacingClassName?: string
}

export function PageContainer({ children, className, spacingClassName = 'space-y-6' }: PageContainerProps) {
  return <div className={[spacingClassName, className].filter(Boolean).join(' ')}>{children}</div>
}
