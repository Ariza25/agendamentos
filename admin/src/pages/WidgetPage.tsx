import { Copy } from 'lucide-react'
import { useEffect, useMemo, useRef, useState } from 'react'

import { fetchCurrentTenant } from '../lib/adminApi'
import { backendOrigin } from '../lib/apiBase'
import { Button } from '../shared/components/Button'
import { Input, Select } from '../shared/components/FormControls'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { Tenant } from '../types'

const widgetScriptUrl = `${backendOrigin}/widget.js`

export function WidgetPage() {
  const { showToast } = useToast()
  const [tenant, setTenant] = useState<Tenant | null>(null)
  const [loading, setLoading] = useState(true)
  const [mode, setMode] = useState<'button' | 'inline'>('button')
  const [label, setLabel] = useState('Agendar horário')
  const [primary, setPrimary] = useState('#4f46e5')

  useEffect(() => {
    fetchCurrentTenant()
      .then(setTenant)
      .catch(() => {
        showToast({ message: 'Não foi possível carregar o estabelecimento.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }, [])

  const slug = tenant?.slug ?? ''

  const snippet = useMemo(() => {
    const attrs = [
      `src="${widgetScriptUrl}"`,
      `data-slug="${slug}"`,
      `data-mode="${mode}"`,
      `data-label="${label}"`,
      `data-primary="${primary}"`,
    ]
    return `<script ${attrs.join(' ')}></script>`
  }, [slug, mode, label, primary])

  async function handleCopy() {
    try {
      await navigator.clipboard.writeText(snippet)
      showToast({ message: 'Código copiado.', tone: 'success' })
    } catch {
      showToast({ message: 'Não foi possível copiar o código.', tone: 'error' })
    }
  }

  return (
    <PageContainer>
      <Section
        subtitle="Deixe seus clientes agendarem direto no seu site, sem sair da página."
        title="Widget de agendamento"
      >
        <div className="mt-2" />
      </Section>

      {loading ? (
        <Section>
          <Spinner />
        </Section>
      ) : (
        <>
          <Section subtitle="Personalize como o widget aparece e copie o código para o seu site." title="Configuração">
            <div className="mt-4 grid gap-4 sm:grid-cols-3">
              <label className="block min-w-0">
                <span className="text-sm font-bold text-[var(--color-ink)]">Modo de exibição</span>
                <Select
                  className="mt-2"
                  onChange={(event) => setMode(event.target.value as 'button' | 'inline')}
                  value={mode}
                >
                  <option value="button">Botão flutuante</option>
                  <option value="inline">Embutido na página</option>
                </Select>
              </label>
              <label className="block min-w-0">
                <span className="text-sm font-bold text-[var(--color-ink)]">Texto do botão</span>
                <Input className="mt-2" onChange={(event) => setLabel(event.target.value)} value={label} />
              </label>
              <label className="block min-w-0">
                <span className="text-sm font-bold text-[var(--color-ink)]">Cor do botão</span>
                <div className="mt-2 flex items-center gap-2">
                  <input
                    className="h-11 w-14 shrink-0 rounded-xl border border-[var(--color-line)] bg-white p-1"
                    onChange={(event) => setPrimary(event.target.value)}
                    type="color"
                    value={primary}
                  />
                  <Input className="font-mono" onChange={(event) => setPrimary(event.target.value)} value={primary} />
                </div>
              </label>
            </div>
          </Section>

          <Section
            actions={
              <Button disabled={!slug} icon={<Copy className="h-4 w-4" />} onClick={() => void handleCopy()}>
                Copiar código
              </Button>
            }
            subtitle={`Identificador do estabelecimento: ${slug || '—'}`}
            title="Código para incorporar"
          >
            <pre className="mt-4 overflow-x-auto rounded-2xl bg-stone-900 p-4 text-sm leading-relaxed text-stone-100">
              <code>{snippet}</code>
            </pre>
          </Section>

          <Section
            subtitle="Veja como o widget fica embutido. O código acima alterna entre botão flutuante e embed inline."
            title="Pré-visualização"
          >
            <WidgetPreview label={label} primary={primary} slug={slug} />
          </Section>
        </>
      )}
    </PageContainer>
  )
}

function WidgetPreview({ label, primary, slug }: { label: string; primary: string; slug: string }) {
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const container = containerRef.current
    if (!container || !slug) {
      return
    }
    container.innerHTML = ''
    const script = document.createElement('script')
    script.src = widgetScriptUrl
    script.setAttribute('data-slug', slug)
    script.setAttribute('data-mode', 'inline')
    script.setAttribute('data-label', label)
    script.setAttribute('data-primary', primary)
    container.appendChild(script)
    return () => {
      container.innerHTML = ''
    }
  }, [label, primary, slug])

  return (
    <div
      className="mt-4 min-h-80 rounded-2xl border border-dashed border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4"
      ref={containerRef}
    />
  )
}
