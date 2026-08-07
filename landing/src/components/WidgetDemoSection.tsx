import { useEffect, useRef } from 'react'

import { WIDGET_SCRIPT_URL } from '../lib/links'

export function WidgetDemoSection() {
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const container = containerRef.current
    if (!container) {
      return
    }
    container.innerHTML = ''
    const script = document.createElement('script')
    script.src = WIDGET_SCRIPT_URL
    script.setAttribute('data-slug', 'default')
    script.setAttribute('data-mode', 'inline')
    script.setAttribute('data-label', 'Agendar horário')
    script.setAttribute('data-primary', '#4f46e5')
    container.appendChild(script)
    return () => {
      container.innerHTML = ''
    }
  }, [])

  return (
    <section className="border-t border-line bg-bg-elevated/60 py-16 sm:py-20">
      <div className="mx-auto max-w-6xl px-6">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="font-headline text-3xl font-extrabold text-ink sm:text-4xl">
            Funciona em qualquer site
          </h2>
          <p className="mt-4 text-lg text-ink-soft">
            Cole uma única linha de código e seus clientes agendam sem sair da sua página.
          </p>
        </div>
        <div className="mx-auto mt-10 max-w-xl rounded-2xl border border-line bg-bg-elevated p-4 shadow-sm">
          <div ref={containerRef} />
        </div>
      </div>
    </section>
  )
}
