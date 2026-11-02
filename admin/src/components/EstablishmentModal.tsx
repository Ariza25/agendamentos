import { Building2, Check, Copy, ExternalLink, Globe } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import { updateCurrentTenant } from '../lib/adminApi'
import { Button } from '../shared/components/Button'
import { Field, Input } from '../shared/components/FormControls'
import { Modal } from '../shared/components/Modal'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { Tenant } from '../types'

type EstablishmentModalProps = {
  open: boolean
  onClose: () => void
  tenant: Tenant | null
  onTenantUpdated: (updated: Tenant) => void
}

function slugify(text: string): string {
  return text
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .slice(0, 50)
}

export function getPublicBookingUrl(slug: string): string {
  if (typeof window === 'undefined') return `/booking/#/booking/${slug}`
  const origin = window.location.origin
  const isPages = window.location.pathname.includes('/agendamentos')
  const base = isPages ? `${origin}/agendamentos/booking` : `${origin}/booking`
  return `${base}/#/booking/${encodeURIComponent(slug || 'agenda')}`
}

export function EstablishmentModal({ open, onClose, tenant, onTenantUpdated }: EstablishmentModalProps) {
  const { showToast } = useToast()

  const [businessName, setBusinessName] = useState(tenant?.name || '')
  const [slug, setSlug] = useState(tenant?.slug || '')
  const [segment, setSegment] = useState<'barbearia' | 'salao' | 'geral'>(
    (tenant?.segment as 'barbearia' | 'salao' | 'geral') || 'barbearia'
  )
  const [isSaving, setIsSaving] = useState(false)
  const [copied, setCopied] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    if (tenant) {
      setBusinessName(tenant.name || '')
      setSlug(tenant.slug || '')
      setSegment((tenant.segment as 'barbearia' | 'salao' | 'geral') || 'barbearia')
    }
  }, [tenant])

  function handleNameChange(value: string) {
    setBusinessName(value)
    // If slug was empty or matched previous auto-generated name, update slug
    if (!slug || slug === slugify(businessName)) {
      setSlug(slugify(value))
    }
  }

  function handleSlugChange(value: string) {
    setSlug(slugify(value))
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
      onTenantUpdated(updated)
      showToast({ message: 'Dados do estabelecimento atualizados com sucesso!', tone: 'success' })
      onClose()
    } catch (err) {
      setError((err as Error).message || 'Falha ao atualizar dados do estabelecimento.')
    } finally {
      setIsSaving(false)
    }
  }

  const publicUrl = getPublicBookingUrl(slug)

  return (
    <Modal
      open={open}
      onClose={onClose}
      size="md"
      subtitle="Gerencie as informações principais do seu estabelecimento e a sua página de agendamento online."
      title="Dados do Estabelecimento"
    >
      <form className="mt-4 grid gap-5" onSubmit={handleSubmit}>
        <div className="rounded-2xl border border-stone-200 bg-white p-5 shadow-sm space-y-4">
          <div className="flex items-center gap-2">
            <Building2 className="h-5 w-5 text-[var(--color-primary)]" />
            <h2 className="text-xs font-black uppercase tracking-wider text-[var(--color-primary)]">
              1. Dados do Estabelecimento
            </h2>
          </div>

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
                onChange={(e) => handleSlugChange(e.target.value)}
                placeholder="meu-estabelecimento"
                required
                value={slug}
              />
              <Globe className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-stone-400" />
            </div>

            {slug ? (
              <div className="mt-3 rounded-xl border border-emerald-200 bg-emerald-50 p-3">
                <div className="flex items-center justify-between gap-2">
                  <div className="min-w-0">
                    <p className="text-[11px] font-bold uppercase tracking-wider text-emerald-800">
                      Link da sua agenda pública:
                    </p>
                    <p className="truncate text-xs font-semibold text-emerald-900 underline">
                      {publicUrl}
                    </p>
                  </div>
                  <div className="flex items-center gap-1.5 shrink-0">
                    <button
                      type="button"
                      onClick={handleCopyLink}
                      className="flex items-center gap-1 rounded-lg border border-emerald-300 bg-white px-2.5 py-1 text-xs font-bold text-emerald-800 shadow-sm transition hover:bg-emerald-100"
                    >
                      {copied ? <Check className="h-3.5 w-3.5 text-emerald-600" /> : <Copy className="h-3.5 w-3.5" />}
                      {copied ? 'Copiado' : 'Copiar'}
                    </button>
                    <a
                      href={publicUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="flex items-center gap-1 rounded-lg border border-emerald-300 bg-white px-2.5 py-1 text-xs font-bold text-emerald-800 shadow-sm transition hover:bg-emerald-100"
                    >
                      <ExternalLink className="h-3.5 w-3.5" />
                      Abrir
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

        {error ? (
          <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">
            {error}
          </p>
        ) : null}

        <div className="flex items-center justify-end gap-3 pt-2">
          <Button onClick={onClose} type="button" variant="secondary">
            Cancelar
          </Button>
          <Button loading={isSaving} type="submit">
            Salvar Dados do Estabelecimento
          </Button>
        </div>
      </form>
    </Modal>
  )
}
