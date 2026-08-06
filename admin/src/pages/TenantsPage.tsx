import { Building2, Plus } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import { createTenant, deleteTenant, fetchTenants, setTenantActive, updateTenant } from '../lib/adminApi'
import { formatDateTime } from '../lib/format'
import { Badge } from '../shared/components/Badge'
import { Button } from '../shared/components/Button'
import { ConfirmModal } from '../shared/components/ConfirmModal'
import { EmptyState } from '../shared/components/EmptyState'
import { Checkbox, Field, Form, Input } from '../shared/components/FormControls'
import { Modal } from '../shared/components/Modal'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { Tenant } from '../types'

const SLUG_PATTERN = /^[a-z0-9](?:[a-z0-9-]{0,60}[a-z0-9])?$/

export function TenantsPage() {
  const { showToast } = useToast()
  const [tenants, setTenants] = useState<Tenant[]>([])
  const [loading, setLoading] = useState(true)
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [editingTenant, setEditingTenant] = useState<Tenant | null>(null)
  const [tenantToDelete, setTenantToDelete] = useState<Tenant | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  function loadTenants() {
    setLoading(true)
    fetchTenants()
      .then(setTenants)
      .catch(() => {
        showToast({ message: 'Não foi possível carregar os estabelecimentos.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  useEffect(loadTenants, [])

  async function handleToggleActive(tenant: Tenant) {
    try {
      await setTenantActive(tenant.id, !tenant.active)
      showToast({ message: tenant.active ? 'Estabelecimento desativado.' : 'Estabelecimento ativado.', tone: 'success' })
      loadTenants()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    }
  }

  async function handleDelete() {
    if (!tenantToDelete) {
      return
    }
    setIsDeleting(true)
    try {
      await deleteTenant(tenantToDelete.id)
      showToast({ message: 'Estabelecimento removido.', tone: 'success' })
      setTenantToDelete(null)
      loadTenants()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    } finally {
      setIsDeleting(false)
    }
  }

  return (
    <PageContainer>
      <Section
        actions={
          <Button
            icon={<Plus className="h-4 w-4" />}
            onClick={() => {
              setEditingTenant(null)
              setIsFormOpen(true)
            }}
          >
            Novo estabelecimento
          </Button>
        }
        subtitle="Gerencie os estabelecimentos do ERP Agendamentos."
        title="Estabelecimentos"
      >
        <div className="mt-2" />
      </Section>

      {loading ? (
        <Section>
          <Spinner />
        </Section>
      ) : tenants.length === 0 ? (
        <Section>
          <EmptyState
            icon={<Building2 className="h-10 w-10" />}
            description="Nenhum estabelecimento cadastrado. Crie um estabelecimento para começar a operar."
            title="Sem estabelecimentos"
          />
        </Section>
      ) : (
        <Section>
          <ul className="mt-2 grid gap-3 md:grid-cols-2 xl:grid-cols-3">
            {tenants.map((tenant) => (
              <li
                className="flex flex-col gap-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4"
                key={tenant.id}
              >
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <p className="truncate font-bold text-[var(--color-ink)]">{tenant.name}</p>
                    <p className="mt-1 truncate text-sm text-[var(--color-ink-soft)]">/{tenant.slug}</p>
                  </div>
                  <Badge tone={tenant.active ? 'success' : 'danger'}>{tenant.active ? 'Ativo' : 'Inativo'}</Badge>
                </div>
                <p className="text-xs font-semibold text-[var(--color-ink-soft)]">
                  Criado em {formatDateTime(tenant.createdAt)}
                </p>
                <div className="mt-auto flex items-center gap-2">
                  <Button
                    onClick={() => {
                      setEditingTenant(tenant)
                      setIsFormOpen(true)
                    }}
                    size="sm"
                    variant="secondary"
                  >
                    Editar
                  </Button>
                  <Button onClick={() => void handleToggleActive(tenant)} size="sm" variant="secondary">
                    {tenant.active ? 'Desativar' : 'Ativar'}
                  </Button>
                  <Button
                    disabled={tenant.id === 'default'}
                    onClick={() => setTenantToDelete(tenant)}
                    size="sm"
                    title={tenant.id === 'default' ? 'O estabelecimento padrão não pode ser removido.' : 'Remover'}
                    variant="danger"
                  >
                    Remover
                  </Button>
                </div>
              </li>
            ))}
          </ul>
        </Section>
      )}

      {isFormOpen ? (
        <TenantFormModal
          onClose={() => setIsFormOpen(false)}
          onSaved={() => {
            setIsFormOpen(false)
            loadTenants()
          }}
          tenant={editingTenant}
        />
      ) : null}

      <ConfirmModal
        confirmLabel="Remover"
        loading={isDeleting}
        message={`O estabelecimento "${tenantToDelete?.name ?? ''}" e todos os dados vinculados a ele serão removidos.`}
        onCancel={() => setTenantToDelete(null)}
        onConfirm={() => void handleDelete()}
        open={Boolean(tenantToDelete)}
        title="Remover estabelecimento?"
      />
    </PageContainer>
  )
}

function TenantFormModal({
  tenant,
  onClose,
  onSaved,
}: {
  tenant: Tenant | null
  onClose: () => void
  onSaved: () => void
}) {
  const { showToast } = useToast()
  const [name, setName] = useState(tenant?.name ?? '')
  const [slug, setSlug] = useState(tenant?.slug ?? '')
  const [active, setActive] = useState(tenant?.active ?? true)
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const trimmedName = name.trim()
    const trimmedSlug = slug.trim().toLowerCase()

    if (!trimmedName) {
      setError('Informe o nome do estabelecimento.')
      return
    }
    if (!trimmedSlug) {
      setError('Informe o identificador (slug).')
      return
    }
    if (!SLUG_PATTERN.test(trimmedSlug)) {
      setError('O identificador deve conter apenas letras minúsculas, números e hífens (ex.: barbearia-norte).')
      return
    }

    setIsSubmitting(true)
    setError('')

    try {
      const input = { name: trimmedName, slug: trimmedSlug, active }
      if (tenant) {
        await updateTenant(tenant.id, input)
        showToast({ message: 'Estabelecimento atualizado.', tone: 'success' })
      } else {
        await createTenant(input)
        showToast({ message: 'Estabelecimento criado.', tone: 'success' })
      }
      onSaved()
    } catch (submitError) {
      setError((submitError as Error).message)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Modal
      footer={
        <>
          <Button disabled={isSubmitting} onClick={onClose} variant="secondary">
            Cancelar
          </Button>
          <Button loading={isSubmitting} type="submit" form="tenant-form">
            Salvar
          </Button>
        </>
      }
      onClose={onClose}
      open
      size="md"
      subtitle={tenant ? 'Atualize os dados do estabelecimento.' : 'Cadastre um novo estabelecimento.'}
      title={tenant ? 'Editar estabelecimento' : 'Novo estabelecimento'}
    >
      <Form id="tenant-form" onSubmit={handleSubmit}>
        <Field label="Nome">
          <Input
            autoFocus
            className="mt-2"
            onChange={(event) => setName(event.target.value)}
            placeholder="Ex.: Barbearia do Centro"
            value={name}
          />
        </Field>
        <Field
          hint="Usado no link público de agendamento (ex.: meubarbearia). Apenas minúsculas, números e hífens."
          label="Identificador (slug)"
        >
          <Input
            className="mt-2"
            disabled={Boolean(tenant)}
            onChange={(event) => setSlug(event.target.value.toLowerCase())}
            placeholder="meubarbearia"
            value={slug}
          />
        </Field>
        <Field label=" ">
          <Checkbox checked={active} onChange={(event) => setActive(event.target.checked)} />
          <span className="ml-2 text-sm font-bold text-[var(--color-ink)]">Estabelecimento ativo</span>
        </Field>
        {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
      </Form>
    </Modal>
  )
}
