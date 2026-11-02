import { Building2, Check, Copy, ExternalLink, Globe, Save } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import { fetchCurrentTenant, updateCurrentTenant } from '../lib/adminApi'
import { Button } from '../shared/components/Button'
import { Field, Form, Input } from '../shared/components/FormControls'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { Tenant } from '../types'
import { getPublicBookingUrl } from '../components/EstablishmentModal'

function slugify(text: string): string {
  return text
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .slice(0, 50)
}

export function EstablishmentPage() {
  const { showToast } = useToast()
  const [tenant, setTenant] = useState<Tenant | null>(null)
  const [loading, setLoading] = useState(true)

  const [businessName, setBusinessName] = useState('')
  const [slug, setSlug] = useState('')
  const [segment, setSegment] = useState<'barbearia' | 'salao' | 'geral'>('barbearia')
  const [isSaving, setIsSaving] = useState(false)
  const [copied, setCopied] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    let active = true
    fetchCurrentTenant()
      .then((data) => {
        if (!active) return
        setTenant(data)
        setBusinessName(data.name || '')
        setSlug(data.slug || '')
        setSegment((data.segment as 'barbearia' | 'salao' | 'geral') || 'barbearia')
      })
      .catch(() => {
        if (active) showToast({ message: 'Não foi possível carregar os dados do estabelecimento.', tone: 'error' })
      })
      .finally(() => {
        if (active) setLoading(false)
      })

    return () => {
      active = false
    }
  }, [])

  function handleNameChange(value: string) {
    setBusinessName(value)
    if (!slug || slug === slugify(businessName)) {
      setSlug(slugify(value))
    }
  }

  async function handleCopyLink() {
    const url = getPublicBookingUrl(slug)
    try {
      await navigator.clipboard.writeText(url)
      setCopied(true)
      showToast({ message: 'Link copiado para a área de transferência!', tone: 'success' })
      setTimeout(() => setCopied(false), 2000)
    } catch {
      showToast({ message: 'Falha ao copiar link.', tone: 'error' })
    }
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError('')

    if (!businessName.trim()) {
      setError('Informe o nome do estabelecimento.')
      return
    }

    if (!slug.trim()) {
      setError('Informe o endereço da sua agenda pública (link/slug).')
      return
    }

    setIsSaving(true)

    try {
      const updated = await updateCurrentTenant({
        name: businessName.trim(),
        slug: slug.trim(),
        segment,
      })
      setTenant(updated)
      showToast({ message: 'Dados do estabelecimento atualizados com sucesso!', tone: 'success' })
    } catch (err) {
      setError((err as Error).message || 'Falha ao atualizar dados do estabelecimento.')
    } finally {
      setIsSaving(false)
    }
  }

  if (loading) {
    return (
      <PageContainer>
        <Section>
          <Spinner />
        </Section>
      </PageContainer>
    )
  }

  const publicUrl = getPublicBookingUrl(slug)

  return (
    <PageContainer>
      <Section
        actions={
          <Button
            form="establishment-form"
            icon={<Save className="h-4 w-4" />}
            loading={isSaving}
            type="submit"
          >
            Salvar Dados do Estabelecimento
          </Button>
        }
        subtitle="Gerencie o nome comercial, segmento e o link de acesso da sua agenda pública para os clientes."
        title="Dados do Estabelecimento"
      >
        <p className="mt-3 text-sm text-[var(--color-ink-soft)]">
          Aqui você define o nome comercial que aparece na página pública, o segmento de atuação da sua equipe e personaliza o link exclusivo pelo qual seus clientes realizam agendamentos.
        </p>
      </Section>

      <Form className="mt-6 grid max-w-3xl gap-6" id="establishment-form" onSubmit={handleSubmit}>
        <Section subtitle="Informações principais de identificação do seu negócio." title="Identificação">
          <div className="space-y-4">
            <Field label="Nome do Estabelecimento">
              <Input
                autoFocus
                className="mt-1 h-11 w-full px-4"
                onChange={(e) => handleNameChange(e.target.value)}
                placeholder="Ex: Barbearia Estilo, Studio Bella"
                required
                value={businessName}
              />
            </Field>

            <Field
              hint="Seus clientes acessarão sua página através deste link."
              label="Endereço da sua agenda pública"
            >
              <div className="relative mt-1">
                <Input
                  className="h-11 w-full pl-9 pr-4 text-sm font-semibold"
                  onChange={(e) => setSlug(slugify(e.target.value))}
                  placeholder="meu-estabelecimento"
                  required
                  value={slug}
                />
                <Globe className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-stone-400" />
              </div>

              {slug ? (
                <div className="mt-3 rounded-2xl border border-emerald-200 bg-emerald-50 p-4">
                  <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
                    <div className="min-w-0">
                      <p className="text-xs font-bold uppercase tracking-wider text-emerald-800">
                        Link público de agendamento:
                      </p>
                      <p className="mt-0.5 truncate text-sm font-semibold text-emerald-950 underline">
                        {publicUrl}
                      </p>
                    </div>
                    <div className="flex items-center gap-2 shrink-0">
                      <button
                        type="button"
                        onClick={handleCopyLink}
                        className="flex items-center gap-1.5 rounded-xl border border-emerald-300 bg-white px-3 py-1.5 text-xs font-bold text-emerald-800 shadow-sm transition hover:bg-emerald-100"
                      >
                        {copied ? <Check className="h-4 w-4 text-emerald-600" /> : <Copy className="h-4 w-4" />}
                        {copied ? 'Copiado!' : 'Copiar link'}
                      </button>
                      <a
                        href={publicUrl}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex items-center gap-1.5 rounded-xl border border-emerald-300 bg-white px-3 py-1.5 text-xs font-bold text-emerald-800 shadow-sm transition hover:bg-emerald-100"
                      >
                        <ExternalLink className="h-4 w-4" />
                        Abrir página
                      </a>
                    </div>
                  </div>
                </div>
              ) : null}
            </Field>

            <Field label="Segmento de Atuação">
              <select
                className="mt-1 h-11 w-full rounded-xl border border-stone-200 bg-white px-3 text-sm font-semibold text-[var(--color-ink)] focus:border-[var(--color-primary)] focus:outline-none"
                onChange={(e) => setSegment(e.target.value as 'barbearia' | 'salao' | 'geral')}
                value={segment}
              >
                <option value="barbearia">Barbearia (Corte, Barba, Combo)</option>
                <option value="salao">Salão de Beleza / Estética (Manicure, Escova, etc.)</option>
                <option value="geral">Clínica / Consultório / Serviços Gerais</option>
              </select>
              <p className="mt-1 text-xs text-[var(--color-ink-soft)]">
                Criaremos serviços e horários sugeridos baseados no seu segmento.
              </p>
            </Field>
          </div>
        </Section>

        {error ? (
          <p className="rounded-xl border border-red-200 bg-red-50 p-4 text-sm font-bold text-red-700">
            {error}
          </p>
        ) : null}
      </Form>
    </PageContainer>
  )
}
